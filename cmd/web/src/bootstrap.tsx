import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import {ApolloClient, ApolloProvider, createHttpLink, InMemoryCache, split} from "@apollo/client";
import {BrowserRouter} from "react-router-dom";
import { WebSocketLink } from "@apollo/client/link/ws";
import {getMainDefinition} from "@apollo/client/utilities";

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
