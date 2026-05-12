/** @type {import('next').NextConfig} */
const nextConfig = {
  // Abaikan error TypeScript dan ESLint saat build (untuk deployment)
  typescript: {
    ignoreBuildErrors: true,
  },
  eslint: {
    ignoreDuringBuilds: true,
  },
};

export default nextConfig;
