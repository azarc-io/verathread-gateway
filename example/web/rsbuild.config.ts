import {defineConfig} from '@rsbuild/core';
import {pluginReact} from '@rsbuild/plugin-react';
import {ModuleFederationPlugin} from '@module-federation/enhanced/rspack';
import {pluginSass} from "@rsbuild/plugin-sass";

export default defineConfig({
    plugins: [pluginReact(), pluginSass()],
    server: {
        port: 3001,
    },
    dev: {
        assetPrefix: 'http://localhost:3001/',
    },
    tools: {
        rspack: (config, {appendPlugins}) => {
            // You need to set a unique value that is not equal to other applications
            config.output!.uniqueName = 'example';
            appendPlugins([
                new ModuleFederationPlugin({
                    name: 'example',
                    filename: 'remoteEntry.js',
                    exposes: {
                        "./Counter": './src/components/Counter.tsx'
                    },
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
