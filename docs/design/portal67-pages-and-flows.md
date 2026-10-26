# Portal67 pages and flows

Companion to the [experience design](portal67-experience-design.md). All screens are proposed design, not a claim that the frontend exists. The named sample courses come from the repository's academic seed records; other example counts, people, and amounts are illustrative unless identified as a calculation fixture.

## 1. End-to-end journeys

### J01 — Choose a section with confidence

**Student intent:** “I need this course, but it has to fit the classes I already registered for.”

1. Open Registration directly or from Today. Confirm the semester beside the title.
2. Search `CSC10003`, its name, or select Specialized topics. Keep the typed query when results refresh.
3. Read course-group rows: section, lecturer, every meeting, seats, and status. Full sections remain discoverable so absence is not mistaken for an unoffered course.
4. Open a section's schedule preview. Existing registrations are solid blocks; the proposed section has a dashed boundary and “Preview only.” An accompanying text list identifies exact clashes.
5. Select **Register for section 01** only after the section identity and all meetings are visible. This is a direct explicit action, without an additional generic “Are you sure?” step.
6. While saving, preserve the row and show “Registering…”; disable repeated submissions for that section. Seat counts elsewhere need not freeze the entire page.
7. On confirmed success, replace the action with Registered, add the section to My registrations, update credit count and schedule, and show an inline receipt with section and semester. The student can continue browsing.
8. If the last seat was taken, show the full state and offer **View other sections** while preserving query and term. If a conflict blocks registration, name both classes and times. If the result is unknown, reconcile the student's registrations before offering a new submission.

**Why:** B0 identifies registration overload and timetable planning as the problems. A fake instant success or an isolated catalogue page would leave both unresolved. **Completion evidence:** a confirmed active enrollment, not a clicked button or green toast.

### J02 — Cancel an existing registration

Open My registrations → choose the specific section → review course, section, semester, and all meetings → read “Your seat will be released. Registering again depends on availability.” → **Cancel registration** → wait for confirmed outcome → retain a Cancelled history entry and update the timetable.

The dialog's safe action is **Keep registration**. Graded records and ended terms cannot be cancelled under I1. Explain these states next to the action. Do not call cancellation an “Undo,” promise seat restoration, or suggest a safe section swap when the system cannot perform one atomically. No tuition refund is implied.

### J03 — Check a class or exam

Today → next class with time and room → full Classes schedule if needed. For exams: Schedule → Exams → chronological entries → select conflict explanation to inspect both exams. Example from the seed: CSC13002 on 5 Jan 2027, 08:00–10:00 overlaps MTH00003, 09:00–10:30. The interface identifies the collision; it does not let the student edit an institutional exam or claim it has reported the issue.

### J04 — Enter grades, then understand the GPA

Lecturer My sections → assigned section → gradebook → edit scores → inspect changed rows and validation → **Save 3 grades** → confirm saved values and time. Student Results → select semester → inspect course scores and graded-credit coverage → expand calculation basis.

Saving demo grades makes them visible under the current contract. There is no fictional Publish step. Blank untouched rows stay ungraded, zero is a valid entered score, and removing an existing score is unsupported until a deletion policy exists. An invalid row prevents the entire submitted batch from saving.

### J05 — Request a document and understand the decision

Student Documents → New request → choose a plain-language type and enter purpose → inspect prefilled identity and short inline summary → **Submit request** → durable detail page with request ID and Pending status. Staff Request queue → oldest pending request → read identity and purpose → select Approve or Reject → write a student-visible comment → review decision → confirm. Student opens the same request detail to see the decision, comment, and event history.

The receipt means the request was saved. Approval means the demo workflow approved it; no legal document or download has been created. A rejected request can be followed by a new request with prefilled type/purpose for revision, but the original remains immutable and the student must submit deliberately.

### J06 — Reconcile tuition and correct a record

Student Tuition → select semester → read charges minus recorded payments → inspect a line item if the balance is unexpected. Staff Tuition records → search student ID/name → confirm identity → select semester → record or correct a simulated payment → review before/after balance → save → inspect confirmed ledger. Voiding retains the row and excludes its amount from paid totals.

Entering a payment is maintaining a sample record, not initiating a money transfer. Registration and tuition are not coupled unless a charge rule is added later. Staff must never infer identity from a name alone.

## 2. Universal screens

### U01 — Sign in

**Job:** enter the correct role workspace with a provided account. **Evidence:** F07, I4. **Components:** C01, C03, C11, C14.

**Content order:** Portal67 wordmark; “Sign in”; a short course-prototype explanation; Email field; Password with Show/Hide; primary **Sign in**; “Need an account?” help explaining that demo accounts are provided by the project team. No public landing-page hero, registration form, social login, or fake SSO.

Labels remain visible after entry. Support password managers, paste, and Enter submission. Initial focus goes to email when appropriate; no keyboard trap. The prototype explanation says “Sample academic records. Not connected to the official HCMUS Portal.” Avoid showing credentials in this document or product UI.

**Outcomes:** successful login goes to the requested allowed deep link or the role default. Wrong credentials show a form-level message; keep email and clear password. An inactive-account response explains access is unavailable without claiming a reset route exists. Service failure says to retry without treating credentials as incorrect. Already signed-in visitors go to their workspace.

**Narrow view:** same short form, maximum about 420 px wide; no split-screen illustration. **Pass condition:** each role reaches its proper start page, and an unauthorized deep link does not reveal another role's data.

### U02 — Account & help

**Job:** verify the signed-in identity and understand prototype limits. **Evidence:** B1, F07. **Components:** C02, C15, C18.

Show name, institutional ID if present, email, and role as read-only text. Omit date of birth because it is unnecessary to these tasks despite being present in the identity model. Missing institutional ID is “Not provided”; never substitute an internal database ID.

Below identity, show concise expandable help for registration conflicts, sample GPA rules, simulated tuition, and document decisions. Provide a verified project contact only when one exists. There is no generic support-chat bubble or assumed ticket system. Sign out is clearly available; successful sign-out clears cached identity-specific screen data and returns to U01. Unsaved edits require a leave decision before sign-out.

**States:** failed identity load offers Retry and Sign out, not an empty invented profile. **Narrow view:** stacked labeled values. **Pass condition:** users can identify their current role and cannot edit unsupported profile fields.

### U03 — Access and recovery pages

| State | Visible content | Recovery |
| --- | --- | --- |
| Not signed in / expired session | “Sign in to continue.” No private data remains visible | Return to the authorized original destination after login; never automatically replay a pending mutation |
| Wrong role | “This page is not available for your account.” | Go to the user's role start page |
| Missing or inaccessible private record | “We couldn't find this record.” | Return to the permitted parent list; do not confirm another student's record exists |
| Whole service unavailable | “Portal67 is temporarily unavailable. Your last action may need checking.” Use the second sentence only for uncertain mutations | Retry a read or return to sign-in; no invented outage ETA |
| One module unavailable | A local failure region replacing only the failed module | Retry that module; other navigation remains usable |
| Unexpected error after a write | “We couldn't confirm whether this was saved.” | Check current record state; follow the C13 reconciliation rule |

An expired editing session obscures private data and holds the in-memory draft only for reauthentication in the same tab by the same account. A different account, sign-out, or tab closure discards it. No promise of draft survival across reloads or browser restarts. **Pass condition:** user-entered data is not exposed to a subsequent account and no uncertain mutation is silently repeated.

## 3. Student screens

### S01 — Today

**Job:** answer “Where next?” and surface exceptions that warrant a visit elsewhere. **Evidence:** F04; D01, D09, D10. **Components:** C02, C04, C07, C12, C14.

**Composition, top to bottom:**

1. Today title, campus date, and the term used for the class agenda.
2. Next scheduled class: time range, course name/code, section, room, and lecturer. Below, remaining classes for the day in chronological order.
3. A compact Needs attention list, up to three direct links, only when supported by actual data: exam overlap first, then class overlap, then a tuition balance that needs inspection. Positive balances are neutral “Balance to review,” not overdue warnings without due dates.
4. One line linking to Registration with registered section/credit count. An upcoming exam may be shown as a dated line after today's classes, clearly distinguished from a class.

No GPA, total-payments KPI, news carousel, weather, inspirational message, or duplicated full request queue. A request's “recently changed” state requires a real read/change marker; the current design does not fabricate unread notifications. Documents remains one navigation choice away.

**States:** no class today says so and links to the weekly schedule; outside a term shows “No current semester schedule”; no enrollments links to Registration. Missing room says “Room not provided.” A failed schedule read says “Couldn't load today's classes,” never “No classes.” Partial module failures stay separate.

**Time truth:** recurring meetings are derived from weekday and term dates. Current records do not include holiday cancellations or make-up sessions; label the agenda “From your weekly timetable” rather than promise live attendance information. Campus time is Asia/Ho_Chi_Minh, UTC+7.

**Narrow view:** next class and same-day rows first; exceptions follow, not a horizontal card carousel. **Pass condition:** a student can identify the next scheduled class and room without opening a calendar grid.

### S02 — Registration

**Job:** find and register a suitable offered section, then verify actual registrations. **Evidence:** F01, F04, F08; D02–D05. **Components:** C04–C08, C12–C14, C16.

**Desktop wireframe:**

| Main region, about two-thirds | Supporting region, about one-third |
| --- | --- |
| Registration + selected semester | My registrations: count and credits |
| Search by code/name; All / Courses / Specialized topics | Compact registered section rows and Cancel actions |
| Results count; optional Available seats filter | “View timetable” link |
| Course groups with comparable section rows | Selected section preview, when open |
| Pagination if needed | Clearly labeled preview summary; no saved-plan claims |

**Catalogue anatomy:** course code and full name are the group heading; credits and kind appear once per course. Each section shows section code, lecturer, all weekday/time pairs with rooms, available seats out of capacity, a text state, and **View section**. The expanded section view offers **Register for section …**. This makes a deliberate click inexpensive without placing ambiguous identical Register buttons beside hidden meeting information.

Default sort is course code then section code. Search matches course code/name; lecturer search is not promised by the current filter contract. All kinds are visible by default. Full sections remain listed with Full; an optional “Has seats” filter is user-controlled. Avoid a filter drawer containing unsupported departments, prerequisite chains, degree levels, or campus locations.

My registrations lists active sections with course, section, credits, and a **Cancel** action. Cancellation history sits in a collapsed “Cancelled registrations” region. The active count excludes cancelled records. Display past semesters read-only with a reason. No artificial countdown or “Registration open” banner is inferred from a term's dates.

**Refresh:** show “Availability checked at 09:41” only after a successful fetch, plus **Refresh availability**. The time describes the fetch, not guaranteed real-time availability. Refresh on return from a mutation and on an explicit user refresh; do not reorder visible rows in the middle of comparison. The server still decides whether a seat exists.

**Edge states:**

| Condition | Treatment |
| --- | --- |
| No matching search | Preserve query; show Clear filters; do not imply no courses exist |
| Term has no offerings | “No sections are available for this semester.” Keep term control |
| Full section | Full text replaces actionable registration; detail still readable |
| Same course in another section | “Already registered in section 01.” Link to that registration; no implied swap |
| Schedule conflict | Explain course, section, day, overlap; allow inspecting both schedules |
| No meeting information | “Meeting times not provided. Schedule compatibility cannot be confirmed.” Do not mark Conflict-free |
| Stale result after submission | Re-fetch and update only affected objects; preserve comparison context |
| Registration save fails clearly | Keep selected section and explain specific failure |
| Save outcome unknown | Pending-verification state; inspect actual enrollment before a repeat |

**Narrow view:** semester and search first, collapsible My registrations summary immediately after, then course groups. Each section is a labeled vertical record; meeting times do not disappear. Preview opens S03 as a full page. **Pass condition:** successful registration and cancellation update both the active list and timetable; blocked attempts never appear as active registrations.

### S03 — Section detail and schedule preview

**Job:** evaluate a particular section before committing. **Evidence:** F01/F04; D03. **Components:** C05–C07, C12, C16.

Show breadcrumb to filtered Registration; course name/code and kind; section and credits; semester; lecturer; seats; complete meeting list. Beneath it, place **Preview against my timetable** and a compact day/time comparison. The side-panel version uses a text schedule comparison if a weekly grid would be too narrow. Never compress seven days into an unreadable 300 px panel.

Existing meetings are identified by course/section. The candidate has dashed outlines plus Preview only text. A conflict callout names the overlap precisely: “Monday 09:30–11:30 overlaps CSC13002 · section 01, 07:30–11:30.” It links to the existing registration. Comparing a candidate does not cancel, move, or reserve anything.

Primary action is **Register for section 01**, followed by Back to results. Known blocking reasons replace the active primary action with an explanation; no tooltip-only disabled button. A future exam clash, if exam data is present, is a warning distinct from the weekly-class blocker because current registration rules do not reject exam overlaps. Do not label a section safe simply because exams have not been published.

**States:** Registered, full, term ended, duplicate course, missing meetings, partial data, loading, unknown save outcome. If a section disappears, preserve Back to results and explain it is unavailable. **Narrow view:** meeting list and conflict explanation precede any diagram; no horizontal dependency for registration. **Pass condition:** users can distinguish Preview only from Registered and identify every meeting affected by registration.

### S04 — Schedule: Classes

**Job:** see the student's committed weekly schedule and locate overlaps. **Evidence:** F04. **Components:** C04, C07, C08, C12, C17.

Header: Schedule; Classes / Exams view navigation; semester; Week grid / Agenda toggle. Default desktop view is a typical weekly timetable, not a dated calendar implying holiday knowledge. A note reads “Repeats during [term dates].” Default mobile view is Agenda, grouped Monday through Sunday.

Grid shows an explicit time axis and all days with meetings, including Saturday and Sunday. Empty weekend columns may collapse but an existing weekend meeting never does. Each block contains course code, section, time range, and room; name is visible in the expanded detail and the agenda. Match the detail's room and lecturer to the selected meeting, not an arbitrary first meeting of the course.

When sessions overlap, draw adjacent lanes with an outlined conflict marker and also list the pair in text above the grid. Do not lay one card over another. Block colors are secondary identifiers; text and position carry meaning. Start/end boundaries that touch are not overlaps.

Selecting a meeting opens a read-only detail with full course, lecturer, room, term, and links to the relevant registration. There is no drag-to-reschedule behavior. Semester dates bound the recurrence; an overlapping time in disjoint terms is not a conflict.

**States:** no active enrollment; no scheduled meetings despite enrollment; missing room; seeded overlap; error; stale data. “No timetable entries” must not assert “No registrations.” **Pass condition:** the desktop grid and text agenda contain the same meetings, and any overlap is understandable without color.

### S05 — Schedule: Exams

**Job:** find an exam and notice collisions before relying on the schedule. **Evidence:** F04. **Components:** C04, C08, C12, C17.

Use a chronological list grouped by date: weekday/date, start–end, course/code, section, exam label, and room. Sorting is date/time ascending within the selected term. Default show upcoming entries with an accessible Show past exams disclosure, but always provide a clear All exams view and count. All entries remain available when the term is past.

When overlap exists, show an explanation band with the number of affected pairs and links to both records. Pair grouping avoids presenting two duplicate error messages as two different incidents. Example: “CSC13002 and MTH00003 overlap from 09:00 to 10:00 on 5 Jan 2027.” Actual times are always visible, not only “In 2 days.”

No edit/reschedule/report-success actions. A plain help explanation can say the prototype displays seeded conflicts; a real contact link requires O08. Missing room is “Room not provided.” No exam records means “No exam entries are available for this semester,” not a promise that there will be no exams.

**Narrow view:** date groups and stacked exam rows; keep time and room near the title. **Pass condition:** a student identifies both overlapping exams and the exact overlap interval.

### S06 — Results

**Job:** understand course results and how much evidence supports the estimated GPA. **Evidence:** F03; D07, D12. **Components:** C04, C08, C09, C15, C17.

**Content order:** Results title; “Estimated from sample rules” label; cumulative GPA with total graded credits; selected semester and semester GPA with graded-credit coverage; course-results table; How this is calculated disclosure. Avoid an oversized GPA that displaces its basis.

The table contains course code/name, section, credits, score out of 10, letter, and points. Missing score is **Not graded**, not 0, a dash without explanation, or Failed. An actual zero score remains 0.00 and contributes its credits. Fail status can be text on the row without turning the entire table red. No rank, percentile, progress-to-degree gauge, or target-GPA prediction.

**Calculation content:** GPA = sum of grade points × course credits, divided by graded credits; ungraded courses are excluded. Show term and cumulative scopes explicitly. Round only the displayed aggregate to two decimals. The current sample bands are F [0,4), D [4,5), C [5,6.5), B [6.5,8), A [8,9), A+ [9,10], mapping to 0, 1, 2, 3, 3.5, 4 points. These are repository sample rules, not a statement of official HCMUS grading.

**Worked seed fixture:** CSC10001: 8.50 → 3.5 points × 4 credits; CSC10002: 6.75 → 3.0 × 4; MTH00003: 3.00 → 0 × 3. GPA = (14 + 12 + 0) / 11 = **2.36**. An ungraded retake in the next term does not remove the earlier grade under the current contract. Do not silently apply a best-attempt or latest-attempt policy.

Rules are configurable in data, but their readable bands are not currently exposed by a dedicated endpoint. O09 must be resolved so the explanation matches the calculation. No rules-settings page is implied.

**States:** no grades means “Not available yet — 0 graded credits”; no enrollments means no course records; partial grades show e.g. “8 of 11 enrolled credits graded”; failed fetch shows an error, not zero GPA. **Narrow view:** a compact summary followed by course records; all scores and credits retain explicit labels. **Pass condition:** students can reproduce the 2.36 example and explain what happens to an ungraded course.

### S07 — Tuition

**Job:** understand the balance of sample records for a semester. **Evidence:** F02; D05, D12. **Components:** C04, C08, C10, C15, C17.

Title and selected semester precede a concise statement: “Simulated payment records. No money is transferred here.” The reconciliation line shows Charges − Recorded payments = Balance. The balance is visually primary; charges and payments remain readable rather than being disconnected KPI cards.

Then show separate Charges and Payment records sections on the same page. Charge rows: description, recorded date, amount. Payment rows: paid date/time, record reference, amount, Recorded/Voided status; expand for note and updated time. Keep voided rows visible with a textual explanation and exclude them from paid totals. No card fields, bank details, QR transfer code, Pay now, or official receipt download.

| Accounting state | Student-facing treatment |
| --- | --- |
| No charges and no payments | “No tuition records for this semester.” Not Paid |
| Positive balance, no effective payments | “Unpaid” and remaining amount |
| Positive balance, some effective payments | “Partly paid” and remaining amount |
| Zero balance with records | “Paid” and 0 VND remaining |
| Negative balance | “Credit balance” and a positive credit amount; do not show negative debt or promise a refund |

Display integer VND with locale grouping and a visible VND suffix. Do not sum different semesters without the explicit All semesters view. That view must show per-term balances, because a net credit in one term must not hide debt in another. No overdue badge without a due-date field and rule.

**States:** individual record missing note, voided payment, no records, error, stale refresh after staff change. **Narrow view:** reconciliation stacked vertically, then labeled records. **Pass condition:** visible effective payments and charges reconcile exactly with the displayed selected-term balance.

### S08 — My document requests

**Job:** find an existing request before submitting another. **Evidence:** F06. **Components:** C08, C11, C17.

Title Documents and primary **New request**. Default list contains all own requests, newest submitted first, with optional Pending / Approved / Rejected filters. Columns: type, request reference, submitted date, status, Open. Status is always text; dates use absolute values with optional relative supplements. Do not display another student's identity columns on a list already scoped to self.

A short description explains the three supported request types. An empty list explains what can be requested and offers New request. A filtered empty result offers Clear filter. Query errors preserve the controls and show Retry. New request is not automatically blocked by an existing pending request because duplicate policy is undefined; S09 warns contextually instead.

**Narrow view:** each request becomes a compact type/reference/status record. **Pass condition:** a student can find a previously submitted request and reach its history without starting another form.

### S09 — New document request

**Job:** make a clear, valid request with minimal repeated input. **Evidence:** F06/F08, I2. **Components:** C03, C11, C14, C15.

Read-only identity strip: name and institutional ID. Fields in order:

| Field | Presentation | Validation / help |
| --- | --- | --- |
| Request type | Radio group: Enrollment confirmation; Transcript; Other document | Required; one sentence explaining each type. Three choices do not need a searchable dropdown |
| Purpose | Textarea with persistent hint and character count | Required after trimming; maximum 1,000 characters; for Other, ask the user to name the document in this field |

Example hint: “Explain what you need and why, such as a transcript for a scholarship application.” Do not demand attachments, address, phone, copy count, or delivery date when the workflow does not use them. A visible note explains that this prototype records decisions and does not issue official documents.

Below the fields, present a compact inline review of identity and type; **Submit request** and **Cancel**. A separate multi-step wizard is unnecessary for two inputs. If a same-type pending request is known, show its reference and a link: “You already have a pending transcript request.” Allow submission because it may be a legitimate separate purpose.

On invalid submission, keep text, focus the error summary, and link each error to its field. On success, navigate to S10 with a persistent Submitted confirmation and ID. On uncertain outcome, inspect recent own requests; if an exact recent match cannot unambiguously identify the operation, explain the uncertainty and do not silently create a second request.

**Narrow view:** same single column. Leaving with edits asks Keep editing / Discard. **Pass condition:** users submit without retyping identity, and neither a form error nor a network error erases purpose text.

### S10 — Document request detail and receipt

**Job:** see what was requested, its current decision, and the record of changes. **Evidence:** F06. **Components:** C11, C12, C15, C18.

Header: type and request reference; current status; submitted date. The top result panel says Pending, Approved, or Rejected, followed immediately by the staff comment when decided. Pending says “Awaiting staff review”; it does not show an invented queue position or completion date.

Follow with read-only submitted purpose and identity, then History: submission event and decision event with actor's display name/role, timestamp, and comment. The current state is authoritative; history is an explanation, not a separate editable object. Prefer a simple vertical event list over a percentage progress bar that implies more lifecycle stages than exist.

For approved: “Approved in this prototype. No official document has been issued.” For rejected: show the reason without hiding it inside history; offer **Start a new request** with editable prefilled content. No Edit submitted request, Withdraw, Reopen, Download transcript, or Mark collected action under the current contract.

**States:** pending; approval; rejection; unavailable record; history load failure with a still-valid current state if obtained independently. If history cannot be loaded, say so instead of displaying an empty history as fact. **Narrow view:** result, purpose, and timeline stack in that order. **Pass condition:** the student can state both the current decision and its explanation, without interpreting approval as issuance.

## 4. Lecturer screens

### L01 — My sections

**Job:** enter the correct roster quickly. **Evidence:** F03/F07. **Components:** C04, C05, C08.

Default lecturer landing page. Header My sections, semester selector, optional local course-code/name filter. Rows show course/name, section, credits, meetings, enrolled count, **Enter grades**. Only assigned sections are listed. Past sections remain available where the grading contract permits; do not infer a grade-entry deadline from the term end date.

Do not show an invented “Awaiting grades” count unless roster data supports it. Do not fetch every roster merely to decorate this list. No course-creation button or resource upload control.

**States:** no assignments for term; no matches; failed read. **Narrow view:** one labeled section record per row. **Pass condition:** lecturer can choose among sections of the same course without opening the wrong roster.

### L02 — Section gradebook

**Job:** enter and correct demo scores accurately across a roster. **Evidence:** F03/F07/F08, I1. **Components:** C03, C08, C13, C14, C19.

Header includes course, section, semester, lecturer, and enrolled count. A note reads “Saved demo grades are visible to students.” Toolbar has student ID/name search, All students / Not graded filter, and dirty-row count. Do not add a separate Publish button.

Table columns: institutional ID, full name, score out of 10, saved/edited state. Email is secondary detail for disambiguation, not a default wide column. Keep ID and name attached to the input at all widths. Sort predictably by institutional ID, with a stable name fallback; editing or saving must not reorder the roster.

Inputs accept 0–10 with up to two decimal places. Show an existing saved value distinctly from an edited value using a text marker, not color alone. Untouched blank means Not graded. Clearing an existing grade is an unsupported delete, so mark it as an error and offer Restore saved score. Tab moves through inputs in document order; Enter must not accidentally commit the entire gradebook. Paste into an individual score field works; spreadsheet multi-cell import is deferred.

A persistent save bar says “3 unsaved grades” and **Save 3 grades**, plus **Discard changes**. The latter asks confirmation when edits exist. Filtering cannot silently exclude edits from the save scope: show “3 unsaved grades, including 1 hidden by filters” and allow Review changes. Review shows student ID/name and old → new score. Saving existing-grade corrections requires that review; first entries may save directly after validation because all identities are already visible.

**Atomic save:** all valid submitted edits save together; any invalid row leaves all edits unsaved. On success, replace local values with the confirmed response and announce the count. On definite failure, retain inputs. On an unknown outcome, read the roster and compare intended values before retry. “Saved at…” is based on confirmed data or the local confirmation event, clearly labeled; it is not invented audit history.

**Concurrency:** reload a fresh roster before overwriting an existing value and flag observed differences. This reduces mistakes but cannot eliminate the race between read and save. Robust stale-edit rejection requires a record revision contract, currently absent. Do not promise protected concurrent editing until that gap is resolved; include a two-editor test in validation.

**Other states:** empty roster; student no longer enrolled; assignment revoked; validation error; session expiry; another tab's newer value. Before switching semester, section, route, or signing out, offer Save changes / Discard / Keep editing. A failed Save must keep the user in context.

**Narrow view:** one student record with ID, name, score field, and state. Keep a reachable save summary that does not cover the focused input. **Pass condition:** a lecturer enters a zero correctly, catches an 11.5 error, saves only after correction, and can tell whether each changed value is saved.

## 5. Administrative staff screens

### A01 — Document request queue

**Job:** choose the next request that needs a decision. **Evidence:** F06/F07. **Components:** C08, C11, C17.

Default staff landing page: Document requests; Pending / All / Approved / Rejected filters. Pending is the default, ordered by submitted time oldest first, then request ID for ties. Show type, student name + ID, submitted date/time, status, reference, and **Review**. “Waiting 2 days” may supplement, not replace, the timestamp. No Overdue marker without an agreed processing deadline.

Search by request reference or student ID/name is proposed UI behavior; current server status filtering can be combined with local search only while the entire small demo list is loaded. If data is paginated later, search must apply across the whole permitted dataset, not just the visible page. Show the filter scope and result count.

No bulk approval: purpose and comment matter per request. No assignee/team controls without an assignment model. Preserve filters and position when returning from A02. A completed decision removes the item from Pending after confirmation; do not automatically open the next student without an explicit **Review next pending request** action.

**States:** no pending requests (“No requests awaiting review” with All requests); empty search; stale row decided elsewhere; fetch failure. **Narrow view:** student identity remains visible with type and age. **Pass condition:** staff can process in a predictable order and distinguish a filtered empty queue from a failed load.

### A02 — Review document request

**Job:** make an informed, accountable decision. **Evidence:** F06/F08, I2. **Components:** C03, C11, C12, C15, C16.

Breadcrumb to the queue; request reference/type/status; identity strip with name, institutional ID, email; submitted date; full purpose; history. A decision panel follows the purpose, not above it. Decision choice uses Approve / Reject radios with no preselection. Comment is required for both choices, maximum 1,000 characters, with hint “The student will see this explanation.”

Use **Review decision** to show a compact summary of request ID, student, decision, and the exact comment. Then **Confirm approval** or **Confirm rejection**. This additional review is justified because decisions are terminal in the current contract. Cancel returns to the still-filled decision form.

After success, switch the page to read-only decided state with actor, time, and comment. Provide Back to queue and Review next pending request. If another staff member already decided it, retain the unsent comment temporarily for reference, show the actual decision, and disable submission; do not overwrite the earlier decision or append the unsent comment to history.

**States:** pending editable; approved/rejected read-only; missing record; invalid comment; uncertain save result; concurrent decision. No extra states such as Needs information, In progress, Ready for collection unless the workflow is extended. **Narrow view:** identity/purpose precede decision fields. **Pass condition:** decisions cannot be submitted without a comment, and a second reviewer sees the first confirmed decision rather than a false success.

### A03 — Tuition student lookup

**Job:** find the correct person's ledger. **Evidence:** F02/F07, I4. **Components:** C03, C08, C14, C17.

Default Tuition records page is a focused search, not a list of all students' balances. Label “Search by student ID, name, or email.” Show matching students with name, institutional ID, and email to distinguish namesakes. Exact institutional-ID matches appear first as a design preference; confirm supported search semantics before promising it.

The initial state asks for a search rather than exposing the full student population. Search by explicit submission or a short settled-input delay; both keyboard Enter and a labeled Search button work. Do not select the first match automatically. Opening a result leads to A04.

**States:** no query, searching, no matches, many matches, unavailable lookup. Keep the query on errors. **Narrow view:** labeled result rows. **Pass condition:** staff can distinguish two identical names and choose by institutional ID.

### A04 — Student tuition record

**Job:** inspect one student's semester ledger and choose a precise maintenance action. **Evidence:** F02/F08, I3. **Components:** C04, C08, C10, C15.

Keep the student's name, institutional ID, and email above the semester selector; these stay visible when opening an edit form. Show the same reconciliation and ledger semantics as S07. Staff-only controls are **Record simulated payment**, **Add charge**, and payment-row **Correct record** / **Void record**. Separate record maintenance from the student's read-only balance.

Charge details expose description, amount, recorded time, and creator. Payment detail includes amount, paid date, recorded by, last updated, note, and status. A voided row is read-only. A record's update timestamp is not a complete before/after audit trail; do not render invented change history. Charge editing/removal is not currently supported and has no UI control.

Filter by semester before adding a charge or payment. In All semesters, require choosing a term before entering values. A student change returns through lookup and clears any unresolved edit form after the leave confirmation.

**States:** no records, paid/partial/unpaid/credit, voided line, stale summary after change, missing student, read error. **Narrow view:** identity, semester, reconciliation, actions, ledger in that order. **Pass condition:** the chosen student's identity and semester remain unambiguous through every maintenance action.

### A05 — Record, correct, or void a simulated payment

**Job:** change a sample payment record with clear consequences. **Evidence:** F02/F08, I3. **Components:** C03, C10, C14–C16.

A focused page opens from A04. At top: student identity, semester, mode, and existing payment reference for corrections. Do not allow changing student or semester in an existing record's correction form; cancel and navigate instead.

| Field | Create | Correct |
| --- | --- | --- |
| Amount in VND | Required integer, 1 to 1,000,000,000 | Prefill existing amount; same bounds |
| Paid date/time | Defaults to now in campus time; visible and editable | Prefill existing paid date/time |
| Note | Optional; up to 500 characters | Prefill; blank explicitly clears note |

No future paid date/time; the current service allows a small clock-skew tolerance, which is not a user-facing scheduling feature. Show grouping in the review, preserve exact integer intent in entry, and accept only unambiguous numeric input. A payment above the outstanding balance prompts a credit-balance explanation but is not automatically rejected, because I3 permits credits.

**Review before save:** identity, semester, paid time, amount, note, and balance before → projected after. Label the latter “Expected balance after saving”; the confirmed response must replace it. Final action is **Save simulated payment** or **Save correction**. Back preserves inputs.

**Void:** read-only summary and consequence “This payment will remain in the history but will no longer count toward the paid total.” Show the projected new balance and **Void payment record** in the destructive style. Safe action is Keep record. Voiding is terminal; no Restore button. Request a reason in the existing note field as a proposed UI rule, while documenting that the backend does not currently require one or retain a separate correction audit.

**Failures:** invalid amount/date, concurrently voided record, lost session, unknown save outcome. A create timeout must not offer an automatic repeat; inspect the latest ledger and follow O10. For correction, compare server values before reporting success. **Narrow view:** full-width fields and stacked before/after values. **Pass condition:** a user can predict the balance change and cannot mistake this for moving real money.

### A06 — Add a semester charge

**Job:** create a sample charge supporting the tuition ledger. **Evidence:** supporting workflow from I3; F02 requires charges to be viewable but does not explicitly require charge authoring. **Components:** C03, C10, C14–C16.

Identity and semester are fixed from A04. Fields: description (required, trimmed, maximum 500 characters) and amount in VND (required positive integer, at most 1,000,000,000). Do not compute fees from course credits without a fee rule. Show review with the exact line item and expected balance after addition; final action **Add charge**.

Explain before submission that the current prototype has no charge-edit or reversal workflow, so the entry must be checked carefully. This is a task-specific consequence, not a general legal disclaimer. A failed save keeps the form; an unknown outcome checks the ledger before another create attempt.

This page is lower priority than all brief Must flows; charges can be seeded if time is constrained. It must not consume time needed for student tuition views or staff payment updates. **Narrow view:** short single-column form. **Pass condition:** submitted description/amount/term match the resulting ledger line and the selected student's balance.

## 6. Permission design

Hiding navigation is helpful orientation; access checks must also apply to direct record visits and every action. This table specifies product permissions, not an implementation recipe.

| Data/action | Student | Lecturer | Staff |
| --- | --- | --- | --- |
| Course offering and meeting information | Browse for registration | Assigned sections in primary workflow | No MVP navigation need |
| Registration/cancellation | Own active account, allowed section/rules | No | No |
| Personal timetable/exams/results | Own only | Own assigned roster/grades, not arbitrary student results | No academic-results tool |
| Grade entry | No | Assigned sections only | No |
| Tuition ledger | Own, read-only | No | Student lookup and maintenance |
| Document submission | Own identity only | No | No impersonated student submissions |
| Document list/detail | Own only | No | All permitted requests |
| Document decision | No | No | Pending requests only |
| Account details | Own read-only | Own read-only | Own read-only; student identity only where task requires it |

Do not expose roster emails in the student catalogue. Display only identity fields necessary to distinguish records. List counts, search results, and confirmation summaries must obey the same scope as detail pages.

## 7. Shared mutation state model

Every create/update flow follows: **Ready → Validating → Saving → Confirmed**, with branches to **Needs correction**, **Definite failure**, or **Outcome unknown**. The selected action is frozen while Saving. Only a confirmed result changes the object to Registered, Saved, Approved, or Recorded.

Unknown is a real state, not an error synonym: the server may have saved the action even if its response was lost. Re-read the relevant object, compare identities and values, and only then offer a repeat. For creates lacking an operation reference, a similar record is not definitive proof; clearly state when confirmation cannot be established. Do not promise exactly-once behavior from a disabled button alone.

These rules connect J01–J06, U03, and C13. They are a direct response to the project's reliability concern, within the honest limits of a course prototype.
