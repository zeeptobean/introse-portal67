import adapter from "@sveltejs/adapter-node";
import { env } from "node:process";
import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig, lazyPlugins } from "vite-plus";

export default defineConfig({
  server: {
    port: 5173,
    strictPort: true,
  },
  // Svelte markup ignores insignificant whitespace so the formatter can break
  // long chains of tags and blocks onto their own lines.
  fmt: { svelte: true, htmlWhitespaceSensitivity: "ignore" },
  test: { include: ["src/**/*.test.ts"] },
  lint: {
    jsPlugins: [{ name: "vite-plus", specifier: "vite-plus/oxlint-plugin" }],
    rules: { "vite-plus/prefer-vite-plus-imports": "error" },
    options: { typeAware: true, typeCheck: true },
  },
  plugins: lazyPlugins(() => [
    sveltekit({
      compilerOptions: {
        // Force runes mode for the project, except for libraries. Can be removed in svelte 6.
        runes: ({ filename }) =>
          filename.split(/[/\\]/).includes("node_modules") ? undefined : true,
      },

      // The frontend now needs a server for authenticated loads and actions.
      adapter: adapter(),
      // SvelteKit 3 uses paths.origin for a fixed, build-time public origin.
      // Without it, adapter-node expects HTTPS or trusted protocol headers.
      paths: { origin: env.ORIGIN },
    }),
  ]),
});
