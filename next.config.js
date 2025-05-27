/** @type {import('next').NextConfig} */
const nextConfig = {
  // Enable static exports
  output: 'standalone',
  
  // Configure webpack to handle framer-motion
  webpack: (config, { isServer }) => {
    // Add framer-motion to the list of external modules
    if (!isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        fs: false,
      };
    }
    
    return config;
  },
  
  // Optimize images
  images: {
    unoptimized: true,
  },
  
  // Disable telemetry
  telemetry: {
    disabled: true,
  },
  
  // Experimental features
  experimental: {
    serverActions: true,
    serverComponentsExternalPackages: ['framer-motion']
  },

  // TypeScript configuration
  typescript: {
    // !! WARN !!
    // Dangerously allow production builds to successfully complete even if
    // your project has type errors.
    // !! WARN !!
    ignoreBuildErrors: true,
  }
}

module.exports = nextConfig 