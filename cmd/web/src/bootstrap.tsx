import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import {ApolloClient, ApolloProvider, InMemoryCache} from "@apollo/client";
import {BrowserRouter} from "react-router-dom";

const client = new ApolloClient({
    uri: 'http://localhost:6010/graphql',
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
