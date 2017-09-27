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

desc_rate "Welcome to the Cilium-Mesos Getting Started Guide demo."
desc_rate "This demo shows a brief introduction of Cilium with Mesos by applying an HTTP policy, enforced between a web-server and client."
desc_rate "The Mesos Master and Slave services as well as Cilium have already been set up."
desc_rate "First, confirm that Cilium is up."
run "cilium status"
desc_rate "Next, start Marathon, for container scheduling".
run "./start_marathon.sh"
desc_rate "Start the web-server application and test it."
run "eval curl -i -H 'Content-Type: application/json' -d @web-server.json 127.0.0.1:8080/v2/apps"
echo ""
sleep 5
./generate_client_file.sh client
desc_rate "Next, start the client task, retrieving URLs from the web-server."
run "eval curl -i -H 'Content-Type: application/json' -d @client.json 127.0.0.1:8080/v2/apps"
echo ""
desc_rate "Cilium represents these workloads as endpoints, as observed with the following output:"
run "cilium endpoint list"
desc_rate "Let's observe what this task is doing by looking at the client log."
trap ' ' INT
run "./tail_client.sh client"
# hit CTRL-c
desc_rate "With no policy enforced, the client can access both /public and /private URLs from the web-server."
desc_rate "Let's apply a Layer-7 policy that allows the client to access only the /public URL."
desc_rate "Here's a closer look at the L7 policy that we will apply:"
run "cat l7-policy.json"
desc_rate "Notice we are only allowing traffic labeled \"client\" to access the /public URL on web-server. Let's import the policy into Cilium".
run "cilium policy import l7-policy.json"
desc_rate "We can observe that the policy got enabled with the following output."
run "cilium endpoint list"
desc_rate "Let's check the client's log again."
run "./tail_client.sh client"
# hit CTRL-c
desc_rate "There you have it! Cilium enforces L7 policy to protect the /private URL on the web-server."
desc_rate "If you want to try out this demo yourself, you can do so by following the steps at: 
  http://www.cilium.io/try-mesos
               
  If you want to learn more about Cilium: 
    - Visit: https://www.cilium.io/
    - Join us on slack: https://cilium.herokuapp.com/
  Thanks for watching!"

trap "echo" EXIT
