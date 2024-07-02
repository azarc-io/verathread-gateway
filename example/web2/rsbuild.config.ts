import {defineConfig} from '@rsbuild/core';
import {pluginVue} from '@rsbuild/plugin-vue';
import {ModuleFederationPlugin} from "@module-federation/enhanced/rspack";
import {pluginSass} from "@rsbuild/plugin-sass";
import {dependencies} from './package.json';

export default defineConfig({
    plugins: [pluginVue(), pluginSass()],
    server: {
        open: false,
        port: 3001,
        host: '0.0.0.0',
    },
    dev: {
        // hmr: false,
        // watchFiles: {
        //     options: {
        //         usePolling: true,
        //         interval: 1000
        //     }
        // },
        client: {
            port: 3001,
            host: 'localhost',
        }
    },
    dev: {
        assetPrefix: 'http://localhost:3001/',
        watchFiles: {
            options: {
                usePolling: true,
                interval: 1000
            }
        }
    },
    tools: {
        rspack: (config, {appendPlugins}) => {
            // You need to set a unique value that is not equal to other applications
            config.output!.uniqueName = 'example';
            config.output!.publicPath = "auto"

            appendPlugins([
                new ModuleFederationPlugin({
                    name: 'example',
                    filename: 'remoteEntry.js',
                    dts: false,
                    dev: {
                        disableDynamicRemoteTypeHints: true,
                    },
                    exposes: {
                        "./Counter": './src/components/Counter.vue'
                    },
                    shared: {
                        ...dependencies,
                        'vue': {
                            singleton: true
                        }
                    },
                }),
            ]);
        },
    }
});
