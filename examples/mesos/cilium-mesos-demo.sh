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
desc_rate "This demo shows a brief introduction to using Cilium with Mesos."
desc_rate "The Mesos Master and Slave services as well as Cilium has already been set up."
desc_rate "First, confirm that Cilium is up."

run "cilium status"

desc_rate "Next, start Marathon, the container scheduler".

run "./start_marathon.sh"

run "curl -i -H 'Content-Type: application/json' -d @web-server.json 127.0.0.1:8080/v2/apps"
run "curl -i -H 'Content-Type: application/json' -d @goodclient.json 127.0.0.1:8080/v2/apps"
run "curl -i -H 'Content-Type: application/json' -d @badclient.json 127.0.0.1:8080/v2/apps"
run "cilium endpoint list"


run "screen -S goodclient"
desc_rate "This is the goodclient's log."


run "screen -S badclient"
desc_rate "This is the badclient's log."

desc_rate "If you want to try out this demo yourself, you can do so by    
 following the steps at: http://www.cilium.io/try-mesos              
                                                                      
 If you want to learn more about Cilium:                                                                                        
  - Visit: https://www.cilium.io/                                                                                             
  - Join us on slack: https://cilium.herokuapp.com/"                                                                         

run "ls -l" 
run "echo foo"
