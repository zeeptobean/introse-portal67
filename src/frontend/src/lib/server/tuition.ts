import { error, fail, redirect, type RequestEvent } from "@sveltejs/kit";
import { superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import {
  paymentFormSchema,
  chargeFormSchema,
  type Charge,
  type Payment,
} from "#lib/schemas/tuition.ts";
import { requireUser } from "#lib/server/auth.ts";
import { formFailure, verifyUncertain } from "#lib/server/forms.ts";
import { selectedTerm } from "#lib/server/terms.ts";
import { ApiError } from "#lib/server/client.ts";
export const localPaymentDate = (value: string) =>
  new Date(Date.parse(value) + 7 * 60 * 60 * 1000).toISOString().slice(0, 16);
export async function tuitionContext(event: RequestEvent) {
  const user = await requireUser(event, "staff");
  const { student, kind, record } = event.params;
  if (
    !student ||
    !["charges", "payments"].includes(kind!) ||
    (kind === "charges" && record !== "new")
  )
    error(404, "Record not found.");
  const terms = await event.locals.api.academics.terms();
  const termId = selectedTerm(event, terms, user.id, "tuition");
  const term = terms.find((t) => t.id === termId);
  if (!term) error(404, "Semester not found.");
  const [summary, payment] = await Promise.all([
    event.locals.api.tuition.student(student, termId),
    record === "new" ? null : event.locals.api.tuition.payment(record!),
  ]);
  if (payment && (payment.student_id !== student || payment.term.id !== termId))
    error(404, "Record not found.");
  return { summary, payment, termId, term, kind: kind! };
}
const ledgerUrl = (saved: Charge | Payment) =>
  `/staff/tuition/students/${saved.student_id}?term=${saved.term.id}&saved=${saved.id}`;
const changedLedger = () =>
  new ApiError("The ledger changed while you were reviewing. Reload it before saving.", 409);
// Each write sets `check` just before it is sent, so a lost response can be
// matched against the records that were not in the ledger it was reviewed on.
export async function saveTuition(event: RequestEvent) {
  await requireUser(event, "staff");
  const { tuition } = event.locals.api;
  if (event.params.kind === "charges") {
    const form = await superValidate(event.request, zod4(chargeFormSchema));
    if (!form.valid) return fail(400, { form });
    let check = async (): Promise<Charge | undefined> => undefined;
    try {
      const context = await tuitionContext(event);
      if (context.summary.balance_vnd !== form.data.balance_vnd) throw changedLedger();
      const input = {
        student_id: context.summary.student.id,
        term_id: context.termId,
        amount_vnd: Number(form.data.amount),
        description: form.data.description,
      };
      const known = new Set(context.summary.terms.flatMap((t) => t.charges.map((c) => c.id)));
      check = async () =>
        (await tuition.student(input.student_id, input.term_id)).terms
          .flatMap((t) => t.charges)
          .find(
            (c) =>
              !known.has(c.id) &&
              c.amount_vnd === input.amount_vnd &&
              c.description === input.description,
          );
      redirect(303, ledgerUrl(await tuition.charge(input)));
    } catch (err) {
      const saved = await verifyUncertain(err, check);
      if (saved) redirect(303, ledgerUrl(saved));
      return formFailure(event, form, err);
    }
  }
  const form = await superValidate(event.request, zod4(paymentFormSchema));
  if (!form.valid) return fail(400, { form });
  let check = async (): Promise<Payment | undefined> => undefined;
  try {
    const context = await tuitionContext(event);
    const { payment, summary, termId } = context;
    if (payment?.status === "voided")
      throw new ApiError("A voided payment cannot be changed.", 409);
    if (
      summary.balance_vnd !== form.data.balance_vnd ||
      (payment && payment.updated_at !== form.data.revision)
    )
      throw changedLedger();
    const values = {
      amount_vnd: Number(form.data.amount),
      paid_at: new Date(`${form.data.paid_at}+07:00`).toISOString(),
      note: form.data.note,
    };
    const voiding = !!payment && form.data.voiding;
    // The backend trims notes and stores a blank one as null.
    const matches = (p: Payment) =>
      p.amount_vnd === values.amount_vnd &&
      Date.parse(p.paid_at) === Date.parse(values.paid_at) &&
      (p.note ?? "") === values.note &&
      (!voiding || p.status === "voided");
    if (payment) {
      check = async () => {
        const actual = await tuition.payment(payment.id);
        return actual.updated_at !== payment.updated_at && matches(actual) ? actual : undefined;
      };
      const saved = await tuition.update(payment.id, {
        ...values,
        ...(voiding ? { status: "voided" as const } : {}),
      });
      redirect(303, ledgerUrl(saved));
    }
    const known = new Set(summary.terms.flatMap((t) => t.payments.map((p) => p.id)));
    check = async () =>
      (await tuition.student(summary.student.id, termId)).terms
        .flatMap((t) => t.payments)
        .find((p) => !known.has(p.id) && matches(p));
    const saved = await tuition.record({
      ...values,
      student_id: summary.student.id,
      term_id: termId,
    });
    redirect(303, ledgerUrl(saved));
  } catch (err) {
    const saved = await verifyUncertain(err, check);
    if (saved) redirect(303, ledgerUrl(saved));
    return formFailure(event, form, err);
  }
}
