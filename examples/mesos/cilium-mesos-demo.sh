#!/bin/bash

if [[ -z "$HOST_IP" ]]; then
    HOST_IP=192.168.44.11
fi
. $(dirname ${BASH_SOURCE})/util.sh

PWD=$(dirname ${BASH_SOURCE})

# Perform any cleanup to reset state and cleanup at end
function cleanup {
        true
}

trap cleanup EXIT
cleanup

desc_rate "Welcome to the Mesos Getting Started Guide demo."
#desc_rate "This demo shows a brief introduction to using Cilium with Mesos."
#desc_rate "The Mesos Master and Slave services as well as Cilium have already been set up."
#desc_rate "First, confirm that Cilium is up."
run "cilium status"
desc_rate "Next, start Marathon, the container scheduler".
run "./start_marathon.sh"
desc_rate "Start the web-server application and test it."
run "eval curl -i -H 'Content-Type: application/json' -d @web-server.json 127.0.0.1:8080/v2/apps"
echo ""
run "export WEB_IP=`cilium endpoint list | grep web-server | awk '{print $6}'`"
run "curl $WEB_IP:8181/api"
desc_rate "Next, start the goodclient task, retrieving URLs from the web-server."
run "eval curl -i -H 'Content-Type: application/json' -d @goodclient.json 127.0.0.1:8080/v2/apps"
echo ""
#desc_rate "Finally, start the badclient task, retrieving URLs from the web-server."
#run "eval curl -i -H 'Content-Type: application/json' -d @badclient.json 127.0.0.1:8080/v2/apps"
#echo ""
desc_rate "Cilium represents these workloads as endpoints, as observed with the following output:"
run "cilium endpoint list"
desc_rate "Let's observe what these tasks are doing by looking at the logs."
#run "screen -S goodclient"
trap ' ' INT
desc_rate "This is the goodclient's log."
run "./tail_client.sh goodclient"
# hit CTRL-c
#run "screen -S badclient"
#desc_rate "This is the badclient's log."
#run "./tail_client.sh badclient"
# hit CTRL-c
desc_rate "With no policy enforced, both the goodclient and badclient can access /public and /private URLs from the web-server."
desc_rate "Let's apply a Layer-7 policy that only allows the goodclient to access the /public URL."
run "cilium policy import l7-policy.json"
desc_rate "We can observe that the policy got enabled with the following output."
run "cilium endpoint list"
sleep 2
desc_rate "Let's check the goodclient's log again."
sleep 2
run "./tail_client.sh goodclient"
# hit CTRL-c
desc_rate "There you have it! Cilium enforces L7 policy to protect the /private URL on the web-server."
desc_rate "If you want to try out this demo yourself, you can do so by    
 following the steps at: http://www.cilium.io/try-mesos              
                                                                      
 If you want to learn more about Cilium: 
  - Visit: https://www.cilium.io/
  - Join us on slack: https://cilium.herokuapp.com/
 Thanks for watching!"

trap "echo" EXIT
