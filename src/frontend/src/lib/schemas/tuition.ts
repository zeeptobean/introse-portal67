import { z } from "zod";
import {
  actorSchema,
  idSchema,
  list,
  personSchema,
  termSummarySchema,
  text,
} from "#lib/schemas/shared.ts";
const amount = z.number().int().safe();
export const totalsSchema = z.object({
  charged_vnd: amount,
  paid_vnd: amount,
  balance_vnd: amount,
  status: z.enum(["no_charges", "unpaid", "partial", "paid", "credit"]),
});
export const chargeSchema = z.object({
  id: idSchema,
  student_id: idSchema,
  term: termSummarySchema,
  description: z.string(),
  amount_vnd: amount,
  created_at: z.string(),
});
export const paymentEventSchema = z.object({
  id: idSchema,
  actor: actorSchema,
  action: z.enum(["snapshot", "recorded", "updated", "voided"]),
  amount_vnd: amount,
  status: z.enum(["recorded", "voided"]),
  paid_at: z.string(),
  note: z.string().nullable(),
  created_at: z.string(),
});
export const paymentSchema = z.object({
  id: idSchema,
  student_id: idSchema,
  term: termSummarySchema,
  amount_vnd: amount,
  status: z.enum(["recorded", "voided"]),
  paid_at: z.string(),
  note: z.string().nullable(),
  updated_at: z.string(),
  history: list(paymentEventSchema),
});
export const tuitionTermSchema = totalsSchema.extend({
  term: termSummarySchema,
  charges: list(chargeSchema),
  payments: list(paymentSchema),
});
export const tuitionSummarySchema = totalsSchema.extend({
  student: personSchema,
  terms: list(tuitionTermSchema),
});
export const amountInput = z
  .string()
  .trim()
  .regex(/^\d+$/, "Enter a whole VND amount.")
  .refine(
    (value) =>
      Number.isSafeInteger(Number(value)) && Number(value) >= 1 && Number(value) <= 1_000_000_000,
    "Enter an amount from 1 to 1,000,000,000 VND.",
  );
export const paymentFormSchema = z.object({
  amount: amountInput,
  paid_at: z
    .string()
    .regex(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/, "Enter a payment date and time.")
    .refine(
      (value) =>
        z.iso.date().safeParse(value.slice(0, 10)).success &&
        /^([01]\d|2[0-3]):[0-5]\d$/.test(value.slice(11)),
      "Enter a valid date and time.",
    ),
  note: text(500, false).default(""),
  voiding: z.boolean().default(false),
  balance_vnd: amount,
  revision: z.string().default(""),
});
export const chargeFormSchema = z.object({
  amount: amountInput,
  description: text(500),
  balance_vnd: amount,
});
export const studentSearchSchema = z.object({ q: z.string().trim().max(200).default("") });
export type Totals = z.infer<typeof totalsSchema>;
export type Charge = z.infer<typeof chargeSchema>;
export type Payment = z.infer<typeof paymentSchema>;
export type TuitionTerm = z.infer<typeof tuitionTermSchema>;
export type TuitionSummary = z.infer<typeof tuitionSummarySchema>;
export type PaymentForm = z.infer<typeof paymentFormSchema>;
export type ChargeForm = z.infer<typeof chargeFormSchema>;
export interface PaymentInput {
  student_id: string;
  term_id: string;
  amount_vnd: number;
  paid_at: string;
  note: string;
}
export interface ChargeInput {
  student_id: string;
  term_id: string;
  amount_vnd: number;
  description: string;
}
