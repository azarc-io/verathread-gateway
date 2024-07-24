#!/bin/bash

# this runs as part of pre-build

echo "on-create start"

# configure git
git config --global --add safe.directory /workspaces/verathread-gateway
git config pull.rebase true

## update hosts file
echo 127.0.0.1 dev.cluster.local | sudo tee -a /etc/hosts
echo 127.0.0.1 k3d-local-registry | sudo tee -a /etc/hosts

## add go bin path to rc
echo 'GOPRIVATE=github.com/azarc-io' >> ~/.zshrc
echo 'GOPRIVATE=github.com/azarc-io' >> ~/.bashrc

### clone the vth dev toolkit project
git clone https://github.com/azarc-io/verathread-dev-toolkit.git /workspaces/verathread-dev-toolkit
git -C /workspaces/verathread-dev-toolkit pull

## install task.dev
go install github.com/go-task/task/v3/cmd/task@latest

## run setup task
task setup

## spin up k3d from the toolkit
pushd /workspaces/verathread-dev-toolkit
task setup
task k3d:create
kubectl wait node --all --for condition=ready --timeout=120s
sleep 15
task k3d:install:charts
popd

echo "on-create complete"
