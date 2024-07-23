import {defineConfig, rspack} from '@rsbuild/core';
import {pluginReact} from '@rsbuild/plugin-react';
import {ModuleFederationPlugin} from '@module-federation/enhanced/rspack';
// @ts-expect-error complains about json import but works
import {dependencies} from './package.json';
import CompressionPlugin from 'compression-webpack-plugin/dist';
import * as path from "node:path";

export default defineConfig(({ env, command, envMode }) => ({
    plugins: [
        pluginReact(),
    ],
    server: {
        open: false,
        port: 3000,
        host: '0.0.0.0',
        compress: true,
        headers: {
            'Access-Control-Allow-Origin': '*',
        },
    },
    dev: {
        assetPrefix: 'http://localhost:3000/',
        client: {
            port: envMode == 'tilt' ? '80' : '3000',
        }
    },
    output: {
        minify: true,
    },
    tools: {
        rspack: (config, {appendPlugins, env}) => {
            // You need to set a unique value that is not equal to other applications
            config.output!.uniqueName = 'mf_shell';
            config.output!.publicPath = "auto";

            config.resolve = Object.assign(config.resolve || {}, {
                alias: {
                    react: path.resolve('./node_modules/react'),
                    "lib-react": path.resolve('./node_modules/react'),
                    "react-dom": path.resolve('./node_modules/react-dom'),
                },
            }),

            appendPlugins([
                new CompressionPlugin(),
                new ModuleFederationPlugin({
                    name: 'mf_shell',
                    filename: 'remoteEntry.js',
                    dts: false,
                    dev: {
                        disableDynamicRemoteTypeHints: true,
                    },
                    exposes: {},
                    remotes: {},
                    shared: {
                        ...dependencies,
                        react: {
                            singleton: true,
                            requiredVersion: dependencies['react']
                        },
                        'react-dom': {
                            singleton: true,
                            requiredVersion: dependencies['react-dom'],
                        },
                        'react-router-dom': {
                            singleton: true,
                            requiredVersion: dependencies['react-router-dom'],
                        },
                    },
                }),
            ]);
        },
    },
}));
