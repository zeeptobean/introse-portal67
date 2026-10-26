import { redirect, error } from "@sveltejs/kit";
import { requireUser } from "#lib/server/auth.ts";
import { ApiError } from "#lib/server/client.ts";
import type { Actions, PageServerLoad } from "./$types";
export const load: PageServerLoad = () => redirect(303, "/account");
export const actions: Actions = {
  default: async (event) => {
    await requireUser(event);
    try {
      await event.locals.api.logout();
    } catch (err) {
      if (!(err instanceof ApiError && err.status === 401))
        error(503, "Could not confirm sign-out. Try again.");
    }
    event.cookies.delete("session", { path: "/" });
    for (const cookie of event.cookies.getAll())
      if (cookie.name.startsWith("portal67-term-"))
        event.cookies.delete(cookie.name, { path: "/" });
    redirect(303, "/sign-in");
  },
};
