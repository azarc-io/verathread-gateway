import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import {ApolloClient, ApolloProvider, createHttpLink, InMemoryCache, split} from "@apollo/client";
import {BrowserRouter} from "react-router-dom";
import { WebSocketLink } from "@apollo/client/link/ws";
import {getMainDefinition} from "@apollo/client/utilities";

let apiBaseUri = import.meta.env.PUBLIC_API_BASE_URL;
let wsBaseUri = import.meta.env.PUBLIC_API_BASE_WS_URL;

if (process.env.NODE_ENV == "production") {
    // const domain = window.location.hostname.split('.')[0];
    const bHost = window.location.host;
    apiBaseUri = `https://${bHost}/query`;
    wsBaseUri = `wss://${bHost}/query`;
}

console.log('env', process.env.NODE_ENV)
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
