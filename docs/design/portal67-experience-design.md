# Portal67 experience design

Version 1.0 · 25 October 2026 · Design proposal for team review

**Portal67 should help a student leave with a definite answer: which classes they are registered for, where they need to be, what their results mean, and what happened to their request.** Lecturers and administrative staff need the corresponding tools to maintain those answers accurately.

The proposed character is a quiet academic workspace: clear course identities, carefully aligned schedules and records, warm neutral surfaces, and a restrained ink-blue accent. Its distinction comes from bringing related decisions together. Registration includes schedule comparison; GPA includes its calculation basis; requests retain the decision and explanation. Decoration does not substitute for any of those capabilities.

This is a product and interaction design specification, not application code or an approved SRS. The briefs establish the scope; user behavior beyond the briefs is a hypothesis to validate. Existing source was inspected for compatibility, not treated as proof of tested behavior. No interviews, usability sessions, load tests, or official HCMUS policy verification were performed for this document.

## Reading map

| Document | What it settles |
| --- | --- |
| **This document** | Evidence, scope, users, design reasoning, navigation, visual direction, and unresolved decisions |
| [Page and flow specifications](portal67-pages-and-flows.md) | Every MVP page, content order, actions, state transitions, responsive behavior, and page-level acceptance |
| [Component and content specification](portal67-components.md) | Reusable component anatomy, interaction contracts, visual tokens, accessibility, and exact example messages |
| [Coverage and validation plan](portal67-design-validation.md) | F01–F08 traceability, source-to-design gaps, usability tasks, design acceptance, and delivery priorities |
| [Visual design reference](portal67-design-reference.html) | Six annotated compositions: registration, Today, results, lecturer gradebook, staff request review, and tuition; responsive reference with illustrative records |

Page IDs (`S02`), component IDs (`C05`), decision IDs (`D03`), and acceptance IDs (`A04`) connect the documents. Screen names and paths below are proposed destinations, not existing frontend routes.

The visual reference is a documentation companion. Its screen picker changes the example being inspected; controls inside each depicted product screen are static. It is not a functional application prototype.

## 1. Evidence and boundaries

### 1.1 Source register

Page references give both the PDF viewer page and the printed page because the files include covers and contents.

| ID | Source | Design-relevant evidence |
| --- | --- | --- |
| B0 | [PA0](../../pa/PA0/T02-PA0.pdf), §2.1–2.4, PDF pp. 5–6 / printed pp. 3–4 | HCMUS context; fragmented GPA/timetable tools; registration overload; students, teachers, and administrative workers; five must-have feature areas |
| B1 | [PA1](../../pa/PA1/T02-PA1.pdf), §2.1, PDF pp. 4–6 / printed pp. 2–4 | Desktop/laptop browser prototype; sample data; no official-system connection, real payment processing, or legally valid documents; three-week development period |
| B2 | PA1, §2.3–2.4, PDF pp. 7–8 / printed pp. 5–6 | F01 registration, F02 tuition, F03 results, F04 scheduling, F06 requests, F07 role permissions, F08 validation are Must; F05 resources is Could |
| B3 | PA1, §4.1–4.3, PDF pp. 10–12 / printed pp. 8–10 | Requirements precede reviewed design; later changes to baselines must be documented; limited design and implementation effort |
| I1 | [Academic model](../../src/backend/internal/academics/model.go), [service](../../src/backend/internal/academics/service.go), [sample records](../../src/backend/db/seeds/00002_demo_academics.sql) | Course versus section, multiple meetings, capacity, course kinds, registration conflict rules, grade entry, estimated GPA |
| I2 | [Document model](../../src/backend/internal/documents/model.go), [service](../../src/backend/internal/documents/service.go) | Three request types; pending/approved/rejected; mandatory purpose and decision comment; student-owned records and history |
| I3 | [Tuition model](../../src/backend/internal/tuition/model.go), [service](../../src/backend/internal/tuition/service.go) | Semester charges, recorded/voided payments, five balance states, corrections; no transfer workflow |
| I4 | [Routes](../../src/backend/internal/app/routes.go), [identity](../../src/backend/internal/auth/user.go), [login](../../src/backend/internal/auth/handler.go) | Student/lecturer/staff roles, assigned-section access, email/password login; no account provisioning or password recovery interface |
| G1 | [WCAG 2.2](https://www.w3.org/TR/WCAG22/) | Accessibility reference for contrast, keyboard operation, focus, target size, reflow, errors, and status feedback |
| G2 | GOV.UK [check answers](https://design-system.service.gov.uk/patterns/check-answers/), [confirmation](https://design-system.service.gov.uk/patterns/confirmation-pages/), [error summary](https://design-system.service.gov.uk/components/error-summary/) | External interaction references for reviewing consequential changes, confirming completion, and locating form errors |

The external references inform interaction quality. They are not evidence that HCMUS users have been observed behaving in a particular way, and they do not dictate Portal67's visual identity.

### 1.2 Requirement language

- **Brief requirement:** explicitly supported by B0–B3; preserve in the MVP.
- **Design decision:** the proposed solution to a requirement or plausible task; review against users.
- **Current contract:** behavior found in I1–I4; may need revision after requirements review.
- **Deferred:** a separate capability that is not necessary to demonstrate the brief's Must features.
- **Open policy:** a rule the briefs do not settle. A visible interface must not invent it.

Unless marked otherwise, layout, ordering, copy, default filters, and interaction choices in this package are design decisions. The page specifications describe the intended experience, including work the frontend has not implemented.

### 1.3 Resolve the brief's contradictions explicitly

| Tension | Working interpretation | Consequence |
| --- | --- | --- |
| PA0's 99.99% uptime expectation versus PA1's explicit prototype limitation | Reliability is a user need; that uptime is not an established service guarantee | Design truthful pending/failure/recovery states. Never display a reliability claim or fake waiting-room position |
| PA1 says “No features are out of scope” while deferring resources and excluding real integration | Use the detailed feature priorities and explicit operating constraints for this design | Record this interpretation for the next PA; it does not silently amend the submitted brief |
| Teachers “manage course content” in stakeholder prose, but F03 specifies grade entry | Assigned sections and grade entry are the MVP lecturer tasks | No assignment editor, learning materials, quizzes, or replacement Moodle inside this design |
| PA0 mentions replacing third-party polls/forms | F06 is a simple document-request workflow | No general form builder, polling system, or campus-wide workflow engine |
| “Timetable planning” could imply a full planner | Comparing a prospective section against actual registrations is the minimum useful planning interaction | No stored multi-plan optimizer or reserved-seat basket in the MVP |

### 1.4 Included, deferred, and deliberately absent

The MVP includes registration and cancellation for courses and specialized topics; weekly classes and exams with overlap explanations; grade entry and understandable estimated GPAs; student tuition views and staff simulated-record maintenance; document submission, decisions, and history; role-aware access and recoverable errors.

Resource booking (F05) remains a later concept: staff would search room/equipment availability, inspect a time slot, request a booking, and track its status. It would need booking ownership, conflict, cancellation, and approval policies before detailed design. Do not put a disabled Resources item in the MVP navigation. A complete MVP design does not imply an approved resource-booking product.

Also deferred: waitlists, transactional section swaps, saved schedule plans, GPA forecasts, degree audits, prerequisite engines, online payment gateways, tuition installments, PDF issuance, bulk imports, notifications/email/SMS, account administration, and official institutional integrations. None should appear as an apparently working control.

## 2. Users and actual jobs

These are behavioral situations, not invented demographic personas. Frequency and device assumptions are hypotheses, not measurements.

| Situation | Trigger and likely behavior | Information needed to act | Design response |
| --- | --- | --- | --- |
| Student choosing sections | Knows a course code or a name; compares time and remaining seats, possibly under time pressure | Course identity, all meetings, section, availability, conflicts, actual registration outcome | Searchable rows; persistent registered summary; one-section schedule preview; no marketing header |
| Student checking the next class | Short visit before or between classes; may use a phone | Day, start/end, room, course, section | Today summary and mobile agenda; no need to decode a miniature weekly grid |
| Student preparing for exams | Checks dates and rooms; must spot simultaneous exams | Chronological list and exact overlapping pairs | Exams are a separate Schedule view with explicit conflict links |
| Student interpreting results | Wants to understand a number and whether it is complete | Course scores, credits, grade points, excluded ungraded courses, sample rule | GPA number beside its coverage; calculation disclosure; no ranking or celebratory score widget |
| Student checking tuition | Wants to know whether a recorded payment changed a balance | Selected semester, charge total, effective payments, balance, record details | A reconciled ledger and clear status; no Pay now button |
| Student requesting a document | Infrequent task; does not necessarily know institutional terminology | Plain-language type choices, purpose, confirmation, decision explanation | Small single-page form, durable detail page, visible history |
| Lecturer entering scores | Repetitive work across a roster; interruptions and typing errors are plausible | Stable student identity, saved values versus edits, errors, save result | Keyboard-friendly table, changed-row count, explicit atomic save |
| Staff reviewing requests | Works through a queue; needs context before deciding | Type, student identity, purpose, age, previous decision/history | Oldest-pending-first queue, detail review, required comment |
| Staff correcting tuition | Starts with a named student/ID; a wrong-person error is costly | Full identity, semester, before/after amounts, recorded status | Search before detail; identity retained through form and review; explicit voiding |

Undergraduates and postgraduates share the shell. Do not invent different registration or GPA rules for them without data and policy. Vietnamese names must remain readable with diacritics; student IDs disambiguate similar names. Study-program segregation, if required later, is a permissions/data question rather than a decorative user preference.

## 3. Decisions that make this Portal67

| ID | Decision | User behavior / evidence | Tradeoff and validation |
| --- | --- | --- | --- |
| D01 | Role-specific starting pages: student Today, lecturer My sections, staff Request queue | B2; users arrive to do different work | No common analytics dashboard. Test that each role can identify its first task immediately |
| D02 | Registration is a comparison workspace with an actual-registration summary | B0's planning problem; I1's sections and meetings | More structured than a card gallery; avoids repeated page changes. Test course-versus-section comprehension |
| D03 | Preview one section against the existing timetable; label it “Preview only” | Choosing a time is a real registration decision | No persistence or seat reservation. Test whether anyone mistakes preview for registration |
| D04 | Show success only after a confirmed saved outcome; unknown outcomes get reconciliation | B0 overload concern and F08 | Less superficially instant than optimistic enrollment; protects trust under slow responses |
| D05 | Scope semester selection to the academic or tuition task, with the semester repeated in confirmations | Registration is term-bound; requests often are not | Slight repetition prevents cross-semester edits. Never hide term changes in a universal filter |
| D06 | Use course/section rows, a timetable, ledgers, and a roster where each fits the data | Comparing seats differs from entering marks or reading a balance | More than one page pattern, but shared typography/status/feedback keep them coherent |
| D07 | GPA includes coverage and the sample calculation basis | B0's inadequate GPA tracking and B2's estimated rules | No speculative forecasts. Test if students can explain why an ungraded course did not reduce GPA |
| D08 | One durable request record holds both the student's progress and the staff decision | F06's status and history | No disconnected confirmation email dependency. “Approved” never implies a file was issued |
| D09 | Important exceptions appear near the affected object; a home summary only links to them | Users need an actionable reason, not a notification count | No global bell without notification infrastructure; no red badge on every module |
| D10 | Preserve a useful amount of information on screen, using disclosure for explanation rather than for essential facts | Need for clean, low-overload UI | Course time, seat state, errors, and document decision remain visible; rules and older history can expand |
| D11 | Restrained academic visual language, with emphasis reserved for decisions | User explicitly requests fresh, clean design without generic gradients | Originality comes from composition and precise domain detail, not a visual effect |
| D12 | Distinguish absence, zero, provisional state, and failed loading | Unreleased grades and unrecorded charges are not zero-valued facts | More deliberate empty-state copy; avoids misleading GPA, payment, and timetable claims |

## 4. Navigation and page architecture

### 4.1 Shared shell

At desktop widths, use a 208 px text-led left navigation, a compact identity area, and one main content region. Place the Portal67 wordmark and a quiet “Course prototype” label above navigation. The footer/account area contains the user's name, role, Account & help, and Sign out. Avoid decorative avatars, global search, and an unrelated activity feed.

Each page begins with a breadcrumb only when it is below a list, a plain title, and a local semester selector where relevant. The page title does not repeat a welcome message. A contextual error occupies a single band below the title. A stable navigation position matters more than reducing the sidebar to unexplained icons.

| Role | Primary navigation, in order | Default entry | Secondary destinations |
| --- | --- | --- | --- |
| Student | Today; Registration; Schedule; Results; Tuition; Documents | S01 Today | Section detail, request form/detail, Account & help |
| Lecturer | My sections | L01 My sections | L02 section gradebook, Account & help |
| Staff | Document requests; Tuition records | A01 Request queue | A02 decision detail, A04 student tuition, record forms, Account & help |

There is no role selector in the ordinary app: I4 gives each account one role. Separate demo accounts can be used for evaluation. A future multi-role account requires a real permission model before adding a switcher.

### 4.2 Destination inventory

| ID | Destination | Proposed location |
| --- | --- | --- |
| U01 | Sign in | `/sign-in` |
| U02 | Account & help | `/account` |
| U03 | Access, missing page, and service-recovery states | Contextual destinations / fallback pages |
| S01 | Today | `/student` |
| S02 | Registration | `/student/registration` |
| S03 | Section detail / schedule preview | `/student/registration/sections/:section` |
| S04 | Schedule: Classes | `/student/schedule/classes` |
| S05 | Schedule: Exams | `/student/schedule/exams` |
| S06 | Results | `/student/results` |
| S07 | Tuition | `/student/tuition` |
| S08 | My document requests | `/student/documents` |
| S09 | New document request | `/student/documents/new` |
| S10 | Document request detail / receipt | `/student/documents/:request` |
| L01 | My sections | `/lecturer/sections` |
| L02 | Section gradebook | `/lecturer/sections/:section/grades` |
| A01 | Document request queue | `/staff/documents` |
| A02 | Review request | `/staff/documents/:request` |
| A03 | Student lookup for tuition | `/staff/tuition` |
| A04 | Student tuition record | `/staff/tuition/students/:student` |
| A05 | Record / correct / void simulated payment | Nested under A04 |
| A06 | Add semester charge | Nested under A04 |

S03 is a wide side panel when there is room and a full page otherwise. It retains a shareable destination and Back behavior. A05 and A06 use focused pages because staff must inspect identity and values without a cramped modal.

### 4.3 Context and Back behavior

Academic pages remember a selected term within the current account session. A specific deep link overrides remembered selection. Default to the term containing today's campus date; if multiple terms qualify, use a deterministic latest-start default and show its full label. With none current, suggest the nearest upcoming term; otherwise the latest past term with a “Past semester” label. Registration still obeys actual rules; term selection is not evidence that registration is open.

Results defaults to the latest term containing the student's enrollments; its cumulative summary explicitly includes all available terms. Tuition defaults to the current term, allowing a legitimate “No records for this semester” state. A staff lookup starts with no selected student; changing student clears the previous student's forms.

Back from detail restores list filters, term, query, and scroll. A local search query belongs to its page and must not filter other modules. A successful mutation updates the affected registration/timetable/results/tuition/request views before claiming that they agree. Failed refreshes are labeled; an old summary cannot silently masquerade as current.

## 5. Visual direction: an academic work surface

### 5.1 Composition

Use a light warm canvas, white working surfaces, dark ink text, hairline separators, and a single deep blue accent. Favor section headings and aligned rows over repeated floating cards. A small vertical active-navigation mark, generous title spacing, tabular numerals, and course-code typography form a recognizable rhythm.

Do not use gradients, glass effects, floating blobs, generic stock campus photos, oversized greetings, ring charts for single balances, pill-shaped containers around every fact, or full-color backgrounds behind every module. A calm home page must not hide density needed for comparing sections or entering grades.

### 5.2 Proposed visual values

These are design values, not framework choices. Detailed component use is in the component specification.

| Property | Proposal | Reason |
| --- | --- | --- |
| Canvas / surface | `#F7F8FA` / `#FFFFFF` | Slight separation without heavy panels |
| Main text / secondary text | `#18232D` / `#52606D` | Strong reading contrast, including table metadata |
| Accent / accent tint | `#234F70` / `#EAF1F6` | Academic ink character; limited to primary actions, selection, and links |
| Separator / control boundary | `#DCE2E7` / `#7C8995` | Quiet row rules; stronger visible boundaries for inputs |
| Success | Text `#225C47`, surface `#EDF5F0` | Confirmed completion, always accompanied by a label |
| Caution | Text `#7A4A0B`, surface `#FFF4DF` | Conflicts and information needing attention |
| Error | Text `#A12D32`, surface `#FCEFF0` | Blocked or failed actions, not general emphasis |
| Typeface | System sans with complete Vietnamese support; Segoe UI on Windows, equivalent system fallback elsewhere | Reliable rendering, low asset cost, familiar readable numerals |
| Type sizes | Title 28/34; section 20/28; body 16/24; dense table 14/20; supporting 13/18 px | Meaningful hierarchy without oversized dashboard metrics |
| Spacing | 4, 8, 12, 16, 24, 32, 48 px | Tight within a record, larger between different tasks |
| Corners | 6 px controls; 8 px contained panels; square table regions | Refined but useful, without excessive rounded cards |
| Motion | 120–180 ms feedback transitions; no essential animated movement | Confirm interaction without delaying repeated work |

Show student names in normal case with full diacritics. Use tabular numerals for seats, scores, credits, and VND. A monospace treatment may distinguish course codes, but not entire tables. Do not truncate the distinguishing part of a course or student name; wrap it.

### 5.3 Density and screen adaptation

Desktop/laptop is the primary environment in B1. At about 1280 px and above, registration can use a main comparison region and a 288–320 px supporting panel. Around 1024 px, collapse that panel below the table and keep a compact registered-count summary. At narrower widths, use a labeled Menu button and a drawer; preserve every task in a single column.

Phone support is a design extension for plausible quick schedule/status checks, not a claim that the brief required a mobile-first product. Class schedules become day agendas; tuition rows become labeled records; registration groups become vertically stacked sections. Grade entry keeps student identity and one editable score visible together, with remaining details expandable. Never shrink a desktop screenshot to fit a phone.

Use 24–32 px content padding on desktop and 16 px on narrow screens. Form text stays around 60–70 characters per line. Wide data work can use the remaining viewport; a narrow form should not stretch to 1500 px. At 200% zoom and at 320 CSS px, primary actions and error messages remain reachable.

### 5.4 Information budget

Student Today contains the next class/day agenda, up to three actionable exceptions, and one compact link to registration. It is not a synopsis of every database table. Registration shows identity, all meeting times, seats, conflict/registration state, and an action; lecturer biographies and long rules are secondary. Request detail prioritizes the current decision and comment before history. Tuition shows the reconciliation before individual lines. Results shows completeness before trends.

These budgets are content priorities, not arbitrary limits on essential information. All affected items remain available through “View all,” pagination, or a destination link. Never conceal an error to maintain a tidy visual count.

## 6. Policies still needing team decisions

Proposed defaults below allow concrete design review. They are not official HCMUS rules or approval claims. The owner is the relevant team role from PA1, not a statement that a named person has accepted the work.

| ID | Decision to settle | Proposed working position | Owner / effect |
| --- | --- | --- | --- |
| O01 | Registration opening/closing dates, prerequisites, credit caps, program restrictions, and retakes | Implement only demonstrated seat, duplicate-course, overlap, and ended-term checks; show no invented deadline or eligibility badge | BA; changes S02/S03 and error states |
| O02 | Exact grading rules, repeat-course treatment, historical rule versions | Display current sample bands and credit-weighted calculation; all graded records count unless a policy changes | BA + lecturer representative; changes S06/L02 |
| O03 | Vietnamese versus English UI | Recommend Vietnamese-first labels with preserved official course names; this package uses English specification prose and provides Vietnamese core vocabulary | Designer + team lead; confirm with likely evaluators/users |
| O04 | Request catalog, supporting evidence, processing time, revision/withdrawal, issuance | Three current types and terminal decisions; no attachments, promised turnaround, withdrawal, or issuance | BA + administrative representative; changes S09/S10/A02 |
| O05 | Meaning and timing of tuition charges, due dates, corrections and audit needs | Charges are explicit records; registration does not automatically create a charge; payment corrections affect totals; no overdue label without a due date | BA + administrative representative; changes S07/A04–A06 |
| O06 | Grade publication and concurrent edits | Explicit Save makes demo grades visible; no draft/publish workflow; propose same-section edit-conflict detection before multi-user use | BA + developer; affects L02 feedback |
| O07 | Account provisioning, recovery, staff permission split | Provided demo accounts; no self-sign-up; one staff role for both duties, per current contract | Team lead; affects U01/U02 and role matrix |
| O08 | Help/contact ownership and operating hours | In-app help explains the prototype; add only a verified project contact | Team lead; never invent an office, email, SLA, or ticket service |
| O09 | How frontend obtains GPA rules | Add a read-only rule description or verified shared sample-rule contract | Developer + BA; essential to truthful S06 disclosure |
| O10 | Retry reconciliation for newly created payments and requests | Show unknown outcome, inspect latest records, and prevent immediate blind resubmission; require a durable operation reference for robust duplicate prevention | Developer; required before claiming safe automatic retries |

## 7. Delivery and review status

This package is intended to feed the PA2 requirement decisions and PA3 UI design review. The recommended implementation order is registration with its timetable comparison, then results with grade entry, then the complete document loop and tuition maintenance; shared identity/error states are required throughout. The validation document explains which refinements may be cut without cutting the briefs' Must features.

**Definition of a useful review:** a teammate can follow every Must feature from its brief ID to a page, action, saved outcome, failure state, and test scenario. A visually attractive screen without its error and completion states is unfinished design.

Authorship: AI-assisted draft prepared at the user's request from the two briefs, the repository contracts, and the linked accessibility/form references. Team review and behavior validation are pending. See the [AI usage note](../ai-usage/2026-10-25-experience-design.md).
