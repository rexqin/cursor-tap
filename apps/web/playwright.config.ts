import { defineConfig, devices } from '@playwright/test';
import path from 'node:path';

const workspaceRoot = path.resolve(__dirname, '../..');
const baseURL = process.env['BASE_URL'] || 'http://localhost:3000';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  reporter: 'list',
  use: {
    baseURL,
    trace: 'on-first-retry',
  },
  webServer: {
    command: 'pnpm dev',
    url: baseURL,
    reuseExistingServer: !process.env.CI,
    cwd: __dirname,
    timeout: 180_000,
  },
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        // Use system Chrome when Playwright-bundled Chromium is not installed
        channel: 'chrome',
      },
    },
  ],
});
