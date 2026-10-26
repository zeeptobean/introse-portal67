import { fail, type Action, type RequestEvent } from "@sveltejs/kit";
import { message, superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import { registrationSchema, cancellationSchema, type Enrollment } from "#lib/schemas/academics.ts";
import { requireUser } from "#lib/server/auth.ts";
import { formFailure, verifyUncertain } from "#lib/server/forms.ts";
const hasEnrollment = (event: RequestEvent, match: (e: Enrollment) => boolean) => async () =>
  (await event.locals.api.academics.enrollments("")).some(match);
export const register: Action = async (event) => {
  await requireUser(event, "student");
  const form = await superValidate(event.request, zod4(registrationSchema), {
    id: "register-section",
  });
  if (!form.valid) return fail(400, { form });
  try {
    const saved = await event.locals.api.academics.register(form.data.section_id);
    return message(form, {
      text: `Registered for ${saved.course.code} · section ${saved.section_code} · ${saved.term.name}. Reference ${saved.id}.`,
      tone: "success",
    });
  } catch (err) {
    const registered = hasEnrollment(
      event,
      (e) => e.section_id === form.data.section_id && e.status === "enrolled",
    );
    if (await verifyUncertain(err, registered))
      return message(form, {
        text: "Registration verified in your current records.",
        tone: "success",
      });
    return formFailure(event, form, err);
  }
};
export const cancel: Action = async (event) => {
  await requireUser(event, "student");
  const form = await superValidate(event.request, zod4(cancellationSchema), {
    id: "cancel-enrollment",
  });
  if (!form.valid) return fail(400, { form });
  try {
    const saved = await event.locals.api.academics.cancel(form.data.enrollment_id);
    return message(form, {
      text: `Cancelled ${saved.course.code} · section ${saved.section_code} · ${saved.term.name}. Your seat has been released.`,
      tone: "success",
    });
  } catch (err) {
    const cancelled = hasEnrollment(
      event,
      (e) => e.id === form.data.enrollment_id && e.status === "cancelled",
    );
    if (await verifyUncertain(err, cancelled))
      return message(form, {
        text: "Cancellation verified in your current records.",
        tone: "success",
      });
    return formFailure(event, form, err);
  }
};
export async function registrationForms() {
  return {
    registerForm: await superValidate(zod4(registrationSchema), { id: "register-section" }),
    cancelForm: await superValidate(zod4(cancellationSchema), { id: "cancel-enrollment" }),
  };
}
