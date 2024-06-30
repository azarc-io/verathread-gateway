/** @type {import('next').NextConfig} */

import { NextFederationPlugin } from '@module-federation/nextjs-mf';

const nextConfig = {
    reactStrictMode: true,
    output: 'export',
    distDir: 'dist',
    webpack: (config, context) => {
        if (!context.isServer) {
            config.plugins.push(
                new NextFederationPlugin({
                    name: 'host',
                    remotes: {},
                    filename: 'static/chunks/remoteEntry.js',
                }),
            )
        }

        return config
    }
};

export default nextConfig;
