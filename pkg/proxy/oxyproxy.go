// Copyright 2017 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cilium/cilium/pkg/lock"
	"github.com/cilium/cilium/pkg/logfields"
	"github.com/cilium/cilium/pkg/node"
	"github.com/cilium/cilium/pkg/policy"
	"github.com/cilium/cilium/pkg/policy/api"
	"github.com/cilium/cilium/pkg/proxy/accesslog"

	"github.com/braintree/manners"
	log "github.com/sirupsen/logrus"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/route"
)

// OxyRedirect implements the Redirect interface for a l7 proxy
type OxyRedirect struct {
	id       string
	toPort   uint16
	epID     uint64
	source   ProxySource
	server   *manners.GracefulServer
	ingress  bool
	nodeInfo accesslog.NodeAddressInfo

	mutex  lock.RWMutex // protecting the fields below
	rules  []string
	router route.Router
}

// ToPort returns the redirect port of an OxyRedirect
func (r *OxyRedirect) ToPort() uint16 {
	return r.toPort
}

func (r *OxyRedirect) updateRules(rules []string) {
	for _, v := range r.rules {
		r.router.RemoveRoute(v)
	}

	r.rules = make([]string, len(rules))
	copy(r.rules, rules)

	for _, v := range r.rules {
		r.router.AddRoute(v, v)
	}
}

func (r *OxyRedirect) getSource() ProxySource {
	return r.source
}

func getOxyPolicyRules(rules []api.PortRuleHTTP) ([]string, error) {
	var l7rules []string

	for _, h := range rules {
		var r string

		if h.Path != "" {
			r = "PathRegexp(\"" + h.Path + "\")"
		}

		if h.Method != "" {
			if r != "" {
				r += " && "
			}
			r += "MethodRegexp(\"" + h.Method + "\")"
		}

		if h.Host != "" {
			if r != "" {
				r += " && "
			}
			r += "HostRegexp(\"" + h.Host + "\")"
		}

		for _, hdr := range h.Headers {
			s := strings.SplitN(hdr, " ", 2)
			if r != "" {
				r += " && "
			}
			r += "Header(\""
			if len(s) == 2 {
				// Remove ':' in "X-Key: true"
				key := strings.TrimRight(s[0], ":")
				r += key + "\",\"" + s[1]
			} else {
				r += s[0]
			}
			r += "\")"
		}

		if !route.IsValid(r) {
			return nil, fmt.Errorf("invalid filter expression: %s", r)
		}
		l7rules = append(l7rules, r)
	}

	return l7rules, nil
}

func translateOxyPolicyRules(l4 *policy.L4Filter) ([]string, error) {
	var l7rules []string

	for _, ep := range l4.L7RulesPerEp {
		rules, err := getOxyPolicyRules(ep.HTTP)
		if err != nil {
			return nil, err
		}
		l7rules = append(rules, l7rules...)
	}

	return l7rules, nil
}

func generateURL(req *http.Request, hostport string) *url.URL {
	newURL := *req.URL
	newURL.Scheme = "http"
	newURL.Host = hostport

	return &newURL
}

// createOxyRedirect creates a redirect with corresponding proxy
// configuration. This will launch a proxy instance.
func createOxyRedirect(l4 *policy.L4Filter, id string, source ProxySource, to uint16) (Redirect, error) {
	for _, ep := range l4.L7RulesPerEp {
		if len(ep.Kafka) > 0 {
			log.Debug("Kafka Parser not supported by Oxy proxy.")
			return nil, fmt.Errorf("unsupported L7 protocol proxy: \"%s\"", l4.L7Parser)
		}
	}

	if l4.L7Parser != policy.ParserTypeHTTP {
		return nil, fmt.Errorf("unknown L7 protocol \"%s\"", l4.L7Parser)
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           ciliumDialerWithContext,
		MaxIdleConns:          2048,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	fwd, err := forward.New(forward.RoundTripper(transport))
	if err != nil {
		return nil, err
	}

	l7rules, err := translateOxyPolicyRules(l4)
	if err != nil {
		return nil, err
	}

	redir := &OxyRedirect{
		id:      id,
		toPort:  to,
		source:  source,
		router:  route.New(),
		ingress: l4.Ingress,
		nodeInfo: accesslog.NodeAddressInfo{
			IPv4: node.GetExternalIPv4().String(),
			IPv6: node.GetIPv6().String(),
		},
	}

	redir.epID = source.GetID()

	redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		record := &accesslog.LogRecord{
			HTTPRequest:       req,
			Timestamp:         time.Now().UTC().Format(time.RFC3339Nano),
			NodeAddressInfo:   redir.nodeInfo,
			TransportProtocol: 6, // TCP's IANA-assigned protocol number
		}

		if redir.ingress {
			record.ObservationPoint = accesslog.Ingress
		} else {
			record.ObservationPoint = accesslog.Egress
		}

		srcIdentity, dstIPPort, err := lookupNewDestFromHttp(req, to)
		if err != nil {
			// FIXME: What do we do here long term?
			log.WithError(err).Error("cannot generate redirect destination url")
			http.Error(w, err.Error(), http.StatusBadRequest)
			record.Info = fmt.Sprintf("cannot generate url: %s", err)
			accesslog.Log(record, accesslog.TypeRequest, accesslog.VerdictError,
				http.StatusBadRequest, accesslog.L7TypeHTTP)
			return
		}

		info, version := getSourceInfo(req.RemoteAddr,
			policy.NumericIdentity(srcIdentity), redir.ingress, redir)
		record.SourceEndpoint = info
		record.IPVersion = version

		if srcIdentity != 0 {
			record.SourceEndpoint.Identity = uint64(srcIdentity)
		}

		record.DestinationEndpoint = getDestinationInfo(dstIPPort,
			redir.ingress, redir)

		// Validate access to L4/L7 resource
		redir.mutex.Lock()
		if len(redir.rules) > 0 {
			rule, _ := redir.router.Route(req)
			if rule == nil {
				http.Error(w, "Access denied", http.StatusForbidden)
				redir.mutex.Unlock()
				accesslog.Log(record, accesslog.TypeRequest, accesslog.VerdictDenied,
					http.StatusForbidden, accesslog.L7TypeHTTP)
				return
			}
			ar := rule.(string)
			log.WithField(logfields.Object,
				logfields.Repr(ar)).Debug("Allowing request based on rule")
			record.Info = fmt.Sprintf("rule: %+v", ar)
		} else {
			log.Debug("Allowing request as there are no rules")
		}
		redir.mutex.Unlock()

		// Reconstruct original URL used for the request
		req.URL = generateURL(req, dstIPPort)

		// log valid request
		accesslog.Log(record, accesslog.TypeRequest, accesslog.VerdictForwarded,
			http.StatusOK, accesslog.L7TypeHTTP)

		ctx := req.Context()
		if ctx != nil {
			marker := GetMagicMark(redir.ingress) | int(record.SourceEndpoint.Identity)
			req = req.WithContext(newIdentityContext(ctx, marker))
		}

		fwd.ServeHTTP(w, req)

		// log valid response
		record.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
		accesslog.Log(record, accesslog.TypeResponse, accesslog.VerdictForwarded,
			http.StatusOK, accesslog.L7TypeHTTP)
	})

	redir.server = manners.NewWithServer(&http.Server{
		Addr:    fmt.Sprintf("[::]:%d", to),
		Handler: redirect,

		// Set a large timeout for ReadTimeout. This timeout controls
		// the time that can pass between accepting the connection and
		// reading the entire request. The default 10 seconds is not
		// long enough.
		ReadTimeout: 120 * time.Second,
	})

	redir.updateRules(l7rules)

	// The following code up until the go-routine is from manners/server.go:ListenAndServe()
	// It was extracted in order to keep the listening on the TCP socket synchronous so that
	// when policies are regenerated, the port is listening for connections before policy
	// revisions get bumped for an endpoint.
	addr := redir.server.Addr
	if addr == "" {
		addr = ":http"
	}

	marker := GetMagicMark(redir.ingress)

	// As ingress proxy, all replies to incoming requests must have the
	// identity of the endpoint we are proxying for
	if redir.ingress {
		marker |= int(source.GetIdentity())
	}

	// Listen needs to be in the synchronous part of this function to ensure that
	// the proxy port is never refusing connections.
	socket, err := listenSocket(addr, marker)
	if err != nil {
		return nil, err
	}

	go func() {
		err := redir.server.Serve(socket.listener)
		if err != nil {
			log.WithError(err).Error("Unable to listen and serve proxy")
		}
	}()

	return redir, nil
}

// UpdateRules replaces old l7 rules of a redirect with new ones.
func (r *OxyRedirect) UpdateRules(l4 *policy.L4Filter) error {
	l7rules, err := translateOxyPolicyRules(l4)
	if err == nil {
		r.mutex.Lock()
		r.updateRules(l7rules)
		r.mutex.Unlock()
	}
	return err
}

// Close the redirect.
func (r *OxyRedirect) Close() {
	r.server.Close()
}
