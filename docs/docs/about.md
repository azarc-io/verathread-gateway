# About The Gateway

The Verathread Gateway project allows us to have a re-usable Gateway that can be deployed using helm.

It enables the following capabilities for other Azarc mono repo based projects:


- Statically or dynamically federate multiple GraphQL servers
- Proxy micro front end resources by url where each microservice serves up it's on front end application and the gateway federates them
- Perform JWT verification and propagation of user context to downstream microservices

![](static/gateway.drawio)


### Roadmap

- [x] GraphQL
    * [x] Federated Queries
    * [x] Federated Mutations
    * [x] Federated Subscriptions
- [x] Statically configured routes
- [x] Dynamic routes
    - [ ] Healthiness & Push changes to UI
    - [ ] Nats based registration
    - [x] GraphQL based registration
    - [x] Dapr Cluster based registration
