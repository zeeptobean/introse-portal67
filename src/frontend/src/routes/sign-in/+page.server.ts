import { fail, redirect } from "@sveltejs/kit";
import { superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import { loginSchema } from "#lib/schemas/auth.ts";
import { currentUser, destination } from "#lib/server/auth.ts";
import { formFailure } from "#lib/server/forms.ts";
import type { Actions, PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  const user = await currentUser(event);
  if (user) redirect(303, destination(user, event.url.searchParams.get("next")));
  return { form: await superValidate(zod4(loginSchema)) };
};
export const actions: Actions = {
  default: async (event) => {
    const form = await superValidate(event.request, zod4(loginSchema));
    if (!form.valid) {
      form.data.password = "";
      return fail(400, { form });
    }
    try {
      const user = await event.locals.api.login(form.data);
      redirect(303, destination(user, event.url.searchParams.get("next")));
    } catch (err) {
      form.data.password = "";
      return formFailure(event, form, err);
    }
  },
};
