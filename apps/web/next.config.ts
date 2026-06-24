import type { NextConfig } from "next";
import path from "node:path";

// pnpm monorepo: Turbopack must resolve from workspace root, not apps/web/app
const monorepoRoot = path.resolve(process.cwd(), "../..");

const nextConfig: NextConfig = {
  turbopack: {
    root: monorepoRoot,
  },
  outputFileTracingRoot: monorepoRoot,
};

export default nextConfig;
