const { PHASE_DEVELOPMENT_SERVER } = require('next/constants')

module.exports = (phase, { defaultConfig }) => {
    if (phase === PHASE_DEVELOPMENT_SERVER) {
        return {
            /* development only config options here */
            async rewrites() {
                return [
                    {
                        source: '/api/:path*',
                        destination: 'http://localhost:8080/:path*' // Proxy to Backend
                    }
                ]
            }
        }
    }

    return {
        /* config options for all phases except development here */
    }
}