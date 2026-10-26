import type { RequestEvent } from "@sveltejs/kit";
import type { Term } from "#lib/schemas/shared.ts";
import { defaultTerm } from "#lib/domain.ts";
export function selectedTerm(
  event: RequestEvent,
  terms: Term[],
  userId: string,
  scope = "academic",
) {
  const name = `portal67-term-${userId}-${scope}`;
  const selected = defaultTerm(
    terms,
    event.url.searchParams.get("term") ?? event.cookies.get(name),
  );
  // Only a semester ID; record data and private drafts never enter cookies.
  if (event.cookies.get(name) !== selected)
    event.cookies.set(name, selected, {
      path: "/",
      httpOnly: true,
      secure: event.url.protocol === "https:",
      sameSite: "lax",
    });
  return selected;
}
