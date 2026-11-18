# Portal67 component and content specification

Read with the [experience design](portal67-experience-design.md) and [page specifications](portal67-pages-and-flows.md). These are design contracts for behavior and appearance, not code components. Each component exists because a named task needs it. Default, focus, loading, unavailable, error, and confirmed states must be designed together.

## 1. Component catalog

### C01 — Prototype identity

**Used:** U01 and the application shell. **Need:** users must understand that this is a course prototype using sample records.

Use a typographic Portal67 wordmark; the “67” can use the ink accent without becoming a decorative badge. Follow it with small, legible “Course prototype” text. On sign-in, explain the scope in one sentence. In tuition and document outcomes, repeat only the relevant limitation at the point it affects interpretation.

No official seal, claimed university endorsement, uptime percentage, or badge implying certified financial/academic records. The label must remain readable on narrow screens. It is static information, not a dismissible warning repeated on every action.

### C02 — Role shell and page header

**Used:** every signed-in page. **Need:** stable orientation across different academic tasks.

Anatomy: wordmark; labeled navigation; active item; account area; page title; local context/actions. The current navigation item uses text weight, an ink-colored leading mark, and a pale background, so color is not the sole signal. Subpages retain the parent's active state and add a short breadcrumb back to the list.

At narrow widths, Menu opens a labeled drawer. Focus moves into it, stays inside while open, Escape closes, and focus returns to Menu. Opening navigation never submits or discards a form. A skip link reaches the main content. Page title, current role, and account identity are not hidden behind an unexplained icon.

Keep one strong primary action per task region. A gradebook's sticky Save is the same action as its contextual save control, not a competing duplicate. An active sidebar or toolbar must never obscure a focused input.

### C03 — Labeled field

**Used:** sign-in, request forms, payment/charge forms, grade inputs, search. **Need:** people must know what to enter and recover without re-entry.

Anatomy: visible label; required/optional status; concise hint when needed; input; field error; character count where a real limit exists. Label precedes the input, not inside a placeholder. Associate label and error with the input. Use comfortable 40–44 px control height; 44 px touch targets are the project preference.

Validate on submit and, after the first failed submission, when the affected value changes or loses focus. Do not show an error while a user has only begun typing a valid value. Trim outer whitespace where the contract does; do not rewrite the meaning of names or purpose text. Preserve all unrelated fields after an error.

Password: support paste/managers and Show/Hide with a stateful accessible label. Numeric score: distinguish blank and zero. Amount: visible VND label, integer entry, explicit grouping on review. A character limit counts user-perceived text as closely as the service permits; the current service counts Unicode code points, so test Vietnamese input and avoid a mismatched byte counter.

### C04 — Semester selector

**Used:** S01–S07, L01/L02 where appropriate, A04–A06. **Need:** the same course and balance may exist in multiple semesters.

Use a labeled select, not a large hero card. Show a full readable semester name and year, with the short code secondary if needed. The selected value remains visible after changing it. Use current/next/past descriptions only when date data supports them.

Changing term updates the affected page only, closes an incompatible section preview, and resets list pagination. Query text may persist if useful, but the user sees its continued presence. Unsaved grade/payment/request data must not be discarded by a term change. A confirmation summary repeats the semester independently of the selector.

If term loading fails, show Retry instead of an empty dropdown interpreted as “no semesters.” The selector does not certify registration is open. On narrow screens it sits beneath the title, never in a horizontally scrolling toolbar.

### C05 — Course and section identity

**Used:** S02/S03/S04, L01/L02. **Need:** people recognize a subject by code/name but register and grade a specific section.

Course group: code, full name, credits, kind. Section row: section code, lecturer, meetings, seats or enrollment state, action. Always pair a section code with its course identity when detached from a group. A plain `01` is not an adequate dialog title.

All meeting times are visible before registration. Room may be absent; use “Room not provided.” Name and code wrap rather than ellipsize essential distinguishing text. Specialized topic is a plain type label using the same registration behavior as other courses, not an unrelated colored module.

Seat count is `12 of 40 seats available`, or compact `12 / 40 available` with an equivalent accessible name. Zero says Full. Do not use scarcity animation or invent a “Low seats” threshold. Capacity unavailable says Availability unavailable and requires a fresh authoritative attempt; it is not zero.

### C06 — Registration state and receipt

**Used:** S02/S03. **Need:** distinguish interest, a request in flight, and a confirmed enrollment.

| State | Visual / action contract |
| --- | --- |
| Offered, no known blocker | View section; explicit Register action after all key facts are visible |
| Preview only | Dashed candidate label; does not affect registered count or credit total |
| Registering | Action label changes; repeat submit blocked; record remains in place |
| Registered | Text + check mark; enrollment reference/details available; Cancel separately named |
| Full | Non-actionable Full text; other sections discoverable |
| Conflict | Warning label and exact conflict explanation; no tooltip dependency |
| Cancelled | Retained in history; excluded from active count and timetable |
| Outcome unknown | Neutral pending-verification band; no false check mark or automatic repeat |

Success receipt contains course, section, semester, and the confirmed action. A toast can supplement it, but the updated object must remain as durable proof after the toast disappears. Do not display a reservation timer, basket total, or unsupported enrollment rank.

### C07 — Timetable and comparison

**Used:** S01/S03/S04. **Need:** compare where time is occupied, without translating a long list mentally.

Anatomy: day headers; labeled time axis; meeting blocks; explicit legend when preview or conflict exists; accessible agenda equivalent. Consistent block time positioning makes comparison meaningful. Avoid free placement or decorative variation.

Registered blocks have a solid rule; the selected candidate uses a dashed rule and Preview only. Conflicts have an icon and text, plus a linked explanation outside the grid. Meeting lanes must keep both overlapping sessions readable. Selecting a block opens details rather than dragging it.

The agenda is available on desktop and default on mobile. If a grid needs internal scrolling, label the region and keep the text alternative easy to find. Keyboard users must reach every meeting without a complex custom calendar interaction. Screen readers should hear course, section, weekday, time range, and room as one coherent record.

### C08 — Record table / record list

**Used:** registration, results, exams, request lists, tuition, rosters. **Need:** compare repetitive facts efficiently.

Use genuine headings and row relationships. Align names and descriptions left; monetary amounts, credits, and grades right with tabular numerals. Use subtle row separators, not a full box around every cell. A hover background can aid tracking but must not imply the whole row is clickable when only the named link is.

Sorting controls show their current direction and scope. Status/action columns do not silently sort themselves after a mutation. Focus and scroll remain stable. For small demo data, show the whole useful set; for larger lists, use explicit pagination with total results and 25-row default. Do not add infinite scrolling to a queue or gradebook.

At narrow widths, transform read-only tables into labeled records while retaining every consequential field. A gradebook retains identity plus score rather than forcing a detached input offscreen. Horizontal scrolling is a last resort for genuinely two-dimensional data, and never the only way to reach a primary action.

### C09 — GPA with evidence

**Used:** S06. **Need:** the number must communicate its limits and be inspectable.

Anatomy: Estimated GPA label; value or Not available; semester/cumulative scope; graded credits; completeness; rule disclosure. Pair the selected semester GPA with e.g. “8 of 11 credits graded.” Pair cumulative GPA with “Across available semesters · 19 graded credits.” Never label a selected-term number Cumulative.

Calculation disclosure includes formula, bands, exclusions, rounding, and a small worked example. It must correspond to the rules that generated the number. No GPA chart is necessary for the MVP; if added later, it must earn its space through a history-comparison task and have a table alternative.

Zero with graded credits is a valid calculated number; zero graded credits means no GPA. Avoid rewarding/punishing color categories or inferred scholarship eligibility without policy. A student's score is sensitive information, not a gamified progress badge.

### C10 — Tuition reconciliation and ledger

**Used:** S07, A04–A06. **Need:** explain why a balance has its value.

Show charges, effective recorded payments, and balance in a single connected group. Explain credit balances in plain language. The status comes from totals; do not let staff manually choose Paid independently of the ledger.

Rows retain voided entries, their status, and exact integer VND amounts. Use a signed mathematical value only in the calculation explanation; the headline can say “Credit balance · 200,000 VND.” Grouped VND formatting should follow the selected language consistently.

In staff review, show Before and Expected after; after confirmation show the server-confirmed ledger. When a refreshed ledger differs from the preview, explain that records changed and require review before proceeding. Audit history beyond current available fields is a documented gap, not a decorative timeline.

### C11 — Document type, status, and event history

**Used:** S08–S10, A01/A02. **Need:** students and staff must share one understanding of a request.

Types have names and one-line descriptions. Enrollment confirmation means a request for confirmation of enrollment; Transcript means a request for academic results documentation; Other document uses the purpose field to identify the need. These are choices in the prototype, not promises of document formats.

Lifecycle: Pending → Approved **or** Pending → Rejected. Do not insert Processing, Printed, Collected, or Completed steps. Pending uses a neutral blue/gray label; Approved uses success text; Rejected uses explicit text and subdued error styling. Status is not communicated by dots alone.

History event anatomy: action/status, actor, absolute campus timestamp, comment if present. Show oldest to newest in the compact two-event history; current state remains above it. Student-visible comments must not be labeled Internal note. Preserve line breaks in purpose/comments without accepting arbitrary rich formatting.

### C12 — Contextual exception

**Used:** conflicts, stale records, service problems, pending verification. **Need:** show what happened, why it matters, and what the user can do.

Anatomy: short title; precise affected object; explanation; one recovery action if available. Examples: “This section overlaps your timetable” with both meetings; “This request was already reviewed” with actual decision; “Availability could not be refreshed” with retry.

Use caution for an identified conflict, error for a failed/invalid action, and neutral information for unknown save outcome. Do not stack a toast, modal, red page banner, and row error describing the same fact. Global banners are reserved for global conditions; otherwise put the message near the affected record.

### C13 — Save feedback and uncertain outcomes

**Used:** every mutation. **Need:** slow responses are a central project concern, and a lost response does not prove a failed save.

Immediately acknowledge the click with a stable action state. If saving takes more than roughly 1 second, keep a textual Saving message; after about 10 seconds, explain “Still waiting for confirmation.” These are proposed feedback thresholds to test, not response-time guarantees. No percentage progress for a request whose progress is unknown.

When transport times out, say “We couldn't confirm whether this was saved.” First inspect the authoritative record. For registration, check the exact student's active enrollment in the section. For grade edits, compare exact intended scores. For request decisions, load the decision/history. For payment corrections, inspect the target record. For new requests/payments/charges, similar values alone cannot prove identity; O10 requires an operation reference for reliable automated deduplication.

Disable repeated actions during reconciliation and explain why. Do not hold unrelated navigation hostage. If the user leaves, the destination must fetch current state rather than assume success or failure. Never repeatedly replay a write in the background.

### C14 — Form error summary

**Used:** multi-field forms and gradebook batch validation. **Need:** users must locate mistakes, especially when fields are below the fold.

Show a summary near the form title and a matching error beside each affected field. Move focus to the summary after a failed submit; each entry focuses the correct field or grade row. Use the same language in both places. A service outage is not “Please correct your fields.” This adapts the [GOV.UK error-summary pattern](https://design-system.service.gov.uk/components/error-summary/).

Example: “Enter a score from 0 to 10 for student DEMO-028.” For long rosters, include identity and a link, not “Invalid value on row 24,” because row order can change. Keep valid edited rows intact. On resubmission, the summary updates to unresolved errors only.

### C15 — Read-only identity and review summary

**Used:** document forms/decisions, staff tuition changes, account. **Need:** a correct value applied to the wrong person is still an incorrect operation.

Show name, institutional ID, and the task-relevant semester/reference. Email is included where staff need disambiguation, not everywhere. Missing fields have explicit labels. A review is human-readable, with ordinary amounts and dates, not backend field names.

For consequential changes, show the exact intended values and edit/back controls immediately before confirmation. Keep entered data when returning to edit. This adapts the [GOV.UK check-answers pattern](https://design-system.service.gov.uk/patterns/check-answers/) to payment changes and terminal decisions; two-field document submission uses an inline review instead of an extra page.

### C16 — Consequence confirmation

**Used:** cancellation, payment voiding, final staff decisions. **Need:** prevent an action whose cost cannot be reliably undone.

Title states the action and object, such as “Cancel CSC13002 · section 01?” Body states the practical consequence. Buttons say Keep registration / Cancel registration, not No / Yes. Initial focus goes to the safe action or heading; Escape cancels; closing returns focus to the trigger. Confirmation is not a license to skip server validation.

Use a dialog for a short consequence with no data entry. Use a dedicated review region/page for a comment or payment correction. Never nest dialogs, make a long roster editable inside one, or use a modal as a substitute for a durable request receipt.

### C17 — Search, filters, sorting, and pagination

**Used:** catalogue, rosters, queues, student lookup, history. **Need:** repeated records need predictable narrowing and return paths.

Search has a persistent descriptive label and explicit clear control. Put only common filters directly beside it; specialized topics is important enough to remain visible in registration. State active filters in plain text and keep them on empty results. Changing a filter resets pagination; browser Back restores the previous query/filter/page/scroll.

Data loads can be delayed after settled typing, but stale responses must not replace a newer query. Announce result counts politely without reading every row aloud. Search is scoped to the module. For server-paginated data, filtering only the currently displayed page is incorrect and must not be represented as a complete result.

### C18 — Receipt and next step

**Used:** S10 and confirmed staff changes. **Need:** users should know what completed and how to find it later.

State the completed action, record reference, relevant identity/term, and a genuine next step. Request detail is a durable receipt that survives refresh. Payment save returns to a ledger with the affected record identifiable. This follows the completion principle of the [GOV.UK confirmation pattern](https://design-system.service.gov.uk/patterns/confirmation-pages/).

Do not say an email was sent, a document was produced, or a bank accepted money unless such an operation actually exists and was confirmed. A transient toast alone is insufficient. A prototype reference must not be styled as a certified certificate number.

### C19 — Grade editing state

**Used:** L02. **Need:** repeated input needs precise save scope and interruption handling.

Track visible states: Not graded, Saved, Edited, Invalid, Saving, and Needs verification. Show changed count and any hidden edits. Editing one score never silently saves another. Discard restores confirmed values, not zeros.

Save only intended changes; preserve stable identity and sort order. A roster refresh while dirty offers review of differences rather than replacing edits. The proposed conflict behavior needs an eventual revision contract to be reliable under concurrent saves. No Auto-saved label appears unless a real autosave design is separately adopted.

## 2. State distinctions that must survive visual design

| Easily confused states | Required distinction |
| --- | --- |
| No data / failed fetch | Empty-state explanation versus Retry error |
| Ungraded / score zero | Not graded versus 0.00 and included credits |
| Preview / registration | Dashed Preview only versus confirmed Registered |
| Paid / no charges | Paid only with ledger evidence; No tuition records otherwise |
| Credit / outstanding debt | Positive credit wording, not an unexplained negative amount |
| Approved request / issued document | Explicit approval of workflow with no issuance claim |
| Pending request / unsaved form | Request ID and history versus local unsaved input |
| Save failed / outcome unknown | Definite rejection versus verification before another write |
| Ended term / registration window closed | Current term-ended rule versus an unimplemented date-window policy |
| Last fetched / last changed | Client retrieval timestamp versus record update timestamp |

## 3. Language and microcopy

Use short, direct sentences and task names people recognize. Do not expose words such as mutation, transaction lock, API, entity, enum, or validation payload. Technical error details belong in diagnostics, not the student's form.

English is used to keep this design specification readable alongside PA1. Vietnamese-first UI is recommended for the HCMUS audience, pending O03. If adopted, translate the whole flow consistently; do not mix English navigation with Vietnamese errors. Preserve course names as supplied by the catalogue. A language switch is deferred unless both translations can be maintained.

| Concept | Proposed Vietnamese UI term | English specification term |
| --- | --- | --- |
| Course | Học phần | Course |
| Section | Lớp học phần | Section |
| Registration | Đăng ký học phần | Registration |
| Timetable | Thời khóa biểu | Class timetable |
| Exams | Lịch thi | Exams |
| Results | Kết quả học tập | Results |
| Estimated GPA | GPA ước tính | Estimated GPA |
| Tuition | Học phí | Tuition |
| Document request | Yêu cầu giấy tờ | Document request |
| Pending | Chờ xét duyệt | Pending review |
| Approved / rejected | Đã duyệt / Từ chối | Approved / Rejected |
| Simulated payment record | Bản ghi thanh toán mô phỏng | Simulated payment record |

Core messages should identify the record and offer the next sensible action:

| Situation | Example copy |
| --- | --- |
| Registration confirmed | “Registered for CSC14005 · section 01 in Semester 2, 2026.” |
| Last seat lost | “This section is now full. View other sections or refresh availability.” |
| Class clash | “Monday 09:30–11:30 overlaps CSC13002 · section 01.” |
| Registered elsewhere | “You are already registered for this course in section 01.” |
| Cancellation consequence | “Your seat will be released. Registering again depends on availability.” |
| No grades | “No grades are available yet. GPA will appear when a course is graded.” |
| Grade input invalid | “Enter a score from 0 to 10, with no more than two decimal places.” |
| Unknown save | “We couldn't confirm whether this was saved. Check the latest record before trying again.” |
| Request saved | “Request #104 has been submitted. You can check its status here.” |
| Request approved | “Approved in this prototype. No official document has been issued.” |
| Same request decided elsewhere | “This request has already been reviewed. The latest decision is shown below.” |
| Void payment | “This record will remain visible but will no longer count toward the paid total.” |
| Partial load | “Couldn't load payment records. The balance cannot be confirmed.” |

Do not use “Something went wrong” alone, “Invalid input” without a remedy, “Congratulations” for routine saves, or “Success!” without an object. Do not blame the user for network/service failures.

Dates in prose use an unambiguous format such as 25 Oct 2026 or 25 thg 10, 2026 in the chosen locale. Input formats must be labeled. Times use 24-hour notation. All academic times are campus-local UTC+7 even when a device is elsewhere. Record histories may add “UTC+7” explicitly. Currency never changes into a decimal-valued foreign format because of an English UI.

## 4. Accessibility specification

Target WCAG 2.2 AA for the implementation. This document is a design target, not a conformance audit. Relevant criteria include 1.4.3 Contrast (Minimum), 1.4.10 Reflow, 1.4.11 Non-text Contrast, 2.1.1 Keyboard, 2.4.7 Focus Visible, 2.4.11 Focus Not Obscured, 2.5.8 Target Size, 3.3.1 Error Identification, 3.3.2 Labels or Instructions, and 4.1.3 Status Messages. [WCAG 2.2 reference](https://www.w3.org/TR/WCAG22/).

Design-specific checks:

- Measure normal text contrast at least 4.5:1, large text 3:1, and meaningful control boundaries/state indicators 3:1. Pale separator rules may be decorative; input boundaries must remain identifiable.
- Use a visible 2 px ink focus outline with 2 px separation. On a dark primary button, combine a light separation edge and dark outer ring. Avoid clipping outlines in table cells.
- Aim for 44 × 44 px interactive hit areas in this product; preserve at least the applicable WCAG 24 × 24 CSS px minimum or permitted spacing exception. Dense text links need adequate spacing.
- Keep reading order consistent with visual priority. Announce result counts and saved status without moving focus unnecessarily. Urgent validation gets focused error summary; ordinary success stays polite.
- Dialogs manage focus and Escape, return focus on close, and have a visible close/safe action. Drawers follow the same rule. Background content is not keyboard-active while a modal is open.
- Every color-coded state also has text; every timetable has an agenda. Do not put essential rules in hover-only tooltips.
- Retain full tasks at 200% zoom and 320 CSS px reflow. Sticky save bars must not cover fields when the mobile keyboard opens.
- Respect reduced motion. Loading uses static or subtle indicators, never flashing. Focus must not jump during availability refresh.
- Confirm Vietnamese diacritics, long course names, duplicate student names, and 1,000-character purposes wrap without clipping.

These are acceptance obligations for later implementation; a screenshot cannot prove keyboard or screen-reader behavior.

## 5. Loading, freshness, and privacy

Use row-shaped skeletons only when a read is in progress; reserve stable title, filters, and column widths. Loading a section list should not erase the term selector or typed query. For slow operations show the text state; do not add a page-covering spinner for a local save.

Previously loaded data can remain visible after a failed refresh with an explicit “Could not refresh” marker. When partial failure prevents a trustworthy derived total, suppress the total and explain which data is missing. Never calculate a reassuring tuition balance from only successfully loaded payments or charges.

No private drafts in permanent browser storage for this MVP. In-memory drafts survive recoverable errors in the same tab, subject to U03's same-account reauthentication rule. Sign-out clears user-specific data. Do not put names, purpose text, grades, or tuition amounts in URL query strings; list context can use non-sensitive term/status identifiers. Institutional ID search should not be echoed into shareable URLs by default.

Minimize shared-screen exposure: no GPA/tuition amount on the global sidebar, no full roster outside the lecturer's section, and no document-purpose excerpt in broad navigation. These choices preserve usefulness while limiting incidental disclosure.

## 6. Visual consistency checks

Use the token values from the main document. Reserve primary ink fill for the next task action, outlines for alternatives, and text links for navigation. Destructive fill belongs only to the final destructive action, not every row's initial Cancel link. Status labels are smaller than actions and are not clickable unless explicitly links.

Use one shadow only where layering needs explanation, such as a floating menu or side panel; static content relies on spacing and rules. Keep short metadata in normal case. Do not center paragraphs or long table text. Allow full comments to wrap. Icons supplement familiar verbs and should not create an icon-only toolbar.

Freshness comes from disciplined composition: a preview joined to registration, a readable explanation attached to GPA, and a decision attached to its request. A screen that could belong to any generic analytics dashboard has lost the product's purpose.
