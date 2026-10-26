import { error } from "@sveltejs/kit";
import { readPage } from "#lib/server/forms.ts";
import { selectedTerm } from "#lib/server/terms.ts";
import type { PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  const { parent, params, depends, locals } = event;
  const academics = locals.api.academics;
  const { terms, user } = await parent();
  if (!["classes", "exams"].includes(params.view)) error(404, "Schedule not found.");
  depends("portal:academics");
  const termId = selectedTerm(event, terms, user.id);
  return readPage(async () => ({
    termId,
    view: params.view,
    entries: params.view === "classes" ? await academics.timetable(termId) : [],
    exams: params.view === "exams" ? await academics.exams(termId) : [],
  }));
};
