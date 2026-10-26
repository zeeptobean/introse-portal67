import { fail } from "@sveltejs/kit";
import { superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import { studentSearchSchema } from "#lib/schemas/tuition.ts";
import { requireUser } from "#lib/server/auth.ts";
import { formFailure } from "#lib/server/forms.ts";
import type { Actions, PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  await requireUser(event, "staff");
  return { form: await superValidate(zod4(studentSearchSchema)) };
};
export const actions: Actions = {
  default: async (event) => {
    await requireUser(event, "staff");
    const form = await superValidate(event.request, zod4(studentSearchSchema));
    if (!form.valid) return fail(400, { form });
    try {
      return { form, students: await event.locals.api.tuition.students(form.data.q) };
    } catch (err) {
      return formFailure(event, form, err);
    }
  },
};
