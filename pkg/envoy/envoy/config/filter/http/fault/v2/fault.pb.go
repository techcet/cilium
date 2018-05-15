// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/config/filter/http/fault/v2/fault.proto

package v2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import route "github.com/cilium/cilium/pkg/envoy/envoy/api/v2/route"
import v2 "github.com/cilium/cilium/pkg/envoy/envoy/config/filter/fault/v2"
import _ "github.com/lyft/protoc-gen-validate/validate"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FaultAbort struct {
	// An integer between 0-100 indicating the percentage of requests/operations/connections
	// that will be aborted with the error code provided.
	Percent uint32 `protobuf:"varint,1,opt,name=percent" json:"percent,omitempty"`
	// Types that are valid to be assigned to ErrorType:
	//	*FaultAbort_HttpStatus
	ErrorType            isFaultAbort_ErrorType `protobuf_oneof:"error_type"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *FaultAbort) Reset()         { *m = FaultAbort{} }
func (m *FaultAbort) String() string { return proto.CompactTextString(m) }
func (*FaultAbort) ProtoMessage()    {}
func (*FaultAbort) Descriptor() ([]byte, []int) {
	return fileDescriptor_fault_f7967c49f0120450, []int{0}
}
func (m *FaultAbort) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FaultAbort.Unmarshal(m, b)
}
func (m *FaultAbort) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FaultAbort.Marshal(b, m, deterministic)
}
func (dst *FaultAbort) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FaultAbort.Merge(dst, src)
}
func (m *FaultAbort) XXX_Size() int {
	return xxx_messageInfo_FaultAbort.Size(m)
}
func (m *FaultAbort) XXX_DiscardUnknown() {
	xxx_messageInfo_FaultAbort.DiscardUnknown(m)
}

var xxx_messageInfo_FaultAbort proto.InternalMessageInfo

type isFaultAbort_ErrorType interface {
	isFaultAbort_ErrorType()
}

type FaultAbort_HttpStatus struct {
	HttpStatus uint32 `protobuf:"varint,2,opt,name=http_status,json=httpStatus,oneof"`
}

func (*FaultAbort_HttpStatus) isFaultAbort_ErrorType() {}

func (m *FaultAbort) GetErrorType() isFaultAbort_ErrorType {
	if m != nil {
		return m.ErrorType
	}
	return nil
}

func (m *FaultAbort) GetPercent() uint32 {
	if m != nil {
		return m.Percent
	}
	return 0
}

func (m *FaultAbort) GetHttpStatus() uint32 {
	if x, ok := m.GetErrorType().(*FaultAbort_HttpStatus); ok {
		return x.HttpStatus
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*FaultAbort) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _FaultAbort_OneofMarshaler, _FaultAbort_OneofUnmarshaler, _FaultAbort_OneofSizer, []interface{}{
		(*FaultAbort_HttpStatus)(nil),
	}
}

func _FaultAbort_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*FaultAbort)
	// error_type
	switch x := m.ErrorType.(type) {
	case *FaultAbort_HttpStatus:
		b.EncodeVarint(2<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.HttpStatus))
	case nil:
	default:
		return fmt.Errorf("FaultAbort.ErrorType has unexpected type %T", x)
	}
	return nil
}

func _FaultAbort_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*FaultAbort)
	switch tag {
	case 2: // error_type.http_status
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.ErrorType = &FaultAbort_HttpStatus{uint32(x)}
		return true, err
	default:
		return false, nil
	}
}

func _FaultAbort_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*FaultAbort)
	// error_type
	switch x := m.ErrorType.(type) {
	case *FaultAbort_HttpStatus:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(x.HttpStatus))
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type HTTPFault struct {
	// If specified, the filter will inject delays based on the values in the
	// object. At least *abort* or *delay* must be specified.
	Delay *v2.FaultDelay `protobuf:"bytes,1,opt,name=delay" json:"delay,omitempty"`
	// If specified, the filter will abort requests based on the values in
	// the object. At least *abort* or *delay* must be specified.
	Abort *FaultAbort `protobuf:"bytes,2,opt,name=abort" json:"abort,omitempty"`
	// Specifies the name of the (destination) upstream cluster that the
	// filter should match on. Fault injection will be restricted to requests
	// bound to the specific upstream cluster.
	UpstreamCluster string `protobuf:"bytes,3,opt,name=upstream_cluster,json=upstreamCluster" json:"upstream_cluster,omitempty"`
	// Specifies a set of headers that the filter should match on. The fault
	// injection filter can be applied selectively to requests that match a set of
	// headers specified in the fault filter config. The chances of actual fault
	// injection further depend on the value of the :ref:`percent
	// <envoy_api_field_config.filter.http.fault.v2.FaultAbort.percent>` field. The filter will
	// check the request's headers against all the specified headers in the filter
	// config. A match will happen if all the headers in the config are present in
	// the request with the same values (or based on presence if the *value* field
	// is not in the config).
	Headers []*route.HeaderMatcher `protobuf:"bytes,4,rep,name=headers" json:"headers,omitempty"`
	// Faults are injected for the specified list of downstream hosts. If this
	// setting is not set, faults are injected for all downstream nodes.
	// Downstream node name is taken from :ref:`the HTTP
	// x-envoy-downstream-service-node
	// <config_http_conn_man_headers_downstream-service-node>` header and compared
	// against downstream_nodes list.
	DownstreamNodes      []string `protobuf:"bytes,5,rep,name=downstream_nodes,json=downstreamNodes" json:"downstream_nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HTTPFault) Reset()         { *m = HTTPFault{} }
func (m *HTTPFault) String() string { return proto.CompactTextString(m) }
func (*HTTPFault) ProtoMessage()    {}
func (*HTTPFault) Descriptor() ([]byte, []int) {
	return fileDescriptor_fault_f7967c49f0120450, []int{1}
}
func (m *HTTPFault) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HTTPFault.Unmarshal(m, b)
}
func (m *HTTPFault) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HTTPFault.Marshal(b, m, deterministic)
}
func (dst *HTTPFault) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HTTPFault.Merge(dst, src)
}
func (m *HTTPFault) XXX_Size() int {
	return xxx_messageInfo_HTTPFault.Size(m)
}
func (m *HTTPFault) XXX_DiscardUnknown() {
	xxx_messageInfo_HTTPFault.DiscardUnknown(m)
}

var xxx_messageInfo_HTTPFault proto.InternalMessageInfo

func (m *HTTPFault) GetDelay() *v2.FaultDelay {
	if m != nil {
		return m.Delay
	}
	return nil
}

func (m *HTTPFault) GetAbort() *FaultAbort {
	if m != nil {
		return m.Abort
	}
	return nil
}

func (m *HTTPFault) GetUpstreamCluster() string {
	if m != nil {
		return m.UpstreamCluster
	}
	return ""
}

func (m *HTTPFault) GetHeaders() []*route.HeaderMatcher {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *HTTPFault) GetDownstreamNodes() []string {
	if m != nil {
		return m.DownstreamNodes
	}
	return nil
}

func init() {
	proto.RegisterType((*FaultAbort)(nil), "envoy.config.filter.http.fault.v2.FaultAbort")
	proto.RegisterType((*HTTPFault)(nil), "envoy.config.filter.http.fault.v2.HTTPFault")
}

func init() {
	proto.RegisterFile("envoy/config/filter/http/fault/v2/fault.proto", fileDescriptor_fault_f7967c49f0120450)
}

var fileDescriptor_fault_f7967c49f0120450 = []byte{
	// 384 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0x4f, 0x6a, 0xdb, 0x40,
	0x18, 0xc5, 0xab, 0x7f, 0x76, 0x3d, 0xc2, 0xd4, 0x4c, 0x17, 0x15, 0x5e, 0x14, 0xd9, 0xdd, 0xa8,
	0x06, 0x8f, 0x8a, 0xba, 0x2c, 0x14, 0x6a, 0x97, 0xe0, 0x4d, 0x42, 0x50, 0xbc, 0xca, 0xc6, 0x8c,
	0xa5, 0x71, 0x2c, 0x50, 0x34, 0x62, 0xf4, 0x49, 0x89, 0xcf, 0x93, 0x4b, 0x84, 0xac, 0xbc, 0xcc,
	0x11, 0x72, 0x05, 0xdf, 0x22, 0xcc, 0x8c, 0x85, 0x37, 0x86, 0x6c, 0xc4, 0xf0, 0x7d, 0xbf, 0xf7,
	0xde, 0xcc, 0x13, 0x9a, 0xb2, 0xa2, 0xe1, 0xbb, 0x30, 0xe1, 0xc5, 0x26, 0xbb, 0x0b, 0x37, 0x59,
	0x0e, 0x4c, 0x84, 0x5b, 0x80, 0x32, 0xdc, 0xd0, 0x3a, 0x87, 0xb0, 0x89, 0xf4, 0x81, 0x94, 0x82,
	0x03, 0xc7, 0x23, 0x85, 0x13, 0x8d, 0x13, 0x8d, 0x13, 0x89, 0x13, 0x4d, 0x35, 0xd1, 0x30, 0x38,
	0xe7, 0x78, 0xce, 0x6c, 0xf8, 0x5d, 0x93, 0xb4, 0xcc, 0xe4, 0x46, 0xf0, 0x1a, 0x98, 0xfe, 0x1e,
	0xf7, 0xdf, 0x1a, 0x9a, 0x67, 0x29, 0x05, 0x16, 0xb6, 0x07, 0xbd, 0x18, 0x3f, 0x22, 0x74, 0x21,
	0x7d, 0xfe, 0xad, 0xb9, 0x00, 0xfc, 0x03, 0x75, 0x4b, 0x26, 0x12, 0x56, 0x80, 0x67, 0xf8, 0x46,
	0xd0, 0x9f, 0xf5, 0x5e, 0x0e, 0x7b, 0xcb, 0x9e, 0x98, 0x5e, 0x1a, 0xb7, 0x1b, 0xfc, 0x0b, 0xb9,
	0xf2, 0x9a, 0xab, 0x0a, 0x28, 0xd4, 0x95, 0x67, 0x2a, 0xb0, 0x2f, 0xc1, 0xcf, 0x93, 0xce, 0xe0,
	0xcd, 0x0e, 0x5e, 0x8d, 0xc5, 0xa7, 0x18, 0x49, 0xe6, 0x46, 0x21, 0xb3, 0xaf, 0x08, 0x31, 0x21,
	0xb8, 0x58, 0xc1, 0xae, 0x64, 0xd8, 0x79, 0x3e, 0xec, 0x2d, 0x63, 0xfc, 0x64, 0xa2, 0xde, 0x62,
	0xb9, 0xbc, 0x56, 0xf1, 0xf8, 0x2f, 0x72, 0x52, 0x96, 0xd3, 0x9d, 0xca, 0x75, 0xa3, 0x80, 0x9c,
	0x6b, 0xa7, 0x2d, 0x86, 0x28, 0xcd, 0x7f, 0xc9, 0xc7, 0x5a, 0x86, 0xe7, 0xc8, 0xa1, 0xf2, 0x09,
	0xea, 0x3a, 0x6e, 0x34, 0x25, 0x1f, 0xb6, 0x4b, 0x4e, 0xef, 0x8e, 0xb5, 0x16, 0xff, 0x44, 0x83,
	0xba, 0xac, 0x40, 0x30, 0x7a, 0xbf, 0x4a, 0xf2, 0xba, 0x02, 0x26, 0x3c, 0xcb, 0x37, 0x82, 0x5e,
	0xfc, 0xa5, 0x9d, 0xcf, 0xf5, 0x18, 0xff, 0x41, 0xdd, 0x2d, 0xa3, 0x29, 0x13, 0x95, 0x67, 0xfb,
	0x56, 0xe0, 0x46, 0xa3, 0x63, 0x22, 0x2d, 0x33, 0x69, 0xae, 0xcb, 0x5f, 0x28, 0xe4, 0x92, 0x42,
	0xb2, 0x65, 0x22, 0x6e, 0x15, 0x32, 0x27, 0xe5, 0x0f, 0xc5, 0x31, 0xa9, 0xe0, 0x29, 0xab, 0x3c,
	0xc7, 0xb7, 0x64, 0xce, 0x69, 0x7e, 0x25, 0xc7, 0x33, 0xfb, 0xd6, 0x6c, 0xa2, 0x75, 0x47, 0xfd,
	0xac, 0xdf, 0xef, 0x01, 0x00, 0x00, 0xff, 0xff, 0x60, 0xe5, 0xb8, 0x7f, 0x63, 0x02, 0x00, 0x00,
}
