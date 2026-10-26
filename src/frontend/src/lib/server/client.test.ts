import { afterEach, describe, expect, it, vi } from "vite-plus/test";
import { z } from "zod";
import { createClient } from "#lib/server/client.ts";
import { idSchema, list } from "#lib/schemas/shared.ts";
const session = (value = "token") => ({ get: vi.fn(() => value), set: vi.fn() });
afterEach(() => vi.restoreAllMocks());
describe("server API boundary", () => {
  it("forwards only the session cookie to the private origin and overwrites the client address", async () => {
    const fetcher = vi.fn().mockResolvedValue(new Response('{"id":"9007199254740993"}'));
    const cookies = session();
    const api = createClient("http://127.0.0.1:8080", cookies, fetcher, "203.0.113.9");
    expect(await api("/me", z.object({ id: idSchema }))).toEqual({ id: "9007199254740993" });
    const [url, init] = fetcher.mock.calls[0];
    expect(url.origin).toBe("http://127.0.0.1:8080");
    expect(init.headers.get("Cookie")).toBe("session=token");
    expect(init.headers.get("X-Forwarded-For")).toBe("203.0.113.9");
    expect(init.redirect).toBe("manual");
    expect(init.cache).toBe("no-store");
    expect(cookies.get).toHaveBeenCalledWith("session");
  });
  it("relays session rotation and expiry while excluding unrelated backend cookies", async () => {
    const headers = new Headers();
    headers.append(
      "Set-Cookie",
      "session=new-token; Path=/; Max-Age=3600; Secure; HttpOnly; SameSite=Strict",
    );
    headers.append("Set-Cookie", "other=private; Path=/");
    const cookies = session();
    const api = createClient(
      "http://localhost:8080",
      cookies,
      vi.fn().mockResolvedValue(new Response("{}", { headers })),
    );
    await api("/auth/login", z.object({}), "POST", {});
    expect(cookies.set).toHaveBeenCalledExactlyOnceWith(
      "session",
      "new-token",
      expect.objectContaining({
        path: "/",
        httpOnly: true,
        secure: true,
        sameSite: "strict",
        maxAge: 3600,
      }),
    );
  });
  it("does not follow upstream redirects with the session", async () => {
    const fetcher = vi
      .fn()
      .mockResolvedValue(
        new Response("{}", { status: 302, headers: { Location: "https://untrusted.test/" } }),
      );
    await expect(
      createClient("http://localhost:8080", session(), fetcher)("/me", z.object({})),
    ).rejects.toMatchObject({ status: 302 });
    expect(fetcher).toHaveBeenCalledTimes(1);
  });
  it("marks a lost write uncertain without replaying it", async () => {
    const fetcher = vi.fn().mockRejectedValue(new TypeError("network"));
    await expect(
      createClient("http://localhost:8080", session(), fetcher)(
        "/document-requests",
        z.object({}),
        "POST",
        {},
      ),
    ).rejects.toMatchObject({ status: 0, uncertain: true });
    expect(fetcher).toHaveBeenCalledTimes(1);
  });
  it("treats unavailable reads and explicit domain conflicts differently from unknown writes", async () => {
    const read = createClient(
      "http://localhost:8080",
      session(),
      vi.fn().mockRejectedValue(new TypeError("network")),
    );
    await expect(read("/terms", z.object({}))).rejects.toMatchObject({ uncertain: false });
    const conflict = createClient(
      "http://localhost:8080",
      session(),
      vi.fn().mockResolvedValue(new Response('{"error":"section is full"}', { status: 409 })),
    );
    await expect(conflict("/enrollments", z.object({}), "POST", {})).rejects.toMatchObject({
      message: "section is full",
      status: 409,
      uncertain: false,
    });
  });
  it("rejects a successful response with the wrong contract and preserves uncertainty after a write", async () => {
    const log = vi.spyOn(console, "error").mockImplementation(() => {});
    const api = createClient(
      "http://localhost:8080",
      session(),
      vi.fn().mockResolvedValue(new Response('{"id":123,"secret":"do not log"}')),
    );
    await expect(api("/enrollments", z.object({ id: idSchema }), "POST", {})).rejects.toMatchObject(
      { status: 502, uncertain: true },
    );
    expect(JSON.stringify(log.mock.calls)).not.toContain("do not log");
  });
  it("keeps concurrent request cookies separate", async () => {
    const seen: string[] = [];
    const fetcher = vi.fn().mockImplementation(async (_url, init) => {
      seen.push(init.headers.get("Cookie"));
      return new Response("{}");
    });
    await Promise.all([
      createClient("http://localhost:8080", session("a"), fetcher)("/me", z.object({})),
      createClient("http://localhost:8080", session("b"), fetcher)("/me", z.object({})),
    ]);
    expect(seen).toEqual(["session=a", "session=b"]);
  });
  it("normalizes nil and omitted collections without accepting malformed records", () => {
    expect(list(z.object({ id: idSchema })).parse(null)).toEqual([]);
    expect(z.object({ history: list(idSchema) }).parse({})).toEqual({ history: [] });
    expect(list(z.object({ id: idSchema })).safeParse([{ id: 123 }]).success).toBe(false);
  });
});
