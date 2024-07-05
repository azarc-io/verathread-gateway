#!/bin/bash

# this runs each time the container starts

echo "post-start start"

## update hosts file
echo 127.0.0.1 dev.cluster.local | sudo tee -a /etc/hosts
echo 127.0.0.1 k3d-local-registry | sudo tee -a /etc/hosts

## add go bin path to rc
echo 'GOPRIVATE=github.com/azarc-io' >> ~/.zshrc
echo 'GOPRIVATE=github.com/azarc-io' >> ~/.bashrc

echo "post-start complete"
