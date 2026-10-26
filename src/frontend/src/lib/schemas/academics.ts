import { z } from "zod";
import { actorSchema, idSchema, list, termSchema } from "#lib/schemas/shared.ts";
export const courseSchema = z.object({
  id: idSchema,
  code: z.string(),
  name: z.string(),
  credits: z.number().int().nonnegative(),
  kind: z.enum(["course", "specialized_topic"]),
});
export const meetingSchema = z.object({
  id: idSchema,
  weekday: z.number().int().min(1).max(7),
  starts_at: z.string().regex(/^\d{2}:\d{2}$/),
  ends_at: z.string().regex(/^\d{2}:\d{2}$/),
  room: z.string().nullable(),
});
export const sectionSchema = z.object({
  id: idSchema,
  term: termSchema,
  course: courseSchema,
  section_code: z.string(),
  lecturer: actorSchema,
  capacity: z.number().int(),
  enrolled_count: z.number().int(),
  seats_available: z.number().int(),
  meetings: list(meetingSchema),
  enrollment_id: idSchema.optional(),
});
export const enrollmentSchema = z.object({
  id: idSchema,
  status: z.enum(["enrolled", "cancelled"]),
  registered_at: z.string(),
  cancelled_at: z.string().nullable(),
  section_id: idSchema,
  section_code: z.string(),
  term: termSchema,
  course: courseSchema,
  lecturer_name: z.string(),
  graded: z.boolean(),
});
export const timetableSchema = meetingSchema.extend({
  enrollment_id: idSchema,
  section_id: idSchema,
  section_code: z.string(),
  term: termSchema,
  course_code: z.string(),
  course_name: z.string(),
  lecturer_name: z.string(),
  conflicts_with: list(idSchema),
});
export const examSchema = z.object({
  id: idSchema,
  label: z.string(),
  starts_at: z.string(),
  ends_at: z.string(),
  room: z.string().nullable(),
  section_id: idSchema,
  section_code: z.string(),
  term: termSchema,
  course_code: z.string(),
  course_name: z.string(),
  conflicts_with: list(idSchema),
});
export const courseResultSchema = z.object({
  enrollment_id: idSchema,
  course: courseSchema,
  section_code: z.string(),
  score: z.number().nullable(),
  letter_grade: z.string().nullable(),
  grade_points: z.number().nullable(),
});
export const termResultsSchema = z.object({
  term: termSchema,
  courses: list(courseResultSchema),
  graded_credits: z.number().int(),
  gpa: z.number().nullable(),
});
export const resultsSchema = z.object({
  terms: list(termResultsSchema),
  graded_credits: z.number().int(),
  cumulative_gpa: z.number().nullable(),
});
export const rosterEntrySchema = z.object({
  enrollment_id: idSchema,
  student_id: idSchema,
  institutional_id: z.string().nullable(),
  full_name: z.string(),
  email: z.string(),
  score: z.number().nullable(),
  graded_at: z.string().nullable(),
});
export const gradeRevisionSchema = z.object({
  id: idSchema,
  enrollment_id: idSchema,
  student_id: idSchema,
  full_name: z.string(),
  score: z.number(),
  graded_by: actorSchema,
  created_at: z.string(),
});
export const scoreInput = z
  .string()
  .trim()
  .min(1, "Enter a score; blank is not zero.")
  .refine(
    (value) => /^\d+(?:\.\d{1,2})?$/.test(value) && Number(value) <= 10,
    "Enter a score from 0 to 10 with at most two decimal places.",
  );
export const gradeBatchSchema = z
  .object({
    grades: z
      .array(z.object({ enrollment_id: idSchema, score: scoreInput }))
      .min(1, "Change at least one score.")
      .max(1000),
  })
  .refine(
    (value) => new Set(value.grades.map((g) => g.enrollment_id)).size === value.grades.length,
    { message: "A student can appear only once in a save.", path: ["grades"] },
  );
export const registrationSchema = z.object({ section_id: idSchema });
export const cancellationSchema = z.object({ enrollment_id: idSchema });
export type Course = z.infer<typeof courseSchema>;
export type Meeting = z.infer<typeof meetingSchema>;
export type Section = z.infer<typeof sectionSchema>;
export type Enrollment = z.infer<typeof enrollmentSchema>;
export type TimetableEntry = z.infer<typeof timetableSchema>;
export type Exam = z.infer<typeof examSchema>;
export type CourseResult = z.infer<typeof courseResultSchema>;
export type TermResults = z.infer<typeof termResultsSchema>;
export type Results = z.infer<typeof resultsSchema>;
export type RosterEntry = z.infer<typeof rosterEntrySchema>;
export type GradeRevision = z.infer<typeof gradeRevisionSchema>;
export type GradeInput = { enrollment_id: string; score: number };
export type GradeForm = z.infer<typeof gradeBatchSchema>;
