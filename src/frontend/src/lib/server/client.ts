import { z } from "zod";
import { parseSetCookie } from "cookie";
import type { Cookies } from "@sveltejs/kit";

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public uncertain = false,
  ) {
    super(message);
  }
}
type SessionCookies = Pick<Cookies, "get" | "set">;
export function createClient(
  origin: string,
  cookies: SessionCookies,
  fetcher: typeof fetch = fetch,
  clientAddress?: string,
) {
  const trusted = new URL(origin);
  return async <S extends z.ZodType>(
    path: string,
    schema: S,
    method = "GET",
    body?: unknown,
  ): Promise<z.output<S>> => {
    const url = new URL(`/api${path}`, trusted);
    if (url.origin !== trusted.origin || !path.startsWith("/") || path.startsWith("//"))
      throw new Error("Invalid backend path.");
    const writing = method !== "GET";
    const headers = new Headers({ Accept: "application/json" });
    const session = cookies.get("session");
    if (session) headers.set("Cookie", `session=${encodeURIComponent(session)}`);
    if (clientAddress) headers.set("X-Forwarded-For", clientAddress);
    if (body !== undefined) headers.set("Content-Type", "application/json");
    let response: Response;
    let json: unknown;
    try {
      response = await fetcher(url, {
        method,
        headers,
        body: body === undefined ? undefined : JSON.stringify(body),
        signal: AbortSignal.timeout(20_000),
        redirect: "manual",
        cache: "no-store",
      });
      // Relay only the session cookie, scoped to the frontend host. Preserve the
      // Go expiry/clear operation; never expose the token as page data.
      for (const raw of response.headers.getSetCookie()) {
        const parsed = parseSetCookie(raw);
        if (parsed.name !== "session") continue;
        cookies.set("session", parsed.value ?? "", {
          path: "/",
          httpOnly: true,
          secure: true,
          sameSite: "strict",
          ...(parsed.maxAge !== undefined ? { maxAge: parsed.maxAge } : {}),
          ...(parsed.expires ? { expires: parsed.expires } : {}),
        });
      }
      json = await response.json();
    } catch {
      throw new ApiError(
        writing
          ? "We couldn't confirm whether this was saved. Check current records before trying again."
          : "Could not reach Portal67. Try again.",
        0,
        writing,
      );
    }
    if (!response.ok) {
      const problem = z.object({ error: z.string() }).safeParse(json);
      throw new ApiError(
        problem.success ? problem.data.error : "Portal67 could not complete this request.",
        response.status,
        writing && response.status >= 500,
      );
    }
    const parsed = schema.safeParse(json);
    if (!parsed.success) {
      console.error("Backend response contract failed", {
        path,
        issues: parsed.error.issues.map((issue) => ({ path: issue.path, code: issue.code })),
      });
      throw new ApiError(
        "The service returned an unexpected record. Refresh before continuing.",
        502,
        writing,
      );
    }
    return parsed.data;
  };
}
export type Client = ReturnType<typeof createClient>;
export const query = (values: Record<string, string | undefined>) => {
  const params = new URLSearchParams();
  for (const [key, value] of Object.entries(values)) if (value) params.set(key, value);
  return params.size ? `?${params}` : "";
};
