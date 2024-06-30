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
                    name: 'example',
                    filename: 'static/runtime/remoteEntry.js',
                    exposes: {
                        "./Counter": "./src/components/Counter.tsx",
                    }
                }),
            )
        }

        return config
    }
};

export default nextConfig;
