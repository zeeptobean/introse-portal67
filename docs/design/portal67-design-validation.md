# Portal67 design coverage and validation

Companion to the [experience design](portal67-experience-design.md), [page/flow specifications](portal67-pages-and-flows.md), and [component specification](portal67-components.md).

This is a proposed review and test plan, not a report of passed product tests. Source inspection establishes what the current contracts describe; it does not establish that the running application behaves correctly. The frontend inspected for this work is a starter page, so the screen designs are forward-looking.

## 1. Brief-to-design traceability

| PA1 ID | Brief obligation | Pages / journeys | Components | Acceptance scenarios |
| --- | --- | --- | --- | --- |
| F01 Must | Browse sections and specialized topics; seats and schedules; submit/cancel; seat/time checks | S02/S03, J01/J02 | C04–C08, C12/C13/C16 | A01–A06 |
| F02 Must | Semester charges, payment records, balances/status; staff simulated payment updates | S07, A03–A05, J06; A06 supporting extension | C03/C04/C08/C10/C15/C16 | A13–A17 |
| F03 Must | Lecturer demo grades; student course results; estimated semester/cumulative GPA under sample rules | L01/L02/S06, J04 | C03/C08/C09/C13/C19 | A09–A12 |
| F04 Must | Weekly timetable, exam list, overlapping sessions, registration conflicts | S01/S03/S04/S05, J01/J03 | C05/C07/C08/C12 | A03/A07/A08 |
| F05 Could | Future room/equipment availability and booking | Deferred concept in main document §1.4 | Not placed in MVP navigation | A28 verifies absence of misleading controls |
| F06 Must | Student request; staff approval/rejection with comments; status/history | S08–S10, A01/A02, J05 | C03/C11/C14–C18 | A18–A22 |
| F07 Must | Different permissions; own student data and assigned lecturer sections | U01–U03 and all record pages | C01/C02/C15 | A23/A24 |
| F08 Must | Invalid/incomplete inputs and understandable errors | All editable pages; U03 | C03/C12–C16/C19 | A02–A06/A10/A16/A19/A22/A25–A27 |

PA0's three pain points also have explicit responses: registration overload → D04/C13 and A05; insufficient timetable planning → D02/D03/S03 and A03; insufficient GPA tracking → D07/S06 and A09/A11. Success is not defined by the number of dashboard widgets.

## 2. Current contract versus required design work

The existing backend is a useful feasibility reference. The final requirements may change it. Do not remove a needed user experience merely because a field is currently missing; document the dependency instead.

| Design area | Observed support | Gap / decision |
| --- | --- | --- |
| Login and roles | Email/password login; one role; assigned-section and student ownership checks | UI and help needed; no password-reset/provisioning flow. O07 |
| Registration | Capacity, duplicate section/course, class overlap, ended-term and graded-cancellation rules | No registration window, prerequisites, credit caps, atomic swap, or waitlist. O01 |
| Schedule comparison | Meeting and enrollment records exist | One-section preview is presentation-only and not a server reservation |
| Day agenda | Weekly meetings and term dates exist | No holiday/cancellation exceptions; S01 labels the derivation |
| Exams | Exam entries and overlapping IDs | UI must show actual affected pairs; no exam-edit/report workflow |
| GPA | Grades, points, credits, semester/cumulative aggregates; sample bands in data | Readable configurable rules not exposed as a dedicated contract; O02/O09 |
| Grade save | Batch validation and all-or-none save | No publish stage, deletion, or revision-based stale-write rejection; O06 |
| Tuition | Charges, payments, updates/voiding, derived balances | No due date, real transfer, charge correction, full correction-event audit, or official receipt; O05 |
| Documents | Three types, required purpose/comment, terminal decisions, event history | No withdrawal, revision, issuance, attachments, expected completion time, or assignment; O04 |
| Lists | Status/term/query filters vary by endpoint | Queue search/sort and pagination behavior require explicit UI/data-scope work; do not filter only a loaded page |
| Retry certainty | Existing object reads can verify many updates | New requests/payments/charges have no operation-reference contract; reliable deduplication needs O10 |
| Notifications | Current records can be read | No unread marker or push channel; no global bell or unread-count claims |

The additional requirements with the greatest impact are unknown-outcome reconciliation, readable GPA rules, and concurrent grade correction. These deserve early requirements review. The design is concrete enough to review while preserving an honest distinction between a proposed behavior and an already supported capability.

## 3. Acceptance scenarios

These are testable design outcomes for wireframe review and later implementation verification. No pass/fail results have been recorded yet.

| ID | Scenario and setup | Required observable result |
| --- | --- | --- |
| A01 | Student finds a specialized topic with available seats | Course kind, section, every meeting, credits, and seat availability are visible; confirmed registration appears in My registrations and timetable |
| A02 | Two students attempt the final available seat | Only the confirmed owner sees Registered; the other sees Full and retains search/term; no phantom enrollment or optimistic success |
| A03 | CSC10003 section 01 is compared against CSC13002 section 01 | Monday 09:30–11:30 overlap is named; candidate is Preview only; actual registration remains unchanged when preview closes |
| A04 | Existing course is selected in another section, or registered section is submitted again | Specific existing registration is identified; no duplicate active course registration; no silent swap |
| A05 | Response is lost after a successful registration write | UI enters unknown-outcome verification; re-read finds the exact enrollment; no duplicate submission and no false failure |
| A06 | Student cancels an active ungraded registration, then tries a graded/ended-term one | Confirmation explains seat loss; successful cancellation updates schedule; blocked cancellation explains its rule |
| A07 | Weekly timetable includes multiple meetings, Saturday, missing room, and a seeded overlap | Every meeting appears in grid and agenda; missing room has text; overlapping blocks remain readable; adjacent end/start times are not conflicts |
| A08 | CSC13002 and MTH00003 exams overlap on 5 Jan 2027 | Both exam identities and 09:00–10:00 intersection are visible; no false resolution or rescheduling control |
| A09 | Seed GPA: 4 credits at 3.5 points, 4 at 3 points, 3 at 0 points | GPA is 2.36 over 11 graded credits; formula and sample scale are discoverable |
| A10 | Lecturer enters blank, 0, 10, 11.5, and three decimal places | Blank remains ungraded; 0 and 10 accepted; invalid entries identified by student; whole batch remains unsaved until fixed |
| A11 | Partial grades, no grades, repeated course, and failed result load | Ungraded credits excluded; no-grade GPA is unavailable; current repeat treatment is explained; load failure is not zero GPA |
| A12 | Lecturer edits, filters out one changed row, navigates away, then returns | Hidden edit count remains; save scope is clear; leave guard protects edits; confirmed save uses returned values. Concurrent-edit limitation is disclosed until version checks exist |
| A13 | Charges 1,000,000 VND and recorded payment 400,000 VND | Balance 600,000 VND, Partly paid; visible ledger reconciles |
| A14 | No records; paid records; then payment exceeds charges | No tuition records differs from Paid; negative mathematical balance is displayed as credit, not overdue debt |
| A15 | Staff voids a 400,000 VND payment against 1,000,000 VND charges | Review predicts balance 1,000,000; confirmed totals agree; voided record remains visible and is read-only |
| A16 | Invalid payment amount/date or two identical student names | Invalid value has a field remedy; correct institutional ID retained through review; no wrong-person save |
| A17 | Payment-create response is lost | UI does not blindly repeat or infer certainty from a similar row; O10 dependency is exposed; staff can inspect latest records |
| A18 | Student submits each supported request type | Own identity is prefilled; purpose required; confirmed receipt has ID, Pending status, and submission history |
| A19 | Missing purpose; oversized comment; server failure | Field-level error and summary point to the input; valid entered content remains intact |
| A20 | Staff approves and rejects separate requests | Both require student-visible comment; detail shows actual decision/actor/time; student sees the same state/history |
| A21 | Student opens an approved request | It is clear that no legally valid document was issued; no nonfunctional download control |
| A22 | Two staff members decide the same pending request | The second sees the first decision, cannot overwrite it, and does not receive false success |
| A23 | Student changes a private-record identifier; lecturer opens another lecturer's roster | No unauthorized data, counts, or mutation; safe missing/access page; shell returns to allowed destination |
| A24 | Session expires during editing; another account signs in | No automatic replay; private draft/data are not disclosed to the other account; correct role destination |
| A25 | A list load fails, a filter finds nothing, and the dataset is actually empty | Three distinct truthful states, each with the right recovery; no manufactured empty dashboard |
| A26 | Keyboard-only navigation, modal open/close, and error correction | Visible unobscured focus, correct order, named controls, summary links, and focus return |
| A27 | 320 px width, 200% zoom, long Vietnamese name, long purpose, mobile keyboard | No lost task/action/essential field; readable wrapping; agenda and roster alternatives; sticky controls do not cover inputs |
| A28 | Review every visible navigation item and action | Every control maps to specified supported behavior; no dead Resources, Pay now, Publish, Download official document, or fake notification controls |

Use isolated synthetic records for destructive examples; do not use the real team member identities from the briefs as demo personas. The GPA and schedule fixtures above are tied to the inspected seed data, but are design verification examples until executed against a running implementation.

## 4. Usability validation before implementation hardens

Recruit, if feasible, 3–5 students with different familiarity with the current portal, one lecturer or credible grade-entry proxy, and one administrative user or process-informed proxy. This is a small formative study, not representative evidence of all HCMUS users. Record proxy participation explicitly. Include a keyboard-only session and a narrow-screen schedule check.

Use realistic tasks without naming the controls to click. Observe behavior before explaining the design. Record wrong turns, misunderstandings, completion, assistance, and the participant's explanation of the outcome. Obtain permission for recording if used; written observations are sufficient.

| Task | What to observe | Proposed design target, not an achieved metric |
| --- | --- | --- |
| Find the room for the next scheduled class | Does the student inspect Today or become lost in modules? | Correct answer within about 15 seconds without coaching |
| Find a suitable section and identify a conflicting alternative | Can they distinguish course, section, preview, and registered state? | Complete within about 90 seconds; zero false beliefs that preview reserves a seat |
| Explain why GPA excludes one course | Do they find coverage and rules without a tutorial? | Correct explanation within about 60 seconds; no confusion between blank and zero |
| Find a rejected request and explain what happened | Is comment more discoverable than status alone? | Correct decision/reason without resubmitting |
| Enter five grades including a zero and correct an invalid value | Does the save scope and student identity remain clear? | Zero wrong-student saves; all participants identify saved versus unsaved values |
| Correct a payment for one of two namesakes | Does identity persist into review? | Correct student and projected balance, no money-transfer misconception |
| Handle an uncertain registration response | Do users repeatedly click or know where to verify? | User can identify that outcome is unknown and check actual registration |

Treat a target miss as a design signal, not a participant failure. If students mistake Preview for Registered, strengthen both label and placement before adding polish. If staff miss the identity in review, move it into the action region. If page-specific terminology causes hesitation, revise vocabulary before adding onboarding.

Do not collect feature usage or session recordings merely because analytics are available. For this prototype, direct task observations and aggregate completion notes are enough. Production telemetry is a separate decision.

## 5. Review artifacts and completeness gate

For each page, the designer's later screen set should include its normal state, meaningful empty/error state, and mutation outcome if applicable. The page specification defines those states even when only representative compositions are illustrated.

Minimum visual review set:

1. Registration desktop: ordinary offering, full section, conflict preview, success, unknown outcome, cancellation confirmation.
2. Classes desktop and mobile agenda; exam conflict pair.
3. Results with mixed grades and the calculation disclosure; no grades.
4. Tuition with partial/credit balance and voided record; staff before/after review.
5. Document form, pending receipt, rejected detail, staff review and already-decided state.
6. Lecturer roster with dirty/error rows, save scope, and narrow layout.
7. Sign-in, wrong-role/missing record, service error, and expired editing session.

The accompanying visual reference illustrates the main composition and visual direction. It is not a tested clickable prototype or a replacement for the full state specifications above.

Before design sign-off, reviewers should be able to answer:

- Does every F01/F02/F03/F04/F06/F07/F08 obligation have a page, action, outcome, and failure behavior?
- Are F05 and other optional features absent from the working navigation?
- Are the status words identical across student and staff views?
- Is course versus section, selected semester versus cumulative, preview versus saved, and approval versus issuance unambiguous?
- Can every important action be completed by keyboard and on a narrow viewport?
- Does every essential visual element have a task reason, and does each claimed fact have data behind it?
- Are unresolved policies labeled, assigned for review, and absent as invented product promises?

## 6. Scope-aware delivery sequence

The briefs allocate only three development weeks, so build coherent vertical journeys rather than many half-working pages. This is sequencing guidance, not a new estimate or replacement for the PA1 schedule.

| Sequence | Complete experience | Safe simplifications |
| --- | --- | --- |
| Foundation | Sign-in, role shell, semester context, fields, status and failure language | Text-led identity; no avatar/theme preferences |
| Academic core | Registration → actual enrollment → timetable → cancellation; exam list/conflicts | Text comparison before a polished mini-grid; no saved planner or waitlist |
| Results loop | Assigned roster → validated grade save → student results and explainable GPA | No import, charts, forecasting, publishing stages |
| Administration loop | Request → staff decision → student history; tuition lookup → record update → reconciled student view | Single staff role; no notifications, workflow builder, or issued files |
| Finishing | Keyboard/reflow/error validation, visual consistency, realistic fixtures, scope documentation | Cut supplementary home modules and charge-authoring polish before any Must flow |

If time is constrained, retain all brief Must behavior and simplify presentation: a schedule agenda may precede an enhanced weekly visualization, a single-column section detail can replace a side panel, and seeded charges can replace A06. Never remove authoritative save feedback, basic conflict visibility, role protection, or required comments to preserve a decorative feature.

## 7. Review status

Completed in this documentation task: both briefs read; key requirement/constraint pages visually inspected; current relevant contracts inspected; page/component decisions documented; requirement mapping and proposed validation scenarios written.

Documentation checks performed on 25 October 2026:

- All local Markdown links resolved; inventory contains 21 page/state specifications, 19 component contracts, and 28 acceptance scenarios, with F01–F08 accounted for.
- The worked GPA example recalculates to 2.36; the illustrative tuition example reconciles to 600,000 VND and excludes the voided payment.
- Representative visual compositions were rendered and inspected at 1100 px, 728 px, 390 px, and 320 px. The reference's screen picker was exercised across all six views at 320 px; each selected exactly one view and reported no horizontal overflow. This checks the documentation reference, not the future product's responsive or accessibility compliance.
- Specified text/background contrast pairs measured between 6.08:1 and 15.00:1, including primary-action and status colors; the control boundary against white measured 3.58:1. Full interactive accessibility testing remains pending.
- Changes are documentation, a static-screen design reference, the AI usage note, and a README entry. No application source was changed.

Pending outside this documentation task: team approval of O01–O10; user/proxy research; detailed production screen iteration; application implementation; runtime, accessibility, concurrency, load, and usability verification. No human agreement or passing implementation test is inferred from this document.
