import { selectedTerm } from "#lib/server/terms.ts";
import { readPage } from "#lib/server/forms.ts";
import type { PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  const { parent, locals } = event;
  const academics = locals.api.academics;
  const { terms, user } = await parent();
  return readPage(async () => ({
    termId: selectedTerm(event, terms, user.id),
    sections: await academics.teaching(),
  }));
};
