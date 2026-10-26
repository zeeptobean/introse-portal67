import { readPage } from "#lib/server/forms.ts";
import { selectedTerm } from "#lib/server/terms.ts";
import type { PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  const { parent, depends, locals } = event;
  const academics = locals.api.academics;
  const { terms, user } = await parent();
  depends("portal:academics");
  const termId = selectedTerm(event, terms, user.id);
  return readPage(async () => {
    const [entries, enrollments, exams] = await Promise.all([
      academics.timetable(termId),
      academics.enrollments(termId),
      academics.exams(termId),
    ]);
    return { termId, entries, enrollments, exams };
  });
};
