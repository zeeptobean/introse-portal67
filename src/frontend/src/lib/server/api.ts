import { z } from "zod";
import { createClient, query } from "#lib/server/client.ts";
import type { Client } from "#lib/server/client.ts";
import { userSchema, personSchema, termSchema, list } from "#lib/schemas/shared.ts";
import * as a from "#lib/schemas/academics.ts";
import { documentSchema } from "#lib/schemas/documents.ts";
import type { DocumentInput, DecisionInput } from "#lib/schemas/documents.ts";
import * as t from "#lib/schemas/tuition.ts";
import type { GradeInput } from "#lib/schemas/academics.ts";
import type { LoginInput } from "#lib/schemas/auth.ts";
export function createApi(call: Client) {
  return {
    me: () => call("/me", userSchema),
    login: (input: LoginInput) => call("/auth/login", userSchema, "POST", input),
    logout: () => call("/auth/logout", z.object({ status: z.literal("logged out") }), "POST"),
    academics: {
      terms: () => call("/terms", list(termSchema)),
      sections: (term = "") => call(`/sections${query({ term_id: term })}`, list(a.sectionSchema)),
      enrollments: (term: string) =>
        call(`/me/enrollments${query({ term_id: term })}`, list(a.enrollmentSchema)),
      timetable: (term: string) =>
        call(`/me/timetable${query({ term_id: term })}`, list(a.timetableSchema)),
      exams: (term: string) => call(`/me/exams${query({ term_id: term })}`, list(a.examSchema)),
      register: (id: string) =>
        call("/enrollments", a.enrollmentSchema, "POST", { section_id: id }),
      cancel: (id: string) => call(`/enrollments/${id}`, a.enrollmentSchema, "DELETE"),
      results: () => call("/me/results", a.resultsSchema),
      teaching: () => call("/lecturer/sections", list(a.sectionSchema)),
      roster: (id: string) => call(`/sections/${id}/grades`, list(a.rosterEntrySchema)),
      saveGrades: (id: string, grades: GradeInput[]) =>
        call(`/sections/${id}/grades`, list(a.rosterEntrySchema), "PUT", { grades }),
      history: (id: string) => call(`/sections/${id}/grade-history`, list(a.gradeRevisionSchema)),
    },
    documents: {
      list: () => call("/document-requests", list(documentSchema)),
      get: (id: string) => call(`/document-requests/${id}`, documentSchema),
      submit: (input: DocumentInput) => call("/document-requests", documentSchema, "POST", input),
      decide: (id: string, input: DecisionInput) =>
        call(`/document-requests/${id}`, documentSchema, "PATCH", input),
    },
    tuition: {
      mine: (term: string) =>
        call(`/me/tuition${query({ term_id: term })}`, t.tuitionSummarySchema),
      students: (q = "") => call(`/staff/students${query({ q })}`, list(personSchema)),
      student: (id: string, term: string) =>
        call(`/staff/students/${id}/tuition${query({ term_id: term })}`, t.tuitionSummarySchema),
      payment: (id: string) => call(`/staff/payments/${id}`, t.paymentSchema),
      record: (input: t.PaymentInput) => call("/staff/payments", t.paymentSchema, "POST", input),
      charge: (input: t.ChargeInput) => call("/staff/charges", t.chargeSchema, "POST", input),
      update: (
        id: string,
        input: Partial<Pick<t.Payment, "amount_vnd" | "paid_at" | "note" | "status">>,
      ) => call(`/staff/payments/${id}`, t.paymentSchema, "PATCH", input),
    },
  };
}
export type Api = ReturnType<typeof createApi>;
export { createClient };
