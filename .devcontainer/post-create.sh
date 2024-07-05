#!/bin/bash

# this runs at Codespace creation - not part of pre-build

echo "post-create start"

# update the repos
git -C /workspaces/verathread-gateway pull

echo "post-create complete"
