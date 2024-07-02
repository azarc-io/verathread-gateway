import {defineConfig} from '@rsbuild/core';
import {pluginVue} from '@rsbuild/plugin-vue';
import {ModuleFederationPlugin} from "@module-federation/enhanced/rspack";
import {pluginSass} from '@rsbuild/plugin-sass';
import {dependencies} from './package.json';


export default defineConfig({
    plugins: [pluginVue(), pluginSass()],
    server: {
        open: false,
        port: 3000,
        host: '0.0.0.0',
    },
    dev: {
        client: {
            port: 3000,
            host: 'localhost',
        }
    },
    output: {},
    tools: {
        rspack: (config, {appendPlugins}) => {
            // You need to set a unique value that is not equal to other applications
            config.output!.uniqueName = 'mf_shell';

            config.devServer = {
                host: '0.0.0.0',
                allowedHosts: 'all',
                hot: true,
                client: {
                    webSocketURL: 'ws://0.0.0.0:3000/ws'
                }
            }

            appendPlugins([
                new ModuleFederationPlugin({
                    name: 'mf_shell',
                    filename: 'remoteEntry.js',
                    exposes: {},
                    remotes: {},
                    dts: false,
                    dev: {
                        disableDynamicRemoteTypeHints: true,
                    },
                    shared: {
                        ...dependencies,
                        vue: {
                            singleton: true,
                            eager: true,
                            requiredVersion: dependencies['vue']
                        },
                        'vue-router': {
                            singleton: true,
                            eager: true,
                            requiredVersion: dependencies['vue-router']
                        },
                        '@vue/apollo-composable': {
                            singleton: true,
                            eager: true,
                            requiredVersion: dependencies['@vue/apollo-composable']
                        }
                    },
                }),
            ]);
        },
    }
});
