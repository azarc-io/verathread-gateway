# Developing With Docker

You can begin working on this project using docker compose without the need to have `golang` and `node` installed on your
machine.

!!! note

    Please note that the first time you start the environment it could take a few minutes for the entire stack
    to become stable and healthy.

### Requirements

- [Docker](https://www.docker.com/get-started/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Commands

!!! info

    You can run the `task` command in the root of the project to get a list of tasks and their descriptions.

Bring up the dev stack:

```shell
task compose:dev
```

Tear down the dev stack:

```shell
task compose:dev:down
```

Rebuild the docker images:

!!! info

    This is only required if you have made changes to the docker files.
    The below command will rebuild all containers, there are more fine grained tasks available if you prefer
    you can run the `task` command to see a list.

```shell
task compose:rebuild
```

---

You can now begin coding, your changes will be hot reloaded for the backend and HMR updates for the front end.
