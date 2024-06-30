import {defineConfig} from '@rsbuild/core';
import {pluginReact} from '@rsbuild/plugin-react';
import {ModuleFederationPlugin} from '@module-federation/enhanced/rspack';

export default defineConfig({
    plugins: [pluginReact()],
    server: {
        port: 3000,
    },
    dev: {
        assetPrefix: 'http://localhost:3000/',
    },
    tools: {
        rspack: (config, {appendPlugins}) => {
            // You need to set a unique value that is not equal to other applications
            config.output!.uniqueName = 'shell';
            appendPlugins([
                new ModuleFederationPlugin({
                    name: 'shell',
                    filename: 'remoteEntry.js',
                    exposes: {},
                    shared: [
                        {
                            react: {
                                singleton: true, // must be specified in each config
                            },
                        },
                        'react-dom'
                    ],
                }),
            ]);
        },
    },
});
