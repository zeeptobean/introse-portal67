import { registrationForms, register, cancel } from "#lib/server/registration.ts";
import { error } from "@sveltejs/kit";
import { readPage } from "#lib/server/forms.ts";
import type { PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  const { parent, params, depends, locals } = event;
  const academics = locals.api.academics;
  await parent();
  depends("portal:academics");
  return readPage(async () => {
    const sections = await academics.sections("");
    const section = sections.find((s) => s.id === params.section);
    if (!section) error(404, "Section not found.");
    return {
      ...(await registrationForms()),
      section,
      timetable: await academics.timetable(section.term.id),
    };
  });
};

export const actions = { register, cancel };
