import { fail, redirect } from "@sveltejs/kit";
import { superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import { decisionInputSchema } from "#lib/schemas/documents.ts";
import { requireUser } from "#lib/server/auth.ts";
import { formFailure, readPage, verifyUncertain } from "#lib/server/forms.ts";
import type { Actions, PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  await requireUser(event, "staff");
  event.depends("portal:documents");
  return readPage(async () => ({
    record: await event.locals.api.documents.get(event.params.request),
    form: await superValidate(zod4(decisionInputSchema)),
  }));
};
export const actions: Actions = {
  default: async (event) => {
    await requireUser(event, "staff");
    const form = await superValidate(event.request, zod4(decisionInputSchema));
    if (!form.valid) return fail(400, { form });
    try {
      await event.locals.api.documents.decide(event.params.request, form.data);
      redirect(303, event.url.pathname);
    } catch (err) {
      const reviewed = async () =>
        (await event.locals.api.documents.get(event.params.request)).status !== "pending";
      // A 409 means another reviewer decided first; show the current record.
      if (await verifyUncertain(err, reviewed, { conflicts: true }))
        redirect(303, `${event.url.pathname}?reviewed=1`);
      return formFailure(event, form, err);
    }
  },
};
