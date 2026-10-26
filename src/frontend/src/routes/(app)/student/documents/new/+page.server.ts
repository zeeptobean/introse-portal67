import { fail, redirect } from "@sveltejs/kit";
import { superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import { documentInputSchema, type DocumentRequest } from "#lib/schemas/documents.ts";
import { requireUser } from "#lib/server/auth.ts";
import { formFailure, verifyUncertain } from "#lib/server/forms.ts";
import type { Actions, PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) => {
  await requireUser(event, "student");
  return { form: await superValidate(zod4(documentInputSchema)) };
};
export const actions: Actions = {
  default: async (event) => {
    await requireUser(event, "student");
    const form = await superValidate(event.request, zod4(documentInputSchema));
    if (!form.valid) return fail(400, { form });
    const { documents } = event.locals.api;
    // Requests made before this submission, so a lost response can be matched
    // against the request it created.
    let check = async (): Promise<DocumentRequest | undefined> => undefined;
    try {
      const known = new Set((await documents.list()).map((r) => r.id));
      check = async () =>
        (await documents.list()).find(
          (r) =>
            !known.has(r.id) &&
            r.request_type === form.data.request_type &&
            r.purpose === form.data.purpose,
        );
      const record = await documents.submit(form.data);
      redirect(303, `/student/documents/${record.id}`);
    } catch (err) {
      const record = await verifyUncertain(err, check);
      if (record) redirect(303, `/student/documents/${record.id}`);
      return formFailure(event, form, err);
    }
  },
};
