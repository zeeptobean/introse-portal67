import type { Handle } from "@sveltejs/kit/hooks";
import { dev } from "$app/env";
import { BACKEND_URL } from "$app/env/private";
import { createApi, createClient } from "#lib/server/api.ts";
export const handle: Handle = async ({ event, resolve }) => {
  let clientAddress: string | undefined;
  try {
    clientAddress = event.getClientAddress();
  } catch (err) {
    // Vite's SSR runner can omit the socket address. In development, keep
    // Go's connection-level limiter rather than inventing/trusting an IP.
    if (!dev) throw err;
  }
  event.locals.api = createApi(
    createClient(BACKEND_URL, event.cookies, event.fetch, clientAddress),
  );
  const response = await resolve(event);
  if (!event.url.pathname.startsWith("/_app/"))
    response.headers.set("Cache-Control", "private, no-store");
  return response;
};
