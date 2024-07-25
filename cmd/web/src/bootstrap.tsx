import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import {ApolloClient, ApolloProvider, createHttpLink, InMemoryCache, split} from "@apollo/client";
import {BrowserRouter} from "react-router-dom";
import {WebSocketLink} from "@apollo/client/link/ws";
import {getMainDefinition} from "@apollo/client/utilities";

const isProduction = import.meta.env.PUBLIC_PRODUCTION === "true" || import.meta.env.PUBLIC_PRODUCTION === true
const bHost = window.location.host;
let apiBaseUri = `http://${bHost}/graphql`;
let wsBaseUri = `ws://${bHost}/graphql`;

if (isProduction && bHost != 'dev.cluster.local' && bHost != 'localhost') {
    console.log('running in production mode')
    // const domain = window.location.hostname.split('.')[0];
    apiBaseUri = `https://${bHost}/graphql`;
    wsBaseUri = `wss://${bHost}/graphql`;
}

console.log('env', process.env.NODE_ENV)
console.log('production', import.meta.env.PUBLIC_PRODUCTION)
console.log('api url', apiBaseUri)
console.log('socket url', wsBaseUri)

const httpLink = createHttpLink({
    // You should use an absolute URL here
    uri: apiBaseUri,
})

const wsLink = new WebSocketLink({
    uri: wsBaseUri,
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

const client = new ApolloClient({
    link: link,
    cache: new InMemoryCache(),
});

const rootEl = document.getElementById('root');

if (rootEl) {
    const root = ReactDOM.createRoot(rootEl);
    root.render(
        <BrowserRouter>
            <ApolloProvider client={client}>
                <App />
            </ApolloProvider>
        </BrowserRouter>
    );
}
