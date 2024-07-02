import {defineConfig} from '@rsbuild/core';
import {pluginReact} from '@rsbuild/plugin-react';
import {ModuleFederationPlugin} from '@module-federation/enhanced/rspack';
// @ts-ignore
import {dependencies} from './package.json';
import CompressionPlugin from 'compression-webpack-plugin/dist';


export default defineConfig({
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
    },
    output: {
        minify: true,
    },
    tools: {
        rspack: (config, {appendPlugins}) => {
            // You need to set a unique value that is not equal to other applications
            config.output!.uniqueName = 'mf_shell';
            config.output!.publicPath = "auto";

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
                            requiredVersion: dependencies['react'],
                            shareKey: 'react'
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
});
