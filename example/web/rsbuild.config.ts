import {defineConfig, rspack} from '@rsbuild/core';
import {pluginReact} from '@rsbuild/plugin-react';
import {ModuleFederationPlugin} from '@module-federation/enhanced/rspack';
import {pluginSass} from "@rsbuild/plugin-sass";
// @ts-ignore
import {dependencies} from './package.json';
import CompressionPlugin from 'compression-webpack-plugin/dist';

export default defineConfig({
    plugins: [pluginReact(), pluginSass()],
    html: {
      crossorigin: 'anonymous',
    },
    server: {
        open: false,
        port: 3001,
        host: '0.0.0.0',
        compress: true,
        headers: {
            'Access-Control-Allow-Origin': '*',
        },
    },
    dev: {
        assetPrefix: 'http://localhost:3001/',
    },
    output: {
      minify: true
    },
    tools: {
        rspack: (config, {appendPlugins}) => {
            // You need to set a unique value that is not equal to other applications
            config.output!.uniqueName = 'example';
            config.output!.publicPath = "auto";

            appendPlugins([
                new CompressionPlugin(),
                new ModuleFederationPlugin({
                    name: 'example',
                    filename: 'remoteEntry.js',
                    dts: false,
                    dev: {
                        disableDynamicRemoteTypeHints: true,
                    },
                    exposes: {
                        "./Counter": './src/components/Counter.tsx'
                    },
                    shared: {
                        ...dependencies,
                        'react': {
                            requiredVersion: dependencies['react'],
                            singleton: true,
                            shareKey: 'react'
                        },
                        'react-dom': {
                            requiredVersion: dependencies['react-dom'],
                            singleton: true,
                        },
                        'react-router-dom': {
                            requiredVersion: dependencies['react-router-dom'],
                            singleton: true,
                        },
                    },
                }),
            ]);
        },
    },
});
