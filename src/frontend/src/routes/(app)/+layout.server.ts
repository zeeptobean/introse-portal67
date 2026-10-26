import { requireUser } from "#lib/server/auth.ts";
import { readPage } from "#lib/server/forms.ts";
import type { LayoutServerLoad } from "./$types";
export const load: LayoutServerLoad = async (event) => {
  event.depends("portal:session");
  const user = await requireUser(event);
  return readPage(async () => ({ user, terms: await event.locals.api.academics.terms() }));
};
