#!/bin/bash

# this runs each time the container starts

echo "post-start start"

echo 127.0.0.1 dev.cluster.local | sudo tee -a /etc/hosts
echo 127.0.0.1 k3d-local-registry | sudo tee -a /etc/hosts

kubectl wait node --all --for condition=ready --timeout=120s

echo "post-start complete"
