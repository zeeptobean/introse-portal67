# Portal67 — Software Requirements Specification (SRS)
## Deliverable PA2: Requirements Engineering Baseline (Synced with PA3 System Design v3.0)

> **Project:** Portal67 — Unified Academic & Administrative Web Portal  
> **Course:** Introduction to Software Engineering (NMCNPM) — Semester 1 (HK1), 2026–2027  
> **Institution:** Faculty of Information Technology, VNU-HCM University of Science (HCMUS)  
> **Team:** Team 02 — **BDD**  
> **Lead Architect / Author:** Vũ Trần Minh Hiếu (MSSV: 24127003)  
> **Document Version:** 3.0  
> **Status:** Formal Requirements Baseline (Synchronized with PA3 System Design & Solution Architecture v3.0)  

---

## 1. Document Overview & System Context

### 1.1. Purpose & Scope
This document specifies the complete software requirements for **Portal67**, an integrated academic and administrative web portal engineered for the students, faculty, and administrative staff of the Faculty of Information Technology at VNU-HCM University of Science (HCMUS). It covers:
- **Functional Requirements (FR):** Comprehensively numbered, logically partitioned by domain/subsystem, and prioritised via the MoSCoW framework.
- **Non-Functional Requirements (NFR):** Rigorously quantified with measurable metrics, target thresholds, and empirical verification procedures (ISO/IEC 25010).
- **Domain Rules (DR):** Official academic and administrative policies modeling HCMUS dual-track curriculum regulations, facility constraints, and separation of duties.
- **Acceptance Criteria (AC):** Verifiable behavior-driven scenarios formulated in Gherkin (Given-When-Then) syntax.
- **Requirements Traceability Matrix (RTM):** Bidirectional traceability connecting requirements, domain rules, non-functional constraints, acceptance tests, and PA3 architectural realizing elements.
- **Requirements Evolution Log (CHG):** Formal change tracking recording deliberate design refinements from initial concept to PA3 architectural baselining.

### 1.2. User Persona & Institutional Actor Model
1. **Student (`STU`):** Enrolled in undergraduate degree programs under one of two operational tracks:
   - **Standard Track (`STANDARD` — Đại trà):** Follows the 2-semester academic calendar (HK1, HK2) managed by Phòng Đào tạo (PĐT). Primarily studies at Campus 2 (Linh Trung, Thủ Đức).
   - **Project-Based Track (`PROJECT_BASED` — CLC / Tiên tiến APCS / Việt - Pháp):** Follows the 3-trimester calendar (HK1, HK2, HK3 summer trimester) managed independently by VP Đề án (CTĐA). Primarily studies at Campus 1 (227 Nguyễn Văn Cừ, Quận 5).
   - *Key capabilities:* Browses partitioned course catalogs, executes atomic course section enrollments with integrated Theory-Practice schedules, plans trial timetables, inspects official grades and projected GPAs, and lodges administrative petitions and grade appeals.
2. **Lecturer & Teaching Assistant (`LEC`):** Instructs assigned course sections.
   - **Primary Lecturer (`is_primary = TRUE`):** Manages class rosters, enters continuous component scores, and holds exclusive authority to lock and submit completed grade sheets to Academic Affairs.
   - **Co-Lecturer / Teaching Assistant (`TA`):** Assists in lab instruction, enters laboratory/coursework scores, but is prohibited from submitting official grade sheets.
3. **Academic & Student Affairs Staff (`ADM`):** Administrative and academic officers split by jurisdiction:
   - **Phòng Đào tạo (PĐT):** Manages Standard track terms, classrooms at CS2, registration windows, and grade finalization.
   - **Văn phòng Đề án (VP Đề án - CTĐA / OEP):** Manages Project-Based terms, CS1 computer labs, trimester registration windows, and cross-registration petitions.
   - **Phòng Công tác Sinh viên (CTSV):** Verifies student statuses, processes certificate applications, and manages document fulfillment.
4. **System Administrator (`SYS`):** Provisions system user accounts, maintains identity mappings, monitors microservice health, and inspects audit ledgers.

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

| FR ID | Subsystem | Requirement Name | Functional Description | Primary Actor | Priority | Target Milestone | Realizing Service |
| :--- | :--- | :--- | :--- | :--- | :---: | :---: | :--- |
| **FR-AUTH-01** | `AUTH` | User Authentication | Authenticate users via institutional identifier (Student ID MSSV / Staff ID MSCB) and password credentials or via institutional OAuth/SSO, returning scoped JWT tokens (Access token 15m, Refresh token 7d). No public self-registration (accounts are pre-provisioned). | All Roles | **Must** | Phase 1 (W8) | `auth-service` |
| **FR-AUTH-02** | `AUTH` | Demo Persona Switcher | Provide pre-configured demo logins representing Student (Standard), Student (CLC), Lecturer, and Academic Staff for presentation and assessment. | All Roles | **Should** | Phase 1 (W8) | `auth-service` |
| **FR-AUTH-03** | `AUTH` | Role-Based Access Control | Enforce backend authorization on all protected routes and API endpoints, disallowing cross-tenant or unprivileged resource access. | All Roles | **Must** | Phase 1 (W8) | `gateway` + all |
| **FR-STU-01** | `STU` | Student Dashboard Overview | Aggregate and render real-time student summary: daily class schedule, cumulative GPA/credits completed, pending document requests, and announcements. | Student | **Must** | Phase 2 (W8) | `academic-service` |
| **FR-STU-02** | `STU` | Biographical & Program Profile | Display student biographical data, academic program track (`STANDARD` vs. `PROJECT_BASED`), cohort, advisor, and allow updating personal contact details. | Student | **Should** | Phase 5 (W10) | `academic-service` |
| **FR-REG-01** | `REG` | Browse Course Offerings | Search and filter published course offerings with smart program-cohort catalog partitioning (`STANDARD` vs `PROJECT_BASED`), department, campus location (CS1 vs CS2), and integrated Theory-Practice session schedule view. | Student | **Must** | Phase 2 (W8) | `registration-service` |
| **FR-REG-02** | `REG` | Atomic Section Enrollment (Add) | Allow students to enroll in open course sections during active registration windows. For courses with labs (`has_lab = TRUE`), registering for a section automatically enrolls and schedules both Theory and Practice sessions, atomically reserving the section seat via high-concurrency in-memory engine. | Student | **Must** | Phase 2 (W8) | `registration-service` |
| **FR-REG-03** | `REG` | Atomic Section Drop | Allow students to drop registered sections within permitted drop window, immediately relinquishing the allocated section seat (releasing both Theory and Practice scheduled slots) in memory and persisting to database. | Student | **Must** | Phase 2 (W8) | `registration-service` |
| **FR-REG-04** | `REG` | Timetable Conflict Advisory | Detect overlapping time slots between newly selected and existing sections; display prominent non-blocking visual warnings (`CONFLICT_DETECTED`) without hard-preventing registration submission (ADR-004). | Student | **Must** | Phase 2 (W8) | `registration-service` |
| **FR-REG-05** | `REG` | Schedule Planner (Draft) | Interactive schedule generator allowing students to assemble trial timetable combinations without committing to official database enrollment or reserving seats. | Student | **Should** | Phase 2 (W9) | `registration-service` |
| **FR-REG-06** | `REG` | Specialization Registration | Enable qualified students (year 2 completion) to submit ranked preferences for academic specialization tracks during designated application periods. | Student | **Could** | Phase 5 (W10) | `registration-service` |
| **FR-REG-07** | `REG` | Cross-Registration Petition | Allow Project-Based students to request enrollment in eligible Standard track sections (`is_cross_enrollment_allowed = TRUE`) with advisory tuition settlement warnings and staff approval. | Student, Staff | **Should** | Phase 5 (W10) | `registration-service` |
| **FR-TT-01** | `TT` | Personal Weekly Timetable | Render interactive grid-based timetable of enrolled sections filtered by semester, week number, and study day, showing room code, building, campus (CS1/CS2), section code, and instructor. | Student, Lecturer | **Must** | Phase 2 (W8) | `academic-service` |
| **FR-TT-02** | `TT` | Timetable Export | Export personalized schedule into standard iCalendar (.ics) format or printable PDF view. | Student, Lecturer | **Could** | Phase 5 (W10) | `academic-service` |
| **FR-GRD-01** | `GRD` | Official Grade Inspection | Render student semester and cumulative grade reports (10-point scale, letter grades, 4-point GPA, accumulated credits) broken down by dynamic component rubrics for finalized academic terms. | Student | **Must** | Phase 3 (W9) | `academic-service` |
| **FR-GRD-02** | `GRD` | GPA Projection / Forecaster | Interactive "What-If" GPA calculation tool enabling students to input hypothetical target grades for current terms to forecast cumulative GPA impact. | Student | **Should** | Phase 3 (W9) | `academic-service` |
| **FR-GRD-03** | `GRD` | Dynamic Rubric Grade Entry | Allow assigned primary lecturers and TAs to input and update continuous component scores across dynamic, configurable assessment rubrics (`academic.grade_components`). | Lecturer, TA | **Must** | Phase 3 (W9) | `academic-service` |
| **FR-GRD-04** | `GRD` | Grade Sheet Submission | Allow only the Primary Lecturer (`is_primary = TRUE`) to lock and submit completed section grade sheets to Academic Affairs, transitioning state from `Draft` to `Submitted`. | Lecturer | **Must** | Phase 3 (W9) | `academic-service` |
| **FR-GRD-05** | `GRD` | Grade Sheet Finalization | Authorize Academic Affairs staff to review, approve, and officially publish submitted grade sheets (`Finalized`), triggering automatic GPA and credit recalculations. | Academic Staff | **Must** | Phase 3 (W9) | `academic-service` |
| **FR-GRD-06** | `GRD` | Grade Amendment Audit Workflow | Enforce an audited change-request workflow for modifying finalized grades, recording previous score, new score, justification reason, and approving officer ID in `Grade_Correction_Log`. | Academic Staff | **Should** | Phase 3 (W9) | `academic-service` |
| **FR-REQ-01** | `REQ` | Document Request Submission | Enable students to apply online for formal certificates (enrollment verification, military postponement, official transcript, tuition extension) with document upload. | Student | **Must** | Phase 4 (W9) | `affairs-service` |
| **FR-REQ-02** | `REQ` | Request Queue Management | Provide administrative staff with a centralized review queue to filter, inspect, approve, return for clarification, or reject student requests. | Academic Staff | **Must** | Phase 4 (W9) | `affairs-service` |
| **FR-REQ-03** | `REQ` | Request Status Tracking | Allow students to track real-time progression of submitted requests (`Submitted`, `Under_Review`, `Requires_Information`, `Approved`, `Rejected`, `Ready_For_Pickup`). | Student | **Must** | Phase 4 (W9) | `affairs-service` |
| **FR-REQ-04** | `REQ` | Grade Appeal (Phúc khảo) | Allow students to lodge formal grade review requests within the appeal window configured and opened by Academic Affairs staff following official grade finalization, routing to departmental re-assessment. | Student | **Should** | Phase 5 (W10) | `affairs-service` |
| **FR-ADM-01** | `ADM` | Registration Window Control | Permit administrators (PĐT or VP Đề án) to configure, open, suspend, and close term course registration periods segregated by academic track and cohort. | Academic Staff | **Must** | Phase 2 (W8) | `registration-service` |
| **FR-ADM-02** | `ADM` | Academic Data Management | Provision and manage courses, section Theory-Practice session schedules, classroom facilities with campus allocations (CS1/CS2), and co-teaching / TA assignments. | Academic Staff | **Must** | Phase 2 (W8) | `academic-service` |
| **FR-ADM-03** | `ADM` | Tuition & Billing Simulation | Display semester tuition obligations calculated by credit hours; permit staff to update simulated payment statuses (Unpaid, Partial, Cleared). | Student, Staff | **Should** | Phase 4 (W10) | `affairs-service` |
| **FR-NOTIF-01**| `NOTIF`| Asynchronous In-App Alerts | Dispatch real-time system alerts to users upon critical events (grade finalized, request status changed, schedule revised, registration window opened) via decoupled background workers. | All Roles | **Must** | Phase 4 (W9) | `notif-service` |
| **FR-NOTIF-02**| `NOTIF`| System Announcement Broadcast | Authorize staff to broadcast global, program-specific, or faculty markdown announcements visible on student and faculty dashboards. | Academic Staff | **Must** | Phase 4 (W9) | `notif-service` |

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
| **NFR-PERF-02** | Performance | Peak Registration Latency | Write latency ($p95$) | Redis Lua reservation: **$\le 10\text{ ms}$**; HTTP write $p95 \le \mathbf{300\text{ ms}}$, max $\le 800\text{ ms}$ | Measured under 100 concurrent virtual users (VUs) executing section registration transactions against `registration-service`. |
| **NFR-SCAL-01** | Scalability | System Throughput | Requests Per Second (RPS) | **$\ge 100\text{ RPS}$** overall; **$\ge 500\text{ RPS}$** on Redis seat reservation engine with Error Rate **$< 0.1\%$** | 3-minute sustained stress test against reverse proxy gateway on Docker containerized test node. |
| **NFR-SCAL-02** | Concurrency | Race-Condition Safety | Over-subscription Count | **Exactly 0 over-allocations** across Theory and Practice seats | 50 concurrent threads attempting to claim the final 1 available seat in a section. Verified via automated DB assertions. |
| **NFR-SEC-01** | Security | Password Hashing | Key derivation algorithm | **Argon2id** (memory: 64MB, iterations: 3, parallelism: 2) or **bcrypt** (cost $\ge 12$) | Static code analysis and cryptographic audit of `auth-service` database records. |
| **NFR-SEC-02** | Security | Token Expiry Lifetime | TTL (Time-To-Live) | Access Token: **$\le 15\text{ minutes}$**; Refresh Token: **$\le 7\text{ days}$** | Verified by unit tests rejecting expired Bearer tokens past $901\text{s}$. |
| **NFR-SEC-03** | Security | Authorization Integrity | IDOR Vulnerability Count | **0 unauthorized data accesses** | Automated integration tests executing cross-student grade retrieval requests; must return HTTP 403 Forbidden. |
| **NFR-USAB-01** | Usability | Responsive Viewport Support | Layout Breakage Count | **0 layout overflows** across viewports: $375\text{px}$ (Mobile), $768\text{px}$ (Tablet), $1280\text{px}+$ (Desktop) | Visual regression testing using Playwright across configured breakpoints. |
| **NFR-USAB-02** | Usability | Frontend Performance | First Contentful Paint (FCP) & Bundle Size | **$\text{FCP} \le 1.5\text{s}$**; Initial production JS+CSS bundle **$\le 450\text{ KB}$** (gzipped) | Google Lighthouse audit under simulated 4G mobile network profile. |
| **NFR-USAB-03** | Usability | Accessibility Compliance | WCAG 2.1 AA Checklist | **Contrast $\ge 4.5:1$**, 100% keyboard navigable, ARIA labels on form controls, Lighthouse Accessibility $\ge 95$ | Audited with Axe-core, Lighthouse, and manual keyboard navigation. |
| **NFR-RELI-01** | Reliability | Data Transaction Integrity | Database Transaction Guarantee | **100% ACID compliance** on grade finalization and enrollment queue persistence | Integration tests verifying automatic rollback of multi-table updates upon simulated failure. |
| **NFR-RELI-02** | Reliability | Service Health Observability | Health Check Endpoint Latency | `/health` responds with HTTP 200 in **$\le 50\text{ ms}$** | Continuous curl polling every 10 seconds across all 5 microservices during staging runtime. |
| **NFR-MAINT-01**| Maintainability | Backend Test Coverage | Code Statement Coverage | **$\ge 70\%$ unit test coverage** on core business domain logic | Go test coverage calculation: `go test -coverprofile=coverage.out ./...`. |
| **NFR-PORT-01** | Portability | One-Command Deployment | Cold-start Deployment Time | System fully reachable via `docker compose up -d` in **$\le 60\text{ seconds}$** | Timed execution from clean git checkout on standard dual-core development host. |
| **NFR-ARCH-01** | Architecture | Microservice Fault Isolation | Failure Cascade Rate | **0% cascade failure** on grading or student affairs when registration experiences peak saturation | Validated by load-testing `registration-service` to 100% CPU while issuing read requests to `/api/grades/mine`. |

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
- **`DR-REG-01` — Active Registration Window Required:** Section registration, change, and dropping can only occur when the registration period assigned to the student's program track (`STANDARD` vs `PROJECT_BASED`) and cohort has the operational status `OPEN`.
- **`DR-REG-02` — No Duplicate Course Code Enrollment:** A student cannot hold simultaneous active enrollments in multiple sections belonging to the same root Course (`course_id`) within the same academic semester, even across differing section codes.
- **`DR-REG-03` — Advisory Non-Blocking Timetable Overlaps (ADR-004):** In accordance with HCMUS prototype baseline policy, time slot collisions between registered sections generate prominent UI warnings (`CONFLICT_DETECTED`) but do **not** programmatically block registration submission.
- **`DR-REG-04` — Hard Section Capacity Limit:** When section capacity enforcement is enabled, total enrolled students must never exceed `section.max_capacity`. Any concurrent registration attempting to exceed capacity must be rejected immediately.
- **`DR-REG-05` — Decoupled Schedule Planning:** Saved schedule plans constructed in the Schedule Planner represent draft advisory structures; they confer no reserved seat rights or official enrollment claims until submitted and committed.
- **`DR-REG-06` — Program Track Cohort Isolation (ADR-008, BR-REG-11):** Standard track students (`STANDARD`) cannot register Project-Based sections (`PROJECT_BASED`). Project-Based students may only enroll in Standard sections when `is_cross_enrollment_allowed = TRUE` under formal academic petition.
- **`DR-REG-07` — Integrated Theory-Practice Enrollment (ADR-008, BR-REG-12):** For courses with laboratory components (`has_lab = TRUE`), each course section encapsulates both Theory and Practice lecture/lab schedules. Registering for the section automatically enrolls and schedules both Theory and Practice sessions without separate or unbundled registration. Selecting a section commits the student to both required schedules.
- **`DR-REG-08` — In-Memory Atomic Reservation (ADR-007, BR-REG-13):** During peak windows, seat quotas are reserved atomically via single-threaded Redis Lua scripts (`seat_reserve.lua`) prior to asynchronous relational persistence, guaranteeing 0 over-allocations without database row lock contention.

### 4.2. Grading & Academic Evaluation Domain Rules (`DR-GRD`)
- **`DR-GRD-01` — Standard 10-Point Component Scale:** All continuous and final grade components must be evaluated on a scale of $0.0$ to $10.0$ with a precision of one decimal place.
- **`DR-GRD-02` — Dynamic Assessment Rubric Weighting (CHG-05):** A course's grade composition is defined dynamically via $k$ grade components (`academic.grade_components`). The sum of all active component weights must equal exactly $100\%$ ($1.00$) prior to grade sheet submission.
- **`DR-GRD-03` — Semester and Cumulative GPA Formula:** 
  $$\text{Semester GPA} = \frac{\sum_{i=1}^{n} (\text{Course Grade}_i \times \text{Credits}_i)}{\sum_{i=1}^{n} \text{Credits}_i}$$
  Calculated to two decimal places. Cumulative GPA considers only the highest earned grade for retaken courses according to faculty credit regulations.
- **`DR-GRD-04` — Strict Separation of Duties (Four-Eyes Principle):** Lecturers hold authority to input and submit section grade sheets (`Draft` $\rightarrow$ `Submitted`). Only authorized Academic Affairs personnel (`ADM`) hold authority to approve, publish, and finalize official grades (`Submitted` $\rightarrow$ `Finalized`).
- **`DR-GRD-05` — Immutability of Finalized Grades & Audit Trail:** A finalized grade sheet cannot be modified via standard update APIs. Any post-finalization grade amendment requires a formal rectification record (`Grade_Correction_Log`) specifying prior score, new score, justification reason, and approving officer ID.
- **`DR-GRD-06` — Multi-Lecturer Assignment Permissions (BR-GRD-10):** In co-taught sections (`academic.section_lecturers`), co-professors and TAs may enter and update continuous component scores; however, only the designated Primary Lecturer (`is_primary = TRUE`) holds permission to lock and submit the completed grade sheet.

### 4.3. Student Affairs & Request Workflow Domain Rules (`DR-REQ`)
- **`DR-REQ-01` — Immutable Request State Machine:** Administrative requests adhere to a unidirectional progression: `Submitted` $\rightarrow$ `Under_Review` $\rightarrow$ `Approved` / `Rejected` $\rightarrow$ `Completed`. Requests cannot be deleted from the audit ledger once submitted.
- **`DR-REQ-02` — Information Revision Cycle:** If an administrative officer marks a request as `Requires_Information`, the request transitions back to a student-editable state while preserving all historical audit notes and timestamps.
- **`DR-REQ-03` — Configurable Grade Appeal Window (`Phúc khảo`):** Grade appeal applications are accepted strictly within the appeal window configured and opened by Academic Affairs staff (`appeal_start_date` to `appeal_end_date`) following official grade finalization. Academic Affairs retains administrative authority to configure the appeal window duration.

### 4.4. Identity & System Access Domain Rules (`DR-ID`)
- **`DR-ID-01` — Immutable Student Identifier:** The student ID (MSSV) serves as the persistent institutional key within the academic domain. An account's email may be updated, but MSSV is immutable.
- **`DR-ID-02` — Protected Academic Profile Fields:** Students are forbidden from directly modifying their academic attributes (MSSV, enrolled program track, cohort, admitted major, cumulative credits, or academic status).
- **`DR-ID-03` — Closed Institutional Portal & Pre-Provisioned Accounts:** Portal67 does not offer public user self-registration. All student, lecturer, and administrative staff accounts are pre-provisioned via university administrative directories. Users authenticate using their institutional identifier (MSSV / MSCB) and password or through university-authorized OAuth/SSO.

---

## 5. Acceptance Criteria (AC) — Behavior-Driven Specifications (BDD)

The following acceptance scenarios follow the **Given-When-Then** structure, establishing verifiable contracts between the business domain rules and system implementation.

---

### Feature: User Authentication & Access Control (`FR-AUTH-01`, `FR-AUTH-03`)

#### Scenario AC-AUTH-01: Successful Authentication with Institutional Identifier (Happy Path)
```gherkin
Scenario: Authenticated user receives valid scoped access and refresh tokens
  Given A pre-provisioned user exists with institutional identifier "24127003" (MSSV) and password "CorrectPassword123!"
    And The account status is "ACTIVE" with role "STUDENT"
  When The user submits a POST request to "/api/auth/login" with identifier "24127003" and password
  Then The system returns HTTP status 200 OK
    And The response body contains an Access Token with expiry <= 15 minutes
    And The response body contains a Refresh Token with expiry <= 7 days
    And The Access Token claims contain "sub", "role: STUDENT", and student MSSV "24127003"
```

#### Scenario AC-AUTH-02: Backend Authorization Guard Rejects Unauthorized Access (Edge Case)
```gherkin
Scenario: Student attempts to access administrative endpoint
  Given A student is authenticated with role "STUDENT" and holds a valid Bearer token
  When The student sends a GET request to "/api/affairs/requests/queue"
  Then The API gateway rejects the request with HTTP status 403 Forbidden
    And Returns an error payload: '{"code": "FORBIDDEN", "message": "Insufficient permissions"}'
    And An audit security log entry is recorded with actor ID, IP address, and timestamp
```

#### Scenario AC-AUTH-03: Successful Authentication via Institutional OAuth/SSO
```gherkin
Scenario: Pre-provisioned user authenticates via Institutional OAuth/SSO
  Given An institutional user account "24127003" is pre-provisioned in Portal67 with role "STUDENT"
    And The user authenticates successfully via university OAuth identity provider with email "24127003@student.hcmus.edu.vn"
  When The OAuth callback returns authorization code to "/api/auth/oauth/callback"
  Then The system verifies the token, resolves identity to student MSSV "24127003"
    And Issues an Access Token (expiry <= 15 minutes) and Refresh Token (expiry <= 7 days)
```

---

### Feature: Course & Section Registration (`FR-REG-01`, `FR-REG-02`, `FR-REG-03`, `FR-REG-04`)

#### Scenario AC-REG-01: Successful Section Enrollment with Integrated Theory and Practice Sessions (Happy Path)
```gherkin
Scenario: Student successfully enrolls in a course section containing Theory and Practice schedules
  Given The student "24127003" is authenticated and on the Course Registration page
    And The registration window for semester "2026-1" is currently "OPEN"
    And Course "CSC10002" (Introduction to Programming) has has_lab = TRUE
    And Section "CSC10002_01" encapsulates both Theory (Monday 07:30 - 10:00) and Practice (Wednesday 13:00 - 16:30)
    And Section "CSC10002_01" has 38 enrolled out of 40 max capacity
    And The student is not currently enrolled in any section of course "CSC10002"
  When The student selects section "CSC10002_01"
    And Clicks "Confirm Registration"
  Then The registration service executes atomic seat reservation in Redis Lua
    And The transaction is committed to the relational database
    And Returns HTTP status 200 OK with message "Successfully enrolled in CSC10002_01"
    And The active seat count becomes 39/40 for section "CSC10002_01"
    And Both Theory lecture and Practice lab sessions appear automatically on the student's personal weekly timetable
```

#### Scenario AC-REG-02: Registration Rejected When Section Capacity is Reached
```gherkin
Scenario: System prevents enrollment when section capacity is exhausted
  Given The student is authenticated and the registration window is "OPEN"
    And Section "CSC10002_01" has reached full capacity of 40/40 enrolled students
  When The student attempts to enroll in section "CSC10002_01"
  Then The Redis Lua script rejects the operation with error "CAPACITY_EXCEEDED"
    And The backend returns HTTP 409 Conflict with payload:
      '{"code": "SECTION_CAPACITY_EXCEEDED", "message": "Selected section CSC10002_01 is full"}'
    And Section seat quota is not decremented
    And The student's enrolled courses remain unchanged
```
    And Neither Theory nor Practice seats are decremented
    And The student's enrolled courses remain unchanged
```

#### Scenario AC-REG-03: Timetable Conflict Advisory Warning Displayed (ADR-004)
```gherkin
Scenario: Student selects an overlapping section without blocking enrollment
  Given The student is already enrolled in section "MTH00001_01" scheduled on "Tuesday 07:30 - 09:10"
    And The student selects section "PHY00002_03" scheduled on "Tuesday 08:20 - 10:00"
  When The student views the registration selection or schedule planner
  Then The system flags a "Schedule Conflict Advisory" for Tuesday Periods 2-3
    And The conflict is visually highlighted in amber with an explanatory warning
    And The "Submit Registration" action remains enabled and accessible
```

#### Scenario AC-REG-04: Cross-Track Cohort Isolation Enforced (ADR-008)
```gherkin
Scenario: Standard track student attempts to enroll in a Project-Based section
  Given The student "24127003" has academic program track "STANDARD"
    And Section "CSC10002_CLC01" has target_program_type = "PROJECT_BASED"
    And The section has is_cross_enrollment_allowed = FALSE
  When The student attempts to submit registration for "CSC10002_CLC01"
  Then The registration service rejects the request with HTTP 403 Forbidden
    And Returns error message: "Enrollment restricted: Section is reserved for Project-Based students"
```

---

### Feature: Grade Submission & Finalization (`FR-GRD-03`, `FR-GRD-04`, `FR-GRD-05`)

#### Scenario AC-GRD-01: Primary Lecturer Submits Complete Dynamic Rubric Grade Sheet
```gherkin
Scenario: Assigned primary lecturer locks and submits section grades for academic review
  Given Lecturer "GV001" is assigned as Primary Lecturer (is_primary = TRUE) to section "CSC10002_01"
    And The section grade sheet is in "Draft" state
    And Grade components (Attendance [10%], Midterm [30%], Lab [20%], Final [40%]) sum to exactly 100%
    And All enrolled students have valid numeric scores populated between 0.0 and 10.0
  When Lecturer "GV001" clicks "Submit to Academic Affairs"
  Then The grade sheet state transitions from "Draft" to "Submitted"
    And Lecturer "GV001" and assigned TAs are restricted to read-only access on the sheet
    And An immutable audit log entry is appended with actor "GV001" and timestamp
    And A notification is dispatched to the Academic Affairs review queue
```

#### Scenario AC-GRD-02: Academic Affairs Finalizes Grades and Recalculates GPA
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

#### Scenario AC-GRD-04: Teaching Assistant Cannot Submit Official Grade Sheet
```gherkin
Scenario: Teaching Assistant attempts to finalize or submit grade sheet
  Given User "TA001" is assigned to section "CSC10002_01" with is_primary = FALSE
  When User "TA001" dispatches a POST request to `/api/grades/sections/CSC10002_01/submit`
  Then The backend rejects the request with HTTP 403 Forbidden
    And Returns error message: "Only the designated Primary Lecturer can submit section grade sheets"
```

---

### Feature: Student Document Request & Grade Appeal (`FR-REQ-01`, `FR-REQ-04`)

#### Scenario AC-REQ-01: Student Submits Enrollment Certificate Request
```gherkin
Scenario: Student submits an administrative document request successfully
  Given Student "24127003" is authenticated and accesses the Document Request portal
  When The student submits a request of type "ENROLLMENT_VERIFICATION" with reason "Military deferment"
  Then A new request record is created with a unique tracking ID (e.g., "REQ-2026-0042")
    And Initial status is set to "Submitted"
    And The request appears in the student's "My Requests" tracking list
    And The request appears in the Student Affairs (CTSV) pending queue
```

#### Scenario AC-REQ-02: Student Submits Grade Appeal Within Configured Appeal Window
```gherkin
Scenario: Student lodges a grade appeal within the window opened by Academic Affairs
  Given Section "CSC10002_01" grade sheet is finalized
    And Academic Affairs has configured and opened an active grade appeal window for this term
    And Student "24127003" has an official grade recorded for this section
  When The student submits a grade appeal for "CSC10002_01" with justification and payment slip
  Then The system accepts the petition and assigns status "Under_Departmental_Review"
    And An alert is dispatched to the FIT Academic Department for re-grading
```

---

## 6. Requirements Traceability Matrix (RTM)

This matrix establishes end-to-end traceability between Functional Requirements (`FR`), Domain Rules (`DR`), Non-Functional Constraints (`NFR`), Acceptance Criteria (`AC`), and realizing architectural components in **PA3**.

| FR ID | Functional Requirement Name | Governing Domain Rules | Bound NFRs | Bound ACs | Realizing PA3 Microservice | API Route / Contract | Target Database Table |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **FR-AUTH-01** | User Authentication | `DR-ID-01`, `DR-ID-03` | `NFR-SEC-01`, `NFR-SEC-02` | `AC-AUTH-01`, `AC-AUTH-03` | `auth-service` | `POST /api/auth/login`, `POST /api/auth/oauth/callback` | `auth.accounts` |
| **FR-AUTH-03** | Role-Based Access Control | `DR-GRD-04`, `DR-GRD-06`, `DR-ID-02` | `NFR-SEC-03` | `AC-AUTH-02`, `AC-GRD-03`, `AC-GRD-04` | `gateway` + all | HTTP Auth Bearer Header | `auth.roles`, `auth.account_roles` |
| **FR-STU-01** | Student Dashboard Overview | `DR-REG-01`, `DR-GRD-03` | `NFR-PERF-01`, `NFR-USAB-01` | `AC-STU-01` | `academic-service` | `GET /api/academic/dashboard/overview` | `academic.enrollments`, `academic.sections` |
| **FR-REG-01** | Browse Course Offerings | `DR-REG-01`, `DR-REG-06` | `NFR-PERF-01`, `NFR-USAB-01` | `AC-REG-01`, `AC-REG-04` | `registration-service` | `GET /api/registration/sections` | Pre-warmed Redis / `academic.sections` |
| **FR-REG-02** | Atomic Section Enrollment (Add) | `DR-REG-01`, `DR-REG-02`, `DR-REG-04`, `DR-REG-07`, `DR-REG-08` | `NFR-PERF-02`, `NFR-SCAL-01`, `NFR-SCAL-02`, `NFR-RELI-01` | `AC-REG-01`, `AC-REG-02`, `AC-REG-04` | `registration-service` | `POST /api/registration/sections/register` | `registration.enrollments` |
| **FR-REG-03** | Atomic Section Drop | `DR-REG-01`, `DR-REG-07`, `DR-REG-08` | `NFR-PERF-02`, `NFR-RELI-01` | `AC-REG-01` | `registration-service` | `DELETE /api/registration/sections/:id` | `registration.enrollments` |
| **FR-REG-04** | Timetable Conflict Advisory | `DR-REG-03`, `DR-REG-05` | `NFR-USAB-01` | `AC-REG-03` | `registration-service` | Client Svelte Engine + Conflict API | `academic.section_schedules`, `academic.classrooms` |
| **FR-REG-05** | Schedule Planner (Draft) | `DR-REG-05` | `NFR-USAB-01` | `AC-REG-03` | `registration-service` | `POST /api/registration/planner/save` | `registration.planner_drafts` |
| **FR-REG-07** | Cross-Registration Petition | `DR-REG-06` | `NFR-SEC-03` | `AC-REG-04` | `registration-service` | `POST /api/registration/cross-petition` | `registration.cross_enrollment_petitions` |
| **FR-TT-01** | Personal Weekly Timetable | `DR-REG-03`, `DR-REG-07` | `NFR-PERF-01`, `NFR-USAB-01` | `AC-TT-01` | `academic-service` | `GET /api/academic/timetable/mine` | `academic.section_schedules`, `academic.classrooms` |
| **FR-GRD-01** | Official Grade Inspection | `DR-GRD-01`, `DR-GRD-02`, `DR-GRD-03` | `NFR-PERF-01`, `NFR-SEC-03` | `AC-GRD-02` | `academic-service` | `GET /api/grades/mine` | `academic.student_grades`, `academic.grade_components` |
| **FR-GRD-03** | Dynamic Rubric Grade Entry | `DR-GRD-01`, `DR-GRD-02`, `DR-GRD-06` | `NFR-SEC-03`, `NFR-RELI-01` | `AC-GRD-01`, `AC-GRD-04` | `academic-service` | `PUT /api/grades/sections/:id` | `academic.student_grades`, `academic.grade_components` |
| **FR-GRD-04** | Grade Sheet Submission | `DR-GRD-02`, `DR-GRD-04`, `DR-GRD-06` | `NFR-RELI-01` | `AC-GRD-01`, `AC-GRD-04` | `academic-service` | `POST /api/grades/sections/:id/submit` | `academic.grade_sheets` |
| **FR-GRD-05** | Grade Sheet Finalization | `DR-GRD-04`, `DR-GRD-05` | `NFR-RELI-01`, `NFR-SEC-03` | `AC-GRD-02`, `AC-GRD-03` | `academic-service` | `POST /api/grades/sections/:id/finalize` | `academic.grade_sheets` |
| **FR-GRD-06** | Grade Amendment Audit Workflow | `DR-GRD-05` | `NFR-RELI-01`, `NFR-SEC-03` | `AC-GRD-03` | `academic-service` | `POST /api/grades/amendments` | `academic.grade_correction_logs` |
| **FR-REQ-01** | Document Request Submission | `DR-REQ-01`, `DR-REQ-02` | `NFR-PERF-01`, `NFR-RELI-01` | `AC-REQ-01` | `affairs-service` | `POST /api/affairs/requests` | `affairs.student_requests` |
| **FR-REQ-02** | Request Queue Management | `DR-REQ-01` | `NFR-SEC-03` | `AC-REQ-01` | `affairs-service` | `GET /api/affairs/requests/queue` | `affairs.student_requests` |
| **FR-REQ-04** | Grade Appeal (Phúc khảo) | `DR-REQ-03`, `DR-GRD-05` | `NFR-PERF-01`, `NFR-SEC-03` | `AC-REQ-02` | `affairs-service` | `POST /api/affairs/grade-appeals` | `affairs.grade_appeals` |
| **FR-ADM-01** | Registration Window Control | `DR-REG-01`, `DR-REG-06` | `NFR-SEC-03` | `AC-ADM-01` | `registration-service` | `POST /api/registration/windows` | `registration.registration_windows` |
| **FR-ADM-02** | Academic Data Management | `DR-REG-06`, `DR-REG-07` | `NFR-SEC-03`, `NFR-RELI-01` | `AC-ADM-02` | `academic-service` | `POST /api/academic/sections` | `academic.sections`, `academic.section_lecturers` |
| **FR-NOTIF-01**| Asynchronous In-App Alerts | `DR-REQ-01`, `DR-GRD-05` | `NFR-PERF-01`, `NFR-ARCH-01` | `AC-NOTIF-01` | `notif-service` | `GET /api/notifications/mine` | `affairs.notifications` |

---

## 7. Requirements Evolution & Refinement Log (PA2 Baseline $\rightarrow$ PA3 Synced v3.0)

This log provides full backward and forward auditability of requirement refinements between the initial PA2 requirements phase and the PA3 detailed architectural design phase:

1. **Refinement `CHG-01` — Advisory Timetable Conflict Policy Clarified (`FR-REG-04`, `DR-REG-03`, `ADR-004`):**
   - *Initial Concept:* Considered programmatically blocking registration submissions that contain period collisions.
   - *Refined Baseline:* Standardized on prominent amber visual advisory warnings without hard blocking, matching authentic university practice where students may attend overlapping courses with advisor consent.
2. **Refinement `CHG-02` — Concurrency-Safe In-Memory Seat Allocation (`FR-REG-02`, `NFR-SCAL-02`, `DR-REG-08`, `ADR-007`):**
   - *Initial Concept:* High-level statement that classes must not exceed capacity via database row-level locking (`SELECT FOR UPDATE`).
   - *Refined Baseline:* Replaced relational locking with atomic Redis Lua scripts (`seat_reserve.lua`) achieving sub-2ms seat reservations and asynchronous write-behind persistence queues (Redis Streams), guaranteeing 0 over-allocations during 8:00 AM registration flash crowds.
3. **Refinement `CHG-03` — Four-Eyes Principle Separation in Grade Lifecycle (`FR-GRD-04`, `FR-GRD-05`, `DR-GRD-04`):**
   - *Initial Concept:* Lecturers submit grades and students view them.
   - *Refined Baseline:* Formalized strict separation of duties: Lecturers can only transition sheets to `Submitted`. A separate administrative role (`Academic Affairs`) must review, approve, and officially publish (`Finalized`) before cumulative GPAs recalculate.
4. **Refinement `CHG-04` — Decoupled Asynchronous Notifications (`FR-NOTIF-01`, `NFR-ARCH-01`):**
   - *Initial Concept:* Core academic service synchronously dispatches emails/alerts upon document status changes.
   - *Refined Baseline:* Extracted notification handling into a dedicated `notification-service` consuming events via background message workers, protecting primary transactional APIs from external network latencies.
5. **Refinement `CHG-05` — Dynamic Assessment Rubrics & Strict BCNF Normalization (`FR-GRD-03`, `DR-GRD-02`):**
   - *Initial Concept:* Assumed 4 hard-coded grade columns (`attendance`, `midterm`, `coursework`, `final`).
   - *Refined Baseline:* Decomposed into dynamic `academic.grade_components` (custom criteria, weights, descriptions) and `academic.student_grades` (per-student score entries tied to enrollments). Eliminates non-superkey functional dependency violations ($\{components\} \rightarrow overall\_score$), reaching 100% Boyce-Codd Normal Form (BCNF) compliance.
6. **Refinement `CHG-06` — Decoupling Course Registration into Dedicated Microservice (`NFR-ARCH-01`, `ADR-007`):**
   - *Initial Concept:* Registration engine operated inside monolithic academic backend.
   - *Refined Baseline:* Extracted course registration into an independent microservice (`registration-service` on Port 8085) with isolated schema (`registration`), ensuring registration flash crowds cannot degrade core institutional operations (grading, transcript viewing, student affairs).
7. **Refinement `CHG-07` — HCMUS Multi-Track Academic Model Adaptation & Integrated Theory-Practice Section Scheduling (`FR-REG-01`, `FR-REG-02`, `FR-REG-07`, `DR-REG-06`, `DR-REG-07`, `ADR-008`):**
   - *Initial Concept:* Homogeneous semester structure, single-lecturer assignment, and monolithic course offerings.
   - *Refined Baseline:* Fully realized HCMUS academic operations:
     a) Segregated term structures: Standard track (2 semesters/yr, Phòng Đào tạo) vs. Project-Based track (3 trimesters/yr including mandatory HK3 summer term, VP Đề án).
     b) Integrated Theory-Practice section schedules: courses with lab components bind both Theory and Practice meeting slots within a single section offering, automatically registering both upon section selection with unified atomic in-memory reservation.
     c) Physical facility normalization (`academic.classrooms`) with campus identification (CS1 Nguyễn Văn Cừ vs CS2 Linh Trung).
     d) Co-teaching faculties and lab TAs via `academic.section_lecturers` ($N:M$).
     e) Smart cohort catalog partitioning (`catalog:{term}:STANDARD` vs `catalog:{term}:PROJECT_BASED`) and cross-enrollment restrictions.

---
*End of PA2 Software Requirements Specification Baseline v3.0 — Team 02 (BDD)*
