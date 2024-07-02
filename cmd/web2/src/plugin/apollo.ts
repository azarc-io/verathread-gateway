import {ApolloClient, createHttpLink, InMemoryCache, split} from '@apollo/client/core'
import { WebSocketLink } from "@apollo/client/link/ws";
import {getMainDefinition} from "@apollo/client/utilities";


// HTTP connection to the API
const httpLink = createHttpLink({
    // You should use an absolute URL here
    uri: 'http://localhost:6010/graphql',
})

const wsLink = new WebSocketLink({
    uri: "ws://localhost:6010/graphql",
    options: {
        reconnect: true
    }
});

const link = split(
    // split based on operation type
    ({ query }) => {
        const definition = getMainDefinition(query)
        return (
            definition.kind === "OperationDefinition" &&
            definition.operation === "subscription"
        )
    },
    wsLink,
    httpLink
)

// Cache implementation
const cache = new InMemoryCache()

// Create the apollo client
export const apolloClient = new ApolloClient({
    link: link,
    cache,
})
