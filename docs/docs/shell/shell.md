# About The Shell

The shell application is designed to federate other applications that are being
federated and discovered by the Application Gateway.

Applications are discovered using service discovery mechanism that has two implementations
that help simplify submissions for two different use cases.

- Applications that form a cluster with the gateway can use the dapr use case from the common library.
- Applications that are not part of the cluster can use the graphql use case from the common library.

The Gateway can federate the GraphQL api of all discovered applications and present them
as a single api.

The Gateway can proxy the micro front ends of each discovered app.

The Gateway serves the configuration required by the shell in order to construct
its navigation structure.
