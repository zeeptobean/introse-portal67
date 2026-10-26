import { error, isHttpError, isRedirect, redirect, type RequestEvent } from "@sveltejs/kit";
import { message, type SuperValidated } from "sveltekit-superforms";
import { ApiError } from "#lib/server/client.ts";
export function formFailure<T extends Record<string, unknown>>(
  event: RequestEvent,
  form: SuperValidated<T>,
  err: unknown,
) {
  if (isRedirect(err) || isHttpError(err)) throw err;
  if (err instanceof ApiError && err.status === 401)
    redirect(303, `/sign-in?next=${encodeURIComponent(event.url.pathname + event.url.search)}`);
  const status =
    err instanceof ApiError && [400, 403, 404, 409, 429].includes(err.status)
      ? (err.status as 400 | 403 | 404 | 409 | 429)
      : 503;
  return message(
    form,
    {
      text: err instanceof Error ? err.message : "Could not complete this action.",
      uncertain: err instanceof ApiError && err.uncertain,
      tone: "error",
    },
    { status },
  );
}
// A write whose outcome is unknown may still have been saved. Read current
// records with `check` and return its result when they show the intended
// change. Pass `conflicts` when a 409 can mean someone else already made it.
// Any other error, or a check that fails, leaves the write uncertain.
export async function verifyUncertain<T>(
  err: unknown,
  check: () => Promise<T>,
  { conflicts = false } = {},
): Promise<T | undefined> {
  if (!(err instanceof ApiError && (err.uncertain || (conflicts && err.status === 409))))
    return undefined;
  try {
    return await check();
  } catch {
    return undefined;
  }
}
export async function readPage<T>(read: () => Promise<T>): Promise<T> {
  try {
    return await read();
  } catch (err) {
    if (isHttpError(err) || isRedirect(err)) throw err;
    const status = err instanceof ApiError && [403, 404].includes(err.status) ? err.status : 503;
    error(
      status,
      status === 404
        ? "We couldn't find this record."
        : err instanceof Error
          ? err.message
          : "Could not load records.",
    );
  }
}
