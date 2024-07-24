# Setup Environment

!!! note

    If you are planning to use Kubernetes then you must first follow the instructions outlined in the
    [Verathread Developer Toolkit](https://dev-toolkit-docs.cloud.azarc.dev/gs_setup/) documentation.

    If you have already setup the toolkit then you can skip the `Task.dev` step below and jump to the
    [Initial Setup](#initial-setup) step.

### Requirements

- [Docker](https://www.docker.com/products/docker-desktop/)
- [Docker Compose](https://docs.docker.com/compose/) - *Not required if using Tilt method*

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

### Initial Setup

Using the task cli run the following command

!!! note

    You might be prompted for a password, this is your os password, set up will try to install some tooling that requires
    elevated permissions.

```shell
task setup
```

This will generate the initial .env file if it does not already exist, download the dapr binary and generate the initial
`tilt_config.json` file in the root of the project.

### Environment Files

Edit your `.env` file and fill out any missing values, you can ask another developer or your manager for these values.

Here are some descriptions of important environment settings:

- NAMESPACE: Set this to the first initial and last name e.g. `wael-dev`
    - **Note**: If you are using the developer toolkit then you must set this to the same namespace as defined in the toolkits `.env` file
- BIND_ADDRESS: You should leave this as default unless you want services to be exposed on a different ip
- MONGO_DSN: The mongodb host, this is a mongo dsn eg. `mongo://localhost:27017`
- NATS_ADDRESS: The nats server address eg. `nats://localhost:4222`
- AUTH_DOMAIN: Ask another dev or your manager for this value, this is the auth0 domain to use
- AUTH_CLIENT_ID: Ask another dev or your manager for this value
- AUTH_CLIENT_SECRET: Ask another dev or your manager for this value
- PROJECT_BASE_DIR: This should be the absolute path to the root of the project where you have it checked out
- TILT_ARCH: Set this to `amd64` if you are on an intel machine otherwise set it to `arm64`

### Pick a dev pipeline

!!! warning

    Please note that you can not run multiple appraoches concurrently, if you want to develop using tilt then
    you must stop anything started in the other for eg. by running `docker compose down` if you were using docker compose and 
    want to switch to tilt. This is because both solutions expose the same ports at the host level.

If you would like to develop live on kubernetes then start [here](tilt.md) otherwise go [here](ide.md).
