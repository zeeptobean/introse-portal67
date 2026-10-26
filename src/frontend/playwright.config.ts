import { defineConfig } from "@playwright/test";
import { env } from "node:process";
export default defineConfig({
  testDir: "tests/browser",
  workers: 1,
  retries: 0,
  timeout: 60_000,
  use: {
    baseURL: env.PLAYWRIGHT_BASE_URL ?? "http://localhost:5173",
    channel: env.PLAYWRIGHT_CHANNEL,
    headless: true,
    trace: "retain-on-failure",
  },
  reporter: "list",
});
