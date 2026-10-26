import { z } from "zod";
export const idSchema = z
  .string()
  .refine(
    (value) => /^[1-9]\d{0,18}$/.test(value) && BigInt(value) <= 9223372036854775807n,
    "Choose a valid record ID.",
  );
export const roleSchema = z.enum(["student", "lecturer", "staff"]);
export const personSchema = z.object({
  id: idSchema,
  full_name: z.string(),
  email: z.string(),
  institutional_id: z.string().nullish(),
});
export const userSchema = personSchema.extend({ role: roleSchema, is_active: z.boolean() });
export const actorSchema = z.object({ id: idSchema, full_name: z.string() });
export const termSummarySchema = z.object({
  id: idSchema,
  code: z.string(),
  name: z.string(),
  starts_on: z.iso.date(),
  ends_on: z.iso.date(),
});
export const termSchema = termSummarySchema.extend({
  registration_opens_at: z.iso.datetime({ offset: true }),
  registration_closes_at: z.iso.datetime({ offset: true }),
  grading_closes_at: z.iso.datetime({ offset: true }),
});
// Go's nil slices serialize as null, and omitempty slices are left out. Every
// collection boundary normalizes both to an empty array, so readers never need
// a fallback.
export const list = <S extends z.ZodType>(schema: S) =>
  z
    .array(schema)
    .nullish()
    .transform((value) => value ?? []);
export const text = (max: number, required = true) =>
  z
    .string()
    .trim()
    .min(required ? 1 : 0, "Enter a value.")
    .max(max, `Use at most ${max} characters.`);
export type Role = z.infer<typeof roleSchema>;
export type Person = z.infer<typeof personSchema>;
export type User = z.infer<typeof userSchema>;
export type Term = z.infer<typeof termSchema>;
