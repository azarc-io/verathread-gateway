import {createApp, h, provide} from 'vue';
import App from './App.vue';
import './index.css';
import router from "./router";
import {DefaultApolloClient} from '@vue/apollo-composable'
import {apolloClient} from "./plugin/apollo";
import {FetchShellConfigDocument, ShellConfigEventType, SubscribeToShellConfigDocument} from "./gql/graphql";
import useFederatedComponent from "./helpers";
import state from "./store/store";


apolloClient.subscribe({
    query: SubscribeToShellConfigDocument,
    variables: {
        tenant: "abc",
        events: [
            ShellConfigEventType.Initial
        ]
    }
}).subscribe(value => {
    const data = value.data?.shellConfiguration.configuration;
    data.categories?.forEach(c => {
        c?.entries?.forEach(e => {
            console.log('register route', e)
            router.addRoute({
                path: e?.module.path,
                component: useFederatedComponent({
                    remoteUrl: 'http://localhost:3001/remoteEntry.js',
                    moduleToLoad: './Counter',
                    remoteName: 'example',
                })
            })
        })
    })

    state.configuration = data;

    createApp({
        setup() {
            provide(DefaultApolloClient, apolloClient)
            provide('state', state)
        },
        render: () => h(App)
    }).use(router).mount('#root');
})
