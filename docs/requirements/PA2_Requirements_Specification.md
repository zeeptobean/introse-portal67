# Portal67 — Software Requirements Specification (SRS)
## Deliverable PA2: Requirements Engineering Baseline

> **Project:** Portal67 — Unified Academic & Administrative Web Portal  
> **Course:** Introduction to Software Engineering (NMCNPM) — Semester 1 (HK1), 2026–2027  
> **Institution:** Faculty of Information Technology, VNU-HCM University of Science (HCMUS)  
> **Team:** Team 02 — **BDD**  
> **Lead Architect / Author:** Vũ Trần Minh Hiếu (MSSV: 24127003)  
> **Document Version:** 2.0  
> **Status:** Formal PA2 Baseline  

---

## 1. Document Overview & System Context

### 1.1. Purpose & Scope
This document specifies the complete software requirements for **Portal67**, an integrated academic and administrative web portal developed for the students, faculty, and administrative staff of HCMUS. It covers:
- **Functional Requirements (FR):** Comprehensively numbered, logically grouped by domain/subsystem, and prioritised via the MoSCoW framework.
- **Non-Functional Requirements (NFR):** Rigorously quantified with measurable metrics, thresholds, and empirical verification methods (ISO/IEC 25010).
- **Domain Rules (DR):** Official academic and administrative policies modeling HCMUS regulations and separation of duties.
- **Acceptance Criteria (AC):** Verifiable behavior-driven scenarios formulated in Gherkin (Given-When-Then) syntax.
- **Requirements Traceability Matrix (RTM):** Bidirectional traceability connecting user requirements, domain rules, and verification criteria.

### 1.2. User Persona & Role Model
1. **Student (`STU`):** Enrolled in undergraduate programs (Standard / Đại trà or High-Quality / CLC / APCS). Performs course registration, tracks weekly timetables, inspects official grades and GPA projections, and submits administrative document requests.
2. **Lecturer (`LEC`):** Instructs assigned course sections. Manages student rosters, inputs continuous assessment grades, and submits final grade sheets for official departmental review.
3. **Academic Affairs / Staff (`ADM`):** Administrative and academic officers. Manages term offerings, supervises course registrations, audits and finalizes semester grade sheets, processes student administrative requests, and updates billing statuses.
4. **System Administrator (`SYS`):** Manages user account provisioning, role mappings, system configuration, and audit logs.

---

## 2. Functional Requirements (FR)

### 2.1. Grouping & Prioritization Conventions
- **ID Schema:** `FR-[SUBSYSTEM]-[SEQUENCE]`
- **Subsystem Codes:**
  - `AUTH`: Authentication, Identity & Access Control
  - `STU`: Student Overview & Academic Profile
  - `REG`: Course & Section Registration, Schedule Planning
  - `TT`: Timetable Visualisation & Scheduling
  - `GRD`: Gradebook, GPA Computation & Auditing
  - `REQ`: Student Affairs, Certifications & Document Workflows
  - `ADM`: Academic Affairs, Term Operations & Administration
  - `NOTIF`: Notifications & System Alerts
- **Prioritization (MoSCoW):**
  - **Must-Have (M):** Mandatory for MVP (Weeks 8–10). Core operational workflow.
  - **Should-Have (S):** High business value; targeted for completion immediately following core MVP.
  - **Could-Have (C):** Optional enhancements if development velocity permits.
  - **Won't-Have / Out of Scope (W):** Explicitly excluded from prototype baseline (e.g., live Moodle grade sync, actual banking settlement).

---

### 2.2. Master Functional Requirements Catalog

| FR ID | Subsystem | Requirement Name | Functional Description | Primary Actor | Priority | Target Milestone |
| :--- | :--- | :--- | :--- | :--- | :---: | :---: |
| **FR-AUTH-01** | `AUTH` | User Authentication | Authenticate users via verified email and password credentials, returning scoped JWT tokens (Access & Refresh tokens). | All Roles | **Must** | Phase 1 (W8) |
| **FR-AUTH-02** | `AUTH` | Demo Persona Switcher | Provide pre-configured demo logins representing Student (Standard), Student (CLC), Lecturer, and Academic Staff for presentation and assessment. | All Roles | **Must** | Phase 1 (W8) |
| **FR-AUTH-03** | `AUTH` | Role-Based Access Control | Enforce backend authorization on all protected routes and API endpoints, disallowing cross-tenant or unprivileged resource access. | All Roles | **Must** | Phase 1 (W8) |
| **FR-STU-01** | `STU` | Student Dashboard Overview | Aggregate and render real-time student summary: daily class schedule, cumulative GPA/credits completed, pending document requests, and announcements. | Student | **Must** | Phase 2 (W8) |
| **FR-STU-02** | `STU` | Biographical & Program Profile | Display student biographical data, academic program (Đại trà vs. CLC), cohort, academic advisor, and allow updating personal contact details. | Student | **Should** | Phase 5 (W10) |
| **FR-REG-01** | `REG` | Browse Course Offerings | Search and filter published course offerings by course code, title, department, program track (Standard/CLC), section type (lecture/practice), and availability. | Student | **Must** | Phase 2 (W8) |
| **FR-REG-02** | `REG` | Section Enrollment (Add) | Allow students to enroll in open lecture and practical sections during active registration windows with immediate validation against seat limits. | Student | **Must** | Phase 2 (W8) |
| **FR-REG-03** | `REG` | Section Drop | Allow students to drop registered sections within the permitted drop window, immediately relinquishing the allocated seat. | Student | **Must** | Phase 2 (W8) |
| **FR-REG-04** | `REG` | Timetable Conflict Advisory | Detect overlapping time slots between newly selected and existing sections; display prominent non-blocking visual warnings without hard-preventing registration. | Student | **Must** | Phase 2 (W8) |
| **FR-REG-05** | `REG` | Schedule Planner (Draft) | Interactive schedule generator allowing students to assemble trial timetable combinations without committing to official database enrollment. | Student | **Should** | Phase 2 (W9) |
| **FR-REG-06** | `REG` | Specialization Registration | Enable qualified students (year 2 completion) to submit ranked preferences for academic specialization tracks during designated application periods. | Student | **Could** | Phase 5 (W10) |
| **FR-TT-01** | `TT` | Personal Weekly Timetable | Render interactive grid-based timetable of enrolled sections filtered by semester, week number, and study day, indicating room, section code, and instructor. | Student, Lecturer | **Must** | Phase 2 (W8) |
| **FR-TT-02** | `TT` | Timetable Export | Export personalized schedule into standard iCalendar (.ics) format or printable PDF view. | Student, Lecturer | **Could** | Phase 5 (W10) |
| **FR-GRD-01** | `GRD` | Official Grade Inspection | Render student semester and cumulative grade reports (10-point scale, letter grades, 4-point GPA, accumulated credits) for finalized academic terms. | Student | **Must** | Phase 3 (W9) |
| **FR-GRD-02** | `GRD` | GPA Projection / Forecaster | Interactive "What-If" GPA calculation tool enabling students to input hypothetical target grades for current terms to forecast cumulative GPA impact. | Student | **Should** | Phase 3 (W9) |
| **FR-GRD-03** | `GRD` | Section Grade Entry | Allow assigned lecturers to input and update continuous component grades (attendance, midterm, coursework, final) for students in their assigned sections. | Lecturer | **Must** | Phase 3 (W9) |
| **FR-GRD-04** | `GRD` | Grade Sheet Submission | Allow lecturers to lock and submit completed section grade sheets to Academic Affairs, transitioning the state from `Draft` to `Submitted`. | Lecturer | **Must** | Phase 3 (W9) |
| **FR-GRD-05** | `GRD` | Grade Sheet Finalization | Authorize Academic Affairs staff to review, approve, and officially publish submitted grade sheets, transitioning state to `Finalized`. | Academic Staff | **Must** | Phase 3 (W9) |
| **FR-GRD-06** | `GRD` | Grade Amendment Audit Workflow | Enforce an audited change-request workflow for modifying finalized grades, recording previous grade, new grade, justification, and approver ID. | Academic Staff | **Should** | Phase 3 (W9) |
| **FR-REQ-01** | `REQ` | Document Request Submission | Enable students to apply online for formal certificates (enrollment verification, military postponement, official transcript, tuition extension). | Student | **Must** | Phase 4 (W9) |
| **FR-REQ-02** | `REQ` | Request Queue Management | Provide administrative staff with a centralized review queue to filter, inspect, approve, return for clarification, or reject student requests. | Academic Staff | **Must** | Phase 4 (W9) |
| **FR-REQ-03** | `REQ` | Request Status Tracking | Allow students to track real-time progression of submitted requests (`Submitted`, `In_Review`, `Requires_Info`, `Approved`, `Rejected`, `Ready_For_Pickup`). | Student | **Must** | Phase 4 (W9) |
| **FR-REQ-04** | `REQ` | Grade Appeal (Phúc khảo) | Allow students to lodge formal grade review requests within the statutory 7-day window following official grade finalization. | Student | **Should** | Phase 5 (W10) |
| **FR-ADM-01** | `ADM` | Registration Window Control | Permit administrators to configure, open, suspend, and close term course registration periods globally or by cohort. | Academic Staff | **Must** | Phase 2 (W8) |
| **FR-ADM-02** | `ADM` | Academic Data Management | Provision and manage courses, course offerings, lecture/practice sections, classroom capacities, and lecturer teaching assignments. | Academic Staff | **Must** | Phase 2 (W8) |
| **FR-ADM-03** | `ADM` | Tuition & Billing Simulation | Display semester tuition obligations calculated by credit hours; permit staff to update simulated payment statuses (Unpaid, Partial, Cleared). | Student, Staff | **Should** | Phase 4 (W10) |
| **FR-NOTIF-01**| `NOTIF`| In-App Notifications | Dispatch real-time system alerts to users upon critical events (grade finalized, request status changed, schedule revised, registration window opened). | All Roles | **Must** | Phase 4 (W9) |
| **FR-NOTIF-02**| `NOTIF`| System Announcement Broadcast | Authorize staff to broadcast global or faculty-specific markdown announcements visible on student and faculty dashboards. | Academic Staff | **Must** | Phase 4 (W9) |

---

## 3. Non-Functional Requirements (NFR) — Measurable Baseline

To satisfy software engineering quality standards (ISO/IEC 25010), all non-functional requirements are specified with concrete metrics, numerical target thresholds, and empirical verification procedures.

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                                 NFR MEASUREMENT MODEL                                  │
│  [Metric] ──► [Quantitative Target] ──► [Test Condition] ──► [Verification Tool]       │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

| NFR ID | Category | Quality Attribute | Measurable Metric | Quantitative Target Threshold | Verification Condition & Method |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **NFR-PERF-01** | Performance | Read Latency | Response time ($p95$) | **$p95 \le 300\text{ ms}$**, **$p99 \le 800\text{ ms}$** | Measured over 1,000 HTTP GET requests to `/api/academic/sections` and `/api/grades/mine` using k6. |
| **NFR-PERF-02** | Performance | Registration Peak Latency | Write latency ($p95$) | **$p95 \le 600\text{ ms}$**, max $\le 1500\text{ ms}$ | Measured under 50 concurrent virtual users (VUs) executing section registration transactions. |
| **NFR-SCAL-01** | Scalability | System Throughput | Requests Per Second (RPS) | **$\ge 100\text{ RPS}$** with Error Rate **$< 0.1\%$** | 3-minute sustained stress test against reverse proxy gateway on Docker containerized test node. |
| **NFR-SCAL-02** | Concurrency | Race-Condition Safety | Over-subscription Count | **Exactly 0 over-allocations** | 20 concurrent threads attempting to claim the final 1 available seat in a section. Verified via DB assertions. |
| **NFR-SEC-01** | Security | Password Hashing | Key derivation algorithm | **Argon2id** (memory: 64MB, iterations: 3, parallelism: 2) or **bcrypt** (cost $\ge 12$) | Static code analysis and cryptographic audit of `auth-service` database records. |
| **NFR-SEC-02** | Security | Token Expiry Lifetime | TTL (Time-To-Live) | Access Token: **$\le 15\text{ minutes}$**; Refresh Token: **$\le 7\text{ days}$** | Verified by unit tests rejecting expired Bearer tokens past $901\text{s}$. |
| **NFR-SEC-03** | Security | Authorization Integrity | IDOR Vulnerability Count | **0 unauthorized data accesses** | Automated integration tests executing cross-student grade retrieval requests; must return HTTP 403 Forbidden. |
| **NFR-USAB-01** | Usability | Responsive Viewport Support | Layout Breakage Count | **0 layout overflows** across viewports: $375\text{px}$ (Mobile), $768\text{px}$ (Tablet), $1280\text{px}+$ (Desktop) | Visual regression testing using Playwright across configured breakpoints. |
| **NFR-USAB-02** | Usability | Frontend Performance | First Contentful Paint (FCP) & Bundle Size | **$\text{FCP} \le 1.5\text{s}$**; Initial production JS+CSS bundle **$\le 450\text{ KB}$** (gzipped) | Google Lighthouse audit under simulated 4G mobile network profile. |
| **NFR-RELI-01** | Reliability | Data Transaction Integrity | Database Transaction Guarantee | **100% ACID compliance** on grade finalization and section enrollment | Integration tests verifying automatic rollback of multi-table updates upon simulated failure. |
| **NFR-RELI-02** | Reliability | Service Health Observability | Health Check Endpoint Latency | `/health` responds with HTTP 200 in **$\le 50\text{ ms}$** | Continuous curl polling every 10 seconds during staging runtime. |
| **NFR-MAINT-01**| Maintainability | Backend Test Coverage | Code Statement Coverage | **$\ge 70\%$ unit test coverage** on core business domain logic | Go test coverage calculation: `go test -coverprofile=coverage.out ./...`. |
| **NFR-PORT-01** | Portability | One-Command Deployment | Cold-start Deployment Time | System fully reachable via `docker compose up -d` in **$\le 60\text{ seconds}$** | Timed execution from clean git checkout on standard dual-core development host. |

---

## 4. Domain Rules (DR) — HCMUS Academic & Operational Policies

Domain rules represent inviolable business invariants and university policies governing Portal67. They exist independently of software mechanics.

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                              DOMAIN RULE TAXONOMY                                      │
│  DR-REG (Registration) │ DR-GRD (Grading & GPA) │ DR-REQ (Student Affairs) │ DR-ID     │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

### 4.1. Course Registration Domain Rules (`DR-REG`)
- **`DR-REG-01` — Active Registration Window Required:** Section registration, section change, and section dropping can only occur when the registration period assigned to the student's cohort/program has the operational status `OPEN`.
- **`DR-REG-02` — No Duplicate Course Code Enrollment:** A student cannot hold simultaneous active enrollments in multiple sections belonging to the same root Course (`course_id`) within the same academic semester, even across differing section codes.
- **`DR-REG-03` — Advisory Non-Blocking Timetable Overlaps:** In accordance with Portal67 solution architecture policy (HCMUS prototype baseline), time slot collisions between registered sections generate prominent UI warnings (`CONFLICT_DETECTED`) but do **not** programmatically block registration submission.
- **`DR-REG-04` — Hard Section Capacity Limit:** When section capacity enforcement is enabled, total enrolled students must never exceed `section.max_capacity`. Any concurrent registration attempting to exceed capacity must be rejected.
- **`DR-REG-05` — Decoupled Schedule Planning:** Saved schedule plans constructed in the Schedule Planner represent client-side or draft advisory structures; they confer no reserved seat rights or official enrollment claims until submitted and committed.

### 4.2. Grading & Academic Evaluation Domain Rules (`DR-GRD`)
- **`DR-GRD-01` — Standard 10-Point Component Scale:** All continuous and final grade components must be evaluated on a scale of $0.0$ to $10.0$ with a precision of one decimal place.
- **`DR-GRD-02` — Valid Component Weighting:** A course's grade composition (e.g., Attendance + Midterm + Lab + Final) must sum to exactly $100\%$ ($1.00$) prior to grade sheet submission.
- **`DR-GRD-03` — Semester and Cumulative GPA Formula:** 
  $$\text{Semester GPA} = \frac{\sum_{i=1}^{n} (\text{Course Grade}_i \times \text{Credits}_i)}{\sum_{i=1}^{n} \text{Credits}_i}$$
  Calculated to two decimal places. Cumulative GPA considers only the highest earned grade for retaken courses according to faculty credit regulations.
- **`DR-GRD-04` — Strict Separation of Duties (Four-Eyes Principle):** Lecturers hold authority to input and submit section grade sheets (`Draft` $\rightarrow$ `Submitted`). Only authorized Academic Affairs personnel (`ADM`) hold authority to approve, publish, and finalize official grades (`Submitted` $\rightarrow$ `Finalized`).
- **`DR-GRD-05` — Immutability of Finalized Grades:** A finalized grade sheet cannot be modified via standard update APIs. Any post-finalization grade amendment requires a formal rectification record (`Grade_Correction_Log`) specifying prior score, new score, justification reason, and approving officer ID.

### 4.3. Student Affairs & Request Workflow Domain Rules (`DR-REQ`)
- **`DR-REQ-01` — Immutable Request State Machine:** Administrative requests adhere to a unidirectional progression: `Submitted` $\rightarrow$ `Under_Review` $\rightarrow$ `Approved` / `Rejected` $\rightarrow$ `Completed`. Requests cannot be deleted from the audit ledger once submitted.
- **`DR-REQ-02` — Information Revision Cycle:** If an administrative officer marks a request as `Requires_Information`, the request transitions back to a student-editable state while preserving all historical audit notes and timestamps.
- **`DR-REQ-03` — Grade Appeal Statutory Window:** Grade appeal applications (`Phúc khảo`) are legally accepted only within $7\text{ calendar days}$ immediately following the official date of grade sheet finalization.

### 4.4. Identity & System Access Domain Rules (`DR-ID`)
- **`DR-ID-01` — Immutable Student Identifier:** The student ID (MSSV) serves as the persistent institutional key within the academic domain. An account's email may be updated, but MSSV is immutable.
- **`DR-ID-02` — Protected Academic Profile Fields:** Students are forbidden from directly modifying their academic attributes (MSSV, enrolled program, cohort, admitted major, cumulative credits, or academic status).

---

## 5. Acceptance Criteria (AC) — Behavior-Driven Specifications (BDD)

The following acceptance scenarios follow the **Given-When-Then** structure, establishing verifiable contracts between the business domain rules and system implementation.

---

### Feature: Course & Section Registration (`FR-REG-02`, `FR-REG-03`, `FR-REG-04`)

#### Scenario AC-REG-01: Successful Section Enrollment (Happy Path)
```gherkin
Scenario: Student successfully enrolls in an available section within the active window
  Given The student "24127003" is authenticated and on the Course Registration page
    And The registration window for semester "2026-1" is currently "OPEN"
    And Section "CSC10002_01" (Introduction to Programming) has 38 enrolled out of 40 max capacity
    And The student is not currently enrolled in any section of course "CSC10002"
  When The student clicks "Enroll" for section "CSC10002_01"
  Then The system records the enrollment transaction successfully
    And Returns HTTP status 200 OK with message "Successfully enrolled in CSC10002_01"
    And The active seat count for section "CSC10002_01" becomes 39/40
    And Section "CSC10002_01" immediately appears on the student's personal weekly timetable
```

#### Scenario AC-REG-02: Registration Rejected When Section Capacity is Reached
```gherkin
Scenario: System prevents enrollment when section capacity is exhausted
  Given The student is authenticated and the registration window is "OPEN"
    And Section "CSC10005_02" has reached full capacity of 50/50 enrolled students
  When The student attempts to enroll in section "CSC10005_02"
  Then The backend rejects the operation with HTTP 409 Conflict
    And Returns an error payload: '{"code": "SECTION_CAPACITY_EXCEEDED", "message": "Class section is full"}'
    And The student's enrolled courses remain unchanged
    And The section enrolled count remains strictly 50/50
```

#### Scenario AC-REG-03: Timetable Conflict Advisory Warning Displayed
```gherkin
Scenario: Student selects an overlapping section without blocking enrollment
  Given The student is already enrolled in section "MTH00001_01" scheduled on "Tuesday 07:30 - 09:10"
    And The student selects section "PHY00002_03" scheduled on "Tuesday 08:20 - 10:00"
  When The student views the registration selection or schedule planner
  Then The system flags a "Schedule Conflict Advisory" for Tuesday Periods 2-3
    And The conflict is visually highlighted in amber with an explanatory warning
    And The "Submit Registration" action remains enabled and accessible
```

---

### Feature: Grade Submission & Finalization (`FR-GRD-03`, `FR-GRD-04`, `FR-GRD-05`)

#### Scenario AC-GRD-01: Lecturer Submits Complete Section Grade Sheet
```gherkin
Scenario: Assigned lecturer locks and submits section grades for academic review
  Given Lecturer "GV001" is assigned to section "CSC10002_01"
    And The grade sheet is in "Draft" state with all component scores populated between 0.0 and 10.0
    And Component weights sum to exactly 100%
  When Lecturer "GV001" clicks "Submit to Academic Affairs"
  Then The grade sheet state transitions from "Draft" to "Submitted"
    And Lecturer "GV001" is restricted to read-only access on the submitted sheet
    And An audit log entry is appended with actor "GV001" and timestamp
    And A notification is dispatched to the Academic Affairs review queue
```

#### Scenario AC-GRD-02: Academic Affairs Finalizes and Publishes Grades
```gherkin
Scenario: Academic Affairs approves submitted grade sheet and makes grades visible to students
  Given A grade sheet for section "CSC10002_01" exists in "Submitted" state
    And An officer with role "Academic Affairs" is logged into the management portal
  When The officer clicks "Approve & Finalize Grade Sheet"
  Then The grade sheet state transitions to "Finalized"
    And The official grades immediately become visible to all enrolled students via their Grade Portal
    And Cumulative GPA and completed credits for affected students are automatically recalculated
```

#### Scenario AC-GRD-03: Direct Mutation of Finalized Grade Rejected
```gherkin
Scenario: User attempts to edit a finalized grade directly via API
  Given Grade for student "24127003" in section "CSC10002_01" is in "Finalized" state
  When A PUT request is dispatched to `/api/grades/sections/CSC10002_01/students/24127003`
  Then The backend rejects the request with HTTP 403 Forbidden
    And Returns error message: "Finalized grades cannot be directly edited. Submit a Grade Rectification Request."
```

---

### Feature: Student Document Request Processing (`FR-REQ-01`, `FR-REQ-02`)

#### Scenario AC-REQ-01: Student Submits Enrollment Certificate Request
```gherkin
Scenario: Student submits an administrative document request successfully
  Given Student "24127003" is authenticated and accesses the Document Request portal
  When The student submits a request of type "ENROLLMENT_VERIFICATION" with reason "Military deferment"
  Then A new request record is created with a unique tracking ID (e.g., "REQ-2026-0042")
    And Initial status is set to "Submitted"
    And The request appears in the student's "My Requests" tracking list
    And The request appears in the Academic Affairs pending queue
```

---

## 6. Requirements Traceability Matrix (RTM)

This matrix establishes end-to-end traceability between Functional Requirements (`FR`), Domain Rules (`DR`), Non-Functional Quality Attributes (`NFR`), and Acceptance Criteria (`AC`).

| FR ID | Functional Requirement Name | Governing Domain Rules | Bound NFRs | Bound Acceptance Criteria |
| :--- | :--- | :--- | :--- | :--- |
| **FR-AUTH-01** | User Authentication | `DR-ID-01` | `NFR-SEC-01`, `NFR-SEC-02` | `AC-AUTH-01` |
| **FR-AUTH-03** | Role-Based Access Control | `DR-GRD-04`, `DR-ID-02` | `NFR-SEC-03` | `AC-AUTH-02`, `AC-GRD-03` |
| **FR-REG-01** | Browse Course Offerings | `DR-REG-01` | `NFR-PERF-01`, `NFR-USAB-01` | `AC-REG-01` |
| **FR-REG-02** | Section Enrollment (Add) | `DR-REG-01`, `DR-REG-02`, `DR-REG-04` | `NFR-PERF-02`, `NFR-SCAL-02`, `NFR-RELI-01`| `AC-REG-01`, `AC-REG-02` |
| **FR-REG-04** | Timetable Conflict Advisory | `DR-REG-03`, `DR-REG-05` | `NFR-USAB-01` | `AC-REG-03` |
| **FR-TT-01** | Personal Weekly Timetable | `DR-REG-03` | `NFR-PERF-01`, `NFR-USAB-01` | `AC-TT-01` |
| **FR-GRD-01** | Official Grade Inspection | `DR-GRD-01`, `DR-GRD-03` | `NFR-PERF-01`, `NFR-SEC-03` | `AC-GRD-02` |
| **FR-GRD-03** | Section Grade Entry | `DR-GRD-01`, `DR-GRD-02`, `DR-GRD-04` | `NFR-SEC-03`, `NFR-RELI-01` | `AC-GRD-01` |
| **FR-GRD-04** | Grade Sheet Submission | `DR-GRD-02`, `DR-GRD-04` | `NFR-RELI-01` | `AC-GRD-01` |
| **FR-GRD-05** | Grade Sheet Finalization | `DR-GRD-04`, `DR-GRD-05` | `NFR-RELI-01`, `NFR-SEC-03` | `AC-GRD-02`, `AC-GRD-03` |
| **FR-REQ-01** | Document Request Submission | `DR-REQ-01`, `DR-REQ-02` | `NFR-PERF-01`, `NFR-RELI-01` | `AC-REQ-01` |
| **FR-REQ-02** | Request Queue Management | `DR-REQ-01` | `NFR-SEC-03` | `AC-REQ-01` |

---
*End of PA2 Requirements Specification Baseline — Team 02 (BDD)*
