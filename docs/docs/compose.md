# Developing With Docker

You can begin working on this project using docker compose without the need to have `golang` and `node` installed on your
machine.

### Requirements

- [Docker](https://www.docker.com/get-started/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Commands

Run the following command in the root of the project to build the docker images:

```shell
docker-compose -f docker-compose-dev.yml build
```
Run the following command to bring up the stack:

```shell
docker-compose -f docker-compose-dev.yml up --watch
```

!!! note

    If you have prevously run docker compose in one of the other Verathread projects then you will need
    to run `docker-compose down` to bring down the existing stack before you bring up the gateways stack.

    If you get an error about a containers name being already in use then run the following command instad.
    ``` shell
    docker-compose -f docker-compose-dev.yml up --watch --remove-orphans
    ```

Run the following command to take down the stack:

```shell
docker-compose down
```

You can now begin coding, your changes will be hot reloaded for the backend and HMR updates for the front end.
