import { fail, type RequestEvent } from "@sveltejs/kit";
import { message, superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import { gradeBatchSchema } from "#lib/schemas/academics.ts";
import { requireUser } from "#lib/server/auth.ts";
import { formFailure, verifyUncertain } from "#lib/server/forms.ts";
export async function saveGrades(event: RequestEvent, verifyOnly = false) {
  await requireUser(event, "lecturer");
  const form = await superValidate(event.request, zod4(gradeBatchSchema));
  if (!form.valid) return fail(400, { form });
  const grades = form.data.grades.map((g) => ({ ...g, score: Number(g.score) }));
  const section = event.params.section!;
  const verified = async () => {
    const actual = await event.locals.api.academics.roster(section);
    return grades.every((g) =>
      actual.some((row) => row.enrollment_id === g.enrollment_id && row.score === g.score),
    );
  };
  try {
    if (verifyOnly) {
      if (!(await verified()))
        return message(
          form,
          {
            text: "Current scores do not match every intended edit. Your changes are preserved; inspect the saved scores before deciding what to do.",
            uncertain: true,
            tone: "caution",
          },
          { status: 409 },
        );
    } else await event.locals.api.academics.saveGrades(section, grades);
    form.data.grades = [];
    return message(form, {
      text: verifyOnly
        ? "Saved scores verified against the current roster."
        : `${grades.length} grades saved. These scores are now visible to students.`,
      tone: "success",
    });
  } catch (err) {
    if (await verifyUncertain(err, verified)) {
      form.data.grades = [];
      return message(form, {
        text: "Saved scores verified against the current roster.",
        tone: "success",
      });
    }
    return formFailure(event, form, err);
  }
}
