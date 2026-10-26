import { superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import { gradeBatchSchema } from "#lib/schemas/academics.ts";
import { saveGrades } from "#lib/server/grades.ts";
import { error } from "@sveltejs/kit";
import { readPage } from "#lib/server/forms.ts";
import type { PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  const { parent, params, depends, locals } = event;
  const academics = locals.api.academics;
  await parent();
  depends("portal:grades");
  return readPage(async () => {
    const sections = await academics.teaching();
    const section = sections.find((s) => s.id === params.section);
    if (!section) error(404, "Section not found.");
    const [roster, history] = await Promise.all([
      academics.roster(params.section),
      academics.history(params.section),
    ]);
    return {
      section,
      roster,
      history,
      form: await superValidate(zod4(gradeBatchSchema)),
    };
  });
};

export const actions = {
  save: saveGrades,
  verify: (event: Parameters<typeof saveGrades>[0]) => saveGrades(event, true),
};
