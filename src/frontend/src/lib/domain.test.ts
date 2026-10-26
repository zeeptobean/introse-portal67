import { describe, expect, it } from "vite-plus/test";
import {
  defaultTerm,
  overlaps,
  conflicts,
  money,
  meetingLabel,
  meetingLanes,
} from "#lib/domain.ts";
import type { Meeting, TimetableEntry } from "#lib/schemas/academics.ts";
import type { Term } from "#lib/schemas/shared.ts";
const term = (id: string, starts_on: string, ends_on: string) =>
  ({ id, starts_on, ends_on }) as Term;
const meeting = (weekday: number, starts_at: string, ends_at: string) =>
  ({ weekday, starts_at, ends_at, room: null }) as Meeting;
describe("campus semester selection", () => {
  const terms = [
    term("past", "2026-02-01", "2026-06-30"),
    term("current", "2026-09-01", "2027-01-31"),
    term("next", "2027-02-01", "2027-06-30"),
  ];
  it("honors a valid deep link over the current semester", () =>
    expect(defaultTerm(terms, "past", new Date("2026-10-25T00:00Z"))).toBe("past"));
  it("defaults to the semester on the campus date, across a UTC day boundary", () =>
    expect(defaultTerm(terms, null, new Date("2026-08-31T18:00Z"))).toBe("current"));
  it("selects the nearest upcoming semester in a gap, regardless of source order", () =>
    expect(defaultTerm([...terms].reverse(), "unknown", new Date("2026-08-01T00:00Z"))).toBe(
      "current",
    ));
  it("falls back to the most recent past semester and handles no terms", () => {
    expect(defaultTerm(terms, null, new Date("2028-01-01T00:00Z"))).toBe("next");
    expect(defaultTerm([])).toBe("");
  });
});
describe("schedule comparison", () => {
  it("separates chained overlaps and reuses a lane at an exact boundary", () => {
    const entries = [
      ["a", "08:00", "10:00"],
      ["b", "09:00", "11:00"],
      ["c", "10:00", "12:00"],
    ].map(([id, start, end]) => ({ ...meeting(1, start, end), id }) as TimetableEntry);
    const lanes = meetingLanes(entries);
    expect([...lanes.values()]).toEqual([
      { index: 0, count: 2 },
      { index: 1, count: 2 },
      { index: 0, count: 2 },
    ]);
  });
  const a = meeting(1, "08:00", "10:00");
  it("allows touching times and meetings on other weekdays", () => {
    expect(overlaps(a, meeting(1, "10:00", "11:00"))).toBe(false);
    expect(overlaps(a, meeting(2, "08:00", "10:00"))).toBe(false);
  });
  it("detects an overlap or a contained meeting", () => {
    expect(overlaps(a, meeting(1, "09:30", "11:00"))).toBe(true);
    expect(overlaps(a, meeting(1, "08:30", "09:00"))).toBe(true);
  });
  it("checks every candidate meeting and reports an existing meeting once", () => {
    const existing = { ...a, id: "m" } as TimetableEntry;
    expect(conflicts([a, a], [existing])).toEqual([existing]);
  });
  it("uses an explicit missing-room label and integer VND formatting", () => {
    expect(meetingLabel(a)).toBe("Monday 08:00–10:00 · Room not provided");
    expect(money(600000)).toBe("600,000 VND");
  });
});
