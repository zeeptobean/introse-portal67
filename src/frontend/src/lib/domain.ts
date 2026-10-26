import type { Meeting, TimetableEntry } from "#lib/schemas/academics.ts";
import type { Role, Term } from "#lib/schemas/shared.ts";
export const roleHome: Record<Role, string> = {
  student: "/student",
  lecturer: "/lecturer/sections",
  staff: "/staff/documents",
};
export const days = [
  "",
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
  "Sunday",
];
export const money = (value: number) => `${new Intl.NumberFormat("en-US").format(value)} VND`;
export const timestamp = (value: string) =>
  new Intl.DateTimeFormat("en-GB", {
    dateStyle: "medium",
    timeStyle: "short",
    timeZone: "Asia/Ho_Chi_Minh",
  }).format(new Date(value));
export const meetingLabel = (m: Meeting) =>
  `${days[m.weekday]} ${m.starts_at}–${m.ends_at} · ${m.room || "Room not provided"}`;
export function defaultTerm(terms: Term[], requested?: string | null, now = new Date()): string {
  if (requested && terms.some((t) => t.id === requested)) return requested;
  const today = new Intl.DateTimeFormat("en-CA", {
    timeZone: "Asia/Ho_Chi_Minh",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(now);
  const current = terms
    .filter((t) => t.starts_on <= today && t.ends_on >= today)
    .sort((a, b) => b.starts_on.localeCompare(a.starts_on));
  return (
    (
      current[0] ??
      [...terms]
        .filter((t) => t.starts_on > today)
        .sort((a, b) => a.starts_on.localeCompare(b.starts_on))[0] ??
      [...terms].sort((a, b) => b.starts_on.localeCompare(a.starts_on))[0]
    )?.id ?? ""
  );
}
export const overlaps = (a: Meeting, b: Meeting) =>
  a.weekday === b.weekday && a.starts_at < b.ends_at && b.starts_at < a.ends_at;
export const conflicts = (meetings: Meeting[], timetable: TimetableEntry[]) =>
  timetable.filter((entry) => meetings.some((m) => overlaps(m, entry)));
export const registrationOpen = (term: Term) =>
  Date.now() >= Date.parse(term.registration_opens_at) &&
  Date.now() <= Date.parse(term.registration_closes_at);

// Reuse a lane only after its previous meeting ends. One lane count per day
// keeps chained overlaps from painting blocks on top of each other.
export function meetingLanes(entries: TimetableEntry[]) {
  const positions = new Map<string, { index: number; count: number }>();
  for (let day = 1; day <= 7; day++) {
    const meetings = entries
      .filter((e) => e.weekday === day)
      .sort((a, b) => a.starts_at.localeCompare(b.starts_at) || a.id.localeCompare(b.id));
    const ends: string[] = [];
    for (const meeting of meetings) {
      let lane = ends.findIndex((end) => end <= meeting.starts_at);
      if (lane === -1) lane = ends.length;
      ends[lane] = meeting.ends_at;
      positions.set(meeting.id, { index: lane, count: 1 });
    }
    for (const meeting of meetings) positions.get(meeting.id)!.count = ends.length;
  }
  return positions;
}
