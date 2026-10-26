import { readPage } from "#lib/server/forms.ts";
import { selectedTerm } from "#lib/server/terms.ts";
import type { PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  const { parent, depends, locals } = event;
  const tuition = locals.api.tuition;
  const { terms, user } = await parent();
  depends("portal:tuition");
  const termId = selectedTerm(event, terms, user.id, "tuition");
  return readPage(async () => ({ termId, summary: await tuition.mine(termId) }));
};
