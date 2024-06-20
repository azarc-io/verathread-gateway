# Setup Environment

!!! note

    If you have already setup k3d, task.dev and other tooling in one of our other projects then
    you can skip this document.

### Task.dev CLI

This is a replacement for make, it provides a much simpler syntax and is well documented.

You can read the task.dev documentation [here](https://taskfile.dev/usage/)

##### OSX

```shell
brew install go-task
```

##### Snap

```shell
sudo snap install task --classic
```

##### NPM

```shell
npm install -g @go-task/cli
```

##### Go Install
```shell
go install github.com/go-task/task/v3/cmd/task@latest
```

!!! info

    You can run `task` on your command line to view a list of available tasks and their descriptions.

### Host Names

In order to make ingress work you will need to add some records to your hosts file, we will be adding two
entries. One is for talking to the ingress and the other is for hosting local docker images.

Start by editing your host file typically found at `/etc/hosts`.

Add these two records and set the ip address to `127.0.0.1` if you are running kubernetes locally, alternatively if you
are running kube on a separate machine you can use that machines ip address.

```text
127.0.0.1 dev.cluster.local
127.0.0.1 k3d-local-registry
```

!!! note

    If you are hosting kube on a separate machine you will need to add that machines ip address to the `cluster.yaml` file
    located in `deployment/k3d/cluster.yaml` 

    In here look for the section starting with:
    ```yaml
      k3s:
        extraArgs:
    ```
    And add another entry like:
    ```yaml
      - arg: --tls-san=127.0.0.1
        nodeFilters:
            - server:*
    ```

### Environment Files

Navigate to the root of the project on your command line and type `task setup` this will install the Tilt cli, copy
over your initial `.env` and `tilt_config.json` files to the root of the project and download the dapr binary.

Edit your `.env` file and fill out any missing values, you can ask another developer or your manager for these values.

Here are some descriptions of important environment settings:

- NAMESPACE: Set this to the first initial and last name e.g. `wael-dev`
- BIND_ADDRESS: You should leave this as default unless you want services to be exposed on a different ip
- NATS_ADDRESS: This is optional unless you have auto app discovery enabled (off by default)
- AUTH_DOMAIN: Ask another dev or your manager for this value
- AUTH_CLIENT_ID: Ask another dev or your manager for this value
- AUTH_CLIENT_SECRET: Ask another dev or your manager for this value
- PROJECT_BASE_DIR: This should be the absolute path to the root of the project where you have it checked out
- DAPRD_BINARY: Make sure you update this whenever a new release is available in the dapr repo
- DAPRD_BIN_PATH: Absolute path to the `.data/bin` in the root of the project
- TILT_ARCH: Set this to `amd64` if you are on an intel machine otherwise set it to `arm64`
