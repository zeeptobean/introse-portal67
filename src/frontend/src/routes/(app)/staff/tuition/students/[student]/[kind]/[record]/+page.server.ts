import { superValidate } from "sveltekit-superforms";
import { zod4 } from "sveltekit-superforms/adapters";
import { paymentFormSchema, chargeFormSchema } from "#lib/schemas/tuition.ts";
import { tuitionContext, saveTuition, localPaymentDate } from "#lib/server/tuition.ts";
import { readPage } from "#lib/server/forms.ts";
import type { Actions, PageServerLoad } from "./$types";
export const load: PageServerLoad = async (event) =>
  readPage(async () => {
    const context = await tuitionContext(event);
    const { payment, summary } = context;
    return {
      ...context,
      paymentForm:
        context.kind === "payments"
          ? await superValidate(
              {
                amount: payment?.amount_vnd.toString() ?? "",
                paid_at: localPaymentDate(payment?.paid_at ?? new Date().toISOString()),
                note: payment?.note ?? "",
                voiding: false,
                balance_vnd: summary.balance_vnd,
                revision: payment?.updated_at ?? "",
              },
              zod4(paymentFormSchema),
              { errors: false },
            )
          : null,
      chargeForm:
        context.kind === "charges"
          ? await superValidate(
              { amount: "", description: "", balance_vnd: summary.balance_vnd },
              zod4(chargeFormSchema),
              { errors: false },
            )
          : null,
    };
  });
export const actions: Actions = { default: saveTuition };
