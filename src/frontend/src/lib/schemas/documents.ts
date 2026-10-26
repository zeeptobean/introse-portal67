import { z } from "zod";
import { actorSchema, idSchema, list, personSchema, text } from "#lib/schemas/shared.ts";
export const requestTypeSchema = z.enum(["enrollment_confirmation", "transcript", "other"]);
export const requestStatusSchema = z.enum(["pending", "approved", "rejected"]);
export const requestEventSchema = z.object({
  id: idSchema,
  actor: actorSchema.extend({ role: z.string() }),
  from_status: requestStatusSchema.nullable(),
  to_status: requestStatusSchema,
  comment: z.string().nullable(),
  created_at: z.string(),
});
export const documentSchema = z.object({
  id: idSchema,
  student: personSchema,
  request_type: requestTypeSchema,
  purpose: z.string(),
  status: requestStatusSchema,
  submitted_at: z.string(),
  decided_at: z.string().nullable(),
  decision_comment: z.string().nullable(),
  history: list(requestEventSchema),
});
export const documentInputSchema = z.object({
  request_type: requestTypeSchema,
  purpose: text(1000),
});
export const decisionInputSchema = z.object({
  status: z.enum(["approved", "rejected"]),
  comment: text(1000),
});
export const requestTypes = {
  enrollment_confirmation: "Enrollment confirmation",
  transcript: "Transcript",
  other: "Other document",
} as const;
export type RequestType = z.infer<typeof requestTypeSchema>;
export type RequestStatus = z.infer<typeof requestStatusSchema>;
export type RequestEvent = z.infer<typeof requestEventSchema>;
export type DocumentRequest = z.infer<typeof documentSchema>;
export type DocumentInput = z.infer<typeof documentInputSchema>;
export type DecisionInput = z.infer<typeof decisionInputSchema>;
