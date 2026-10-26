import { registrationForms, register, cancel } from "#lib/server/registration.ts";
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
    const [sections, enrollments, timetable] = await Promise.all([
      academics.sections(termId),
      academics.enrollments(termId),
      academics.timetable(termId),
    ]);
    return {
      ...(await registrationForms()),
      termId,
      sections,
      enrollments,
      timetable,
    };
  });
};

export const actions = { register, cancel };
