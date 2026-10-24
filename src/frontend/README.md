# Portal67 frontend

SvelteKit with TypeScript and Vite+. See the [root README](../../README.md) for the full local setup.

```sh
pnpm install --frozen-lockfile
pnpm dev
```

The frontend runs at http://localhost:5173. During development, `/api` requests are proxied to the Go backend at http://127.0.0.1:8080. Update `server.proxy` in `vite.config.ts` if you change the backend address.

```sh
pnpm check       # Svelte and TypeScript diagnostics
pnpm lint        # Vite+ format, lint, and TypeScript checks
pnpm format      # Format supported files
pnpm build       # Production build
pnpm preview     # Preview the production build locally
```

For deployment, choose a SvelteKit adapter for the hosting environment and configure the web server to forward `/api` to the Go backend. The development proxy is not production routing.
