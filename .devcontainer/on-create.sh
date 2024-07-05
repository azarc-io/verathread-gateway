#!/bin/bash

# this runs as part of pre-build

echo "on-create start"

# clone repos
git clone https://github.com/azarc-io/verathread-gateway.git /workspaces/verathread-gateway

## install golang
#wget -c https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
#sudo tar -C /usr/local/ -xzf go1.22.0.linux-amd64.tar.gz
## shellcheck disable=SC2016
## dont want variables to expand
#echo '
## set PATH so it includes /usr/local/go/bin if it exists
#if [ -d "/usr/local/go/bin" ] ; then
#    PATH="/usr/local/go/bin:$PATH"
#fi
#' >> ~/.bashrc

# add go bin path to rc
echo 'PATH=$PATH:/go/bin' >> ~/.zshrc
echo 'PATH=$PATH:/go/bin' >> ~/.bashrc

# install task.dev
go install github.com/go-task/task/v3/cmd/task@latest

# run setup task
task setup

echo "on-create complete"
