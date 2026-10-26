import { selectedTerm } from "#lib/server/terms.ts";
import { readPage } from "#lib/server/forms.ts";
import type { PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  const { parent, depends, locals } = event;
  const academics = locals.api.academics;
  const { user } = await parent();
  depends("portal:academics");
  return readPage(async () => {
    const results = await academics.results();
    return {
      results,
      termId: selectedTerm(
        event,
        results.terms.map((t) => t.term),
        user.id,
      ),
    };
  });
};
