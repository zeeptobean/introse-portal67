import { error, redirect, type RequestEvent } from "@sveltejs/kit";
import { ApiError } from "#lib/server/client.ts";
import type { Role, User } from "#lib/schemas/shared.ts";
import { roleHome } from "#lib/domain.ts";
export async function currentUser(event: RequestEvent): Promise<User | null> {
  if (!event.cookies.get("session")) return null;
  if (!event.locals.userPromise)
    event.locals.userPromise = event.locals.api.me().catch((err) => {
      if (err instanceof ApiError && err.status === 401) {
        event.cookies.delete("session", { path: "/" });
        return null;
      }
      throw err;
    });
  return event.locals.userPromise;
}
export async function requireUser(event: RequestEvent, role?: Role): Promise<User> {
  let user: User | null;
  try {
    user = await currentUser(event);
  } catch {
    error(503, "Portal67 is temporarily unavailable. Try again.");
  }
  if (!user)
    redirect(303, `/sign-in?next=${encodeURIComponent(event.url.pathname + event.url.search)}`);
  if (role && user.role !== role) error(403, "This page is not available for your account.");
  return user;
}
export function destination(user: User, next: string | null) {
  return next &&
    (next.startsWith(`/${user.role}/`) || next === `/${user.role}` || next === "/account")
    ? next
    : roleHome[user.role];
}
