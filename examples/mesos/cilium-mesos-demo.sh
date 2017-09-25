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
desc_rate "

desc_rate "If you want to try out this demo yourself, you can do so by    
 following the steps at: http://www.cilium.io/try-mesos              
                                                                      
 If you want to learn more about Cilium:                                                                                        
  - Visit: https://www.cilium.io/                                                                                             
  - Join us on slack: https://cilium.herokuapp.com/"                                                                         

run "ls -l" 
run "echo foo"
