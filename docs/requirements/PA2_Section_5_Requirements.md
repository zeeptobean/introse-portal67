# Section 5 — Software Requirements Specification

> **Course:** Introduction to Software Engineering (NMCNPM) — Semester 1 (HK1), 2026–2027  
> **Institution:** Faculty of Information Technology, VNU-HCM University of Science (HCMUS)  
> **Team:** Team 02 — **BDD**  
> **Document:** Software Requirements Specification (SRS) — Section 5 Deliverable  
> **Author:** Vũ Trần Minh Hiếu (MSSV: 24127003)  
> **Traceability Reference:** [PA2_Requirements_Specification.md](file:///D:/Hieu/university/3rd_year/HK1/NMCNPM/Project/introse-portal67/docs/requirements/PA2_Requirements_Specification.md) (v3.0)  

---

## 5.1 Functional Requirements

*[Use "The system shall ..." statements. Number the requirements and group them by feature. Give each a MoSCoW priority and a source (stakeholder, interview, document).]*

### Feature Definitions:
- **`F01` — User Authentication, Identity & Access Control** (`AUTH`)
- **`F02` — Student Overview & Academic Profile** (`STU`)
- **`F03` — Course & Section Registration, Schedule Planning** (`REG`)
- **`F04` — Timetable Visualisation & Scheduling** (`TT`)
- **`F05` — Gradebook, Dynamic Evaluation & Auditing** (`GRD`)
- **`F06` — Student Affairs, Certifications & Document Workflows** (`REQ`)
- **`F07` — Academic Operations, Term Management & Administration** (`ADM`)
- **`F08` — Asynchronous Notifications & System Announcements** (`NOTIF`)

---

### Master Functional Requirements Table

| ID        | Requirement                                                                                                                                                                                                                                  | Feature | Priority   | Source                                                               |
| :--------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :-------: | :----------: | :-------------------------------------------------------------------- |
| **FR-01** | The system shall authenticate users via institutional identifier (Student ID MSSV / Staff ID MSCB) and password credentials or via institutional OAuth/SSO, returning cryptographically signed JWT tokens (Access Token lifetime $\le 15$ min, Refresh Token lifetime $\le 7$ days). Public self-registration is strictly omitted; all accounts are pre-provisioned by administration. | `F01`   | **Must**   | HCMUS Identity Management Standards &amp; PA1 Project Baseline       |
| **FR-02** | The system shall provide pre-configured demo switchers representing Student (Standard), Student (Project-Based CLC), Lecturer, and Academic Affairs Staff for course evaluation and testing.                                                 | `F01`   | **Should** | NMCNPM Course Project Assessment Guidelines                          |
| **FR-03** | The system shall enforce backend role-based access control (RBAC) on all protected API routes, rejecting unprivileged or cross-tenant data requests with HTTP 403 Forbidden.                                                                 | `F01`   | **Must**   | ISO/IEC 27001 Access Control Policy &amp; Solution Architecture v3.0 |
| **FR-04** | The system shall aggregate and render a real-time student dashboard displaying daily class schedules, cumulative GPA, total earned credits, pending document requests, and active announcements.                                             | `F02`   | **Must**   | Stakeholder Interview: FIT Undergraduate Students                    |
| **FR-05** | The system shall display the student's biographical data, academic program track (`STANDARD` vs. `PROJECT_BASED`), cohort, and advisor, while allowing students to update personal contact details.                                          | `F02`   | **Should** | HCMUS Student Handbook 2026 (*Sổ tay sinh viên*)                     |
| **FR-06** | The system shall allow students to browse, search, and filter course offerings with smart program-cohort catalog partitioning (`STANDARD` vs `PROJECT_BASED`), department, campus location (CS1 vs CS2), and integrated Theory-Practice session schedules. | `F03`   | **Must**   | FIT Course Offering &amp; Registration Regulations &amp; ADR-008     |
| **FR-07** | The system shall allow students to enroll in open course sections during active registration windows; for courses with lab components (`has_lab = TRUE`), enrolling in a section automatically registers and schedules both the Theory and Practice sessions of that section, atomically reserving the section seat via high-concurrency in-memory engine. | `F03`   | **Must**   | FIT Credit-Based Training Guidelines &amp; ADR-007                   |
| **FR-08** | The system shall allow students to drop registered sections within permitted drop windows, atomically relinquishing the allocated section seat (releasing both Theory and Practice scheduled slots) in memory and updating relational records. | `F03`   | **Must**   | HCMUS Academic Affairs Drop Regulations                              |
| **FR-09** | The system shall detect overlapping time slots between newly selected and enrolled sections, displaying prominent non-blocking amber advisory warnings (`CONFLICT_DETECTED`) without preventing submission.                                  | `F03`   | **Must**   | HCMUS Academic Advisory Policy &amp; ADR-004                         |
| **FR-10** | The system shall provide an interactive schedule planner allowing students to generate and inspect trial timetable combinations without committing to official database enrollment or claiming seats.                                        | `F03`   | **Should** | Undergraduate Student Council Survey &amp; Solution Architecture     |
| **FR-11** | The system shall enable qualified students completing Year 2 to submit ranked preferences for academic specialization tracks during designated application periods.                                                                          | `F03`   | **Could**  | FIT Academic Handbook — Specialization Allocation Procedures         |
| **FR-12** | The system shall allow Project-Based students to submit cross-registration petitions for eligible Standard track sections (`is_cross_enrollment_allowed = TRUE`) with advisory tuition settlement warnings.                                  | `F03`   | **Should** | VP Đề án (CTĐA) Cross-Enrollment Policy &amp; ADR-008                |
| **FR-13** | The system shall render an interactive weekly timetable of enrolled sections filtered by semester, week number, and study day, indicating room code, building, campus (CS1/CS2), section code, and instructor.                               | `F04`   | **Must**   | HCMUS Timetable Standards &amp; ADR-008                              |
| **FR-14** | The system shall export the user's personalized weekly schedule into standard iCalendar (.ics) format or printable PDF view.                                                                                                                 | `F04`   | **Could**  | Usability Improvement Proposal                                       |
| **FR-15** | The system shall render student semester and cumulative grade reports (10-point scale, letter grades, 4-point GPA, accumulated credits) broken down by dynamic component rubrics for finalized academic terms.                               | `F05`   | **Must**   | MoET Circular 08/2021/TT-BGDĐT on Undergraduate Grading              |
| **FR-16** | The system shall provide an interactive "What-If" GPA projection tool enabling students to simulate hypothetical target scores for current terms and evaluate impact on cumulative GPA.                                                      | `F05`   | **Should** | FIT Academic Advisor Committee Recommendations                       |
| **FR-17** | The system shall allow assigned primary lecturers and teaching assistants to input and update continuous component scores across dynamic, configurable assessment rubrics (`academic.grade_components`).                                     | `F05`   | **Must**   | FIT Lecturer Grading Regulations &amp; CHG-05                        |
| **FR-18** | The system shall permit only the designated Primary Lecturer (`is_primary = TRUE`) to lock and submit completed section grade sheets to Academic Affairs, transitioning status from `Draft` to `Submitted`.                                  | `F05`   | **Must**   | FIT Department Regulations on Separation of Grading Duties           |
| **FR-19** | The system shall authorize Academic Affairs staff to review, approve, and officially publish submitted grade sheets (`Finalized`), automatically triggering GPA and credit recalculations for all enrolled students.                         | `F05`   | **Must**   | HCMUS Academic Affairs Grade Administration Policy                   |
| **FR-20** | The system shall enforce an immutable change-request workflow for modifying finalized grades, recording prior score, new score, justification reason, and approving officer ID in `Grade_Correction_Log`.                                    | `F05`   | **Should** | ISO/IEC 27001 Audit Trail &amp; MoET Examination Inspection Protocol |
| **FR-21** | The system shall enable students to submit formal administrative certificate applications online (enrollment verification, military postponement, official transcript, tuition extension) with digital document upload.                      | `F06`   | **Must**   | HCMUS Student Affairs Office (Phòng CTSV) Regulations                |
| **FR-22** | The system shall provide administrative staff with a centralized review queue to filter, inspect, approve, return for clarification, or reject student requests.                                                                             | `F06`   | **Must**   | CTSV Officer Workflow Interview &amp; Specifications                 |
| **FR-23** | The system shall enable students to track real-time progression of submitted requests through a defined unidirectional status lifecycle.                                                                                                     | `F06`   | **Must**   | Student Experience Improvement Baseline                              |
| **FR-24** | The system shall allow students to lodge formal grade review requests (*Phúc khảo*) within the appeal window configured and opened by Academic Affairs staff following official grade finalization.                                         | `F06`   | **Should** | HCMUS Academic Regulations on Grade Appeals (*Quy chế Phúc khảo*)    |
| **FR-25** | The system shall permit administrators (PĐT or VP Đề án) to configure, open, suspend, and close term course registration periods segregated by academic program track and cohort.                                                            | `F07`   | **Must**   | PĐT &amp; VP Đề án Operational Workflow Manual                       |
| **FR-26** | The system shall allow administrative staff to provision and manage courses, section Theory-Practice session schedules, classroom facilities with campus allocations (CS1/CS2), and co-teaching/TA assignments.                             | `F07`   | **Must**   | FIT Curriculum Administration Standards &amp; ADR-008                |
| **FR-27** | The system shall display semester tuition obligations calculated by enrolled credit hours and allow staff to update simulated payment statuses.                                                                                              | `F07`   | **Should** | HCMUS Tuition Billing Regulations                                    |
| **FR-28** | The system shall dispatch real-time asynchronous in-app notifications upon critical lifecycle events (grade finalized, request status changed, registration window opened) via decoupled background workers.                                 | `F08`   | **Must**   | Solution Architecture v3.0 &amp; ADR-007                             |
| **FR-29** | The system shall authorize administrative staff to broadcast global, program-specific, or faculty-wide markdown announcements visible on user dashboards.                                                                                    | `F08`   | **Must**   | FIT Communication &amp; Public Information Guidelines                |

---

## 5.2 Non-Functional Requirements

*[Every non-functional requirement must be verifiable: state a metric and the conditions. Example: "Under 200 concurrent users, the search page shall return results in $\le 1$ s at P95."]*

| ID | Category | Requirement (measurable) | How It Will Be Verified |
| :--- | :--- | :--- | :--- |
| **NFR-01** | Performance | Under 1,000 sustained HTTP GET requests to `/api/academic/sections` and `/api/grades/mine`, the system shall return responses with $p95 \le 300\text{ ms}$ and $p99 \le 800\text{ ms}$. | Measured using k6 load-testing scripts executing against the containerized reverse proxy gateway. |
| **NFR-02** | Performance | Under 100 concurrent virtual users executing section enrollment transactions, the `registration-service` shall execute Redis Lua seat reservations in $\le 10\text{ ms}$ and complete end-to-end HTTP requests with $p95 \le 300\text{ ms}$ and max $\le 800\text{ ms}$. | Measured using k6 concurrent write scenarios targeting `/api/registration/sections/register`. |
| **NFR-03** | Scalability | Under a 3-minute sustained stress test, the system shall sustain $\ge 100\text{ RPS}$ overall and $\ge 500\text{ RPS}$ on the in-memory seat reservation engine with an HTTP error rate $< 0.1\%$. | Verified via automated k6 stress-test harness running in standard Docker Compose test node. |
| **NFR-04** | Concurrency / Integrity | When 50 concurrent threads simultaneously attempt to claim the final 1 remaining seat in a section, the system shall maintain **exactly 0 over-allocations** for the section (encompassing both integrated Theory and Practice schedules). | Automated multi-threaded integration test issuing concurrent POST requests; verified via database count assertions. |
| **NFR-05** | Security | The system shall hash 100% of stored user passwords using **Argon2id** (memory: 64MB, iterations: 3, parallelism: 2) or **bcrypt** (work factor $\ge 12$). | Static source code security review and direct cryptographic audit of `auth.accounts` database records. |
| **NFR-06** | Security | User access tokens shall have a maximum Time-To-Live (TTL) of **$\le 15\text{ minutes}$**, and refresh tokens shall have a maximum TTL of **$\le 7\text{ days}$**. | Automated unit test suite verifying token expiration enforcement past 901 seconds. |
| **NFR-07** | Security | The system shall prevent Insecure Direct Object References (IDOR); cross-student grade retrieval or unauthorized data access shall yield **0 data leaks** and return HTTP 403 Forbidden. | Automated security integration test suite executing unauthorized cross-tenant requests. |
| **NFR-08** | Usability | The web interface shall display **0 horizontal scrollbars or broken layouts** across responsive breakpoints: $375\text{px}$ (Mobile), $768\text{px}$ (Tablet), and $1280\text{px}+$ (Desktop). | Visual regression testing using Playwright across configured viewport matrix. |
| **NFR-09** | Usability | Under simulated 4G mobile network conditions, the web application shall achieve First Contentful Paint (FCP) **$\le 1.5\text{ seconds}$**, with an initial production JS+CSS bundle **$\le 450\text{ KB}$** (gzipped). | Automated Google Lighthouse CI audit running against production frontend build. |
| **NFR-10** | Usability / Accessibility | The portal shall achieve **WCAG 2.1 AA compliance** with text color contrast $\ge 4.5:1$, 100% keyboard operability, visible focus indicators, and an automated Lighthouse Accessibility score **$\ge 95/100$**. | Automated Axe-core audit, Lighthouse Accessibility scan, and manual NVDA/ChromeVox screen reader evaluation. |
| **NFR-11** | Reliability | The system shall ensure **100% ACID compliance** on grade finalization and course registration write-behind persistence queues. | Integration tests simulating transient database failures during transaction processing, asserting zero corrupted partial states. |
| **NFR-12** | Observability | The `/health` endpoint of each microservice shall respond with HTTP status 200 in **$\le 50\text{ ms}$**. | Continuous curl health-check polling every 10 seconds across all 5 backend microservices during staging runtime. |
| **NFR-13** | Maintainability | The core backend business domain packages shall maintain a statement test coverage of **$\ge 70\%$**. | Go test coverage analysis: `go test -coverprofile=coverage.out ./...`. |
| **NFR-14** | Portability | The entire multi-service application environment shall boot from a clean git checkout to a fully reachable operational state via `docker compose up -d` in **$\le 60\text{ seconds}$**. | Timed clean deployment script on a standard dual-core development machine with pre-cached Docker images. |
| **NFR-15** | Architecture | Flash-crowd traffic spikes saturating the `registration-service` (100% CPU) shall result in **0% cascade failure** on grade inspection and document submission services. | Chaos engineering test: running full load against registration while simultaneously verifying HTTP 200 on `/api/grades/mine`. |

---

## 5.3 Domain Requirements and Business Rules

*[Identify inviolable business policies, institutional regulations, academic bylaws, and legal standards governing system logic.]*

| ID | Rule | Source (policy, law, standard) |
| :--- | :--- | :--- |
| **BR-01** | **Active Registration Window Required:** Section registration, modification, and dropping are permitted strictly when the registration period assigned to the student's program track (`STANDARD` vs `PROJECT_BASED`) and cohort is in the operational status `OPEN`. | HCMUS Academic Affairs Policy, Decision No. 128/QĐ-KHTN on Course Registration |
| **BR-02** | **No Duplicate Course Code Enrollment:** A student cannot hold simultaneous active enrollments in multiple sections belonging to the same root course (`course_id`) within the same academic semester, even under different section codes. | MoET Circular 08/2021/TT-BGDĐT on Credit-Based Higher Education, Art. 5 |
| **BR-03** | **Advisory Non-Blocking Timetable Overlaps:** Time slot collisions between newly selected and enrolled sections generate prominent UI warnings (`CONFLICT_DETECTED`) but do not programmatically block registration submission. | HCMUS Academic Advisor Guidance Policy & ADR-004 |
| **BR-04** | **Hard Section Capacity Enforcement:** Total enrolled students in any section must never exceed `section.max_capacity`. Any concurrent registration attempting to exceed capacity must be rejected immediately. | FIT Classroom & Laboratory Safety Capacity Regulations |
| **BR-05** | **Decoupled Schedule Planning:** Saved schedule plans constructed in the Schedule Planner represent draft advisory combinations; they confer no reserved seat rights or official enrollment claims until officially submitted and committed. | Portal67 Solution Architecture v3.0 Specification |
| **BR-06** | **Program Track Cohort Isolation:** Standard track students (`STANDARD`) cannot register Project-Based sections (`PROJECT_BASED`). Project-Based students may only enroll in Standard sections when `is_cross_enrollment_allowed = TRUE` under formal academic petition. | VP Đề án (CTĐA) Training Scheme Regulations & ADR-008 |
| **BR-07** | **Integrated Theory-Practice Enrollment:** For courses with laboratory components (`has_lab = TRUE`), each course section encapsulates both Theory and Practice lecture/lab schedules. Registering for the section automatically enrolls and schedules both Theory and Practice sessions for that single section without separate or unbundled registration. Selecting a section commits the student to both required schedules. | FIT Laboratory Training & Course Syllabus Policy & ADR-008 |
| **BR-08** | **In-Memory Atomic Seat Reservation:** During peak registration periods, seat quotas are reserved atomically via Redis Lua scripts (`seat_reserve.lua`) prior to asynchronous relational persistence, guaranteeing 0 over-allocations without database row lock contention. | High-Concurrency System Design Standards & ADR-007 |
| **BR-09** | **Standard 10-Point Component Grading Scale:** All continuous and final grade components must be evaluated on a scale of $0.0$ to $10.0$ with a precision of one decimal place. | MoET Circular 08/2021/TT-BGDĐT, Art. 9 (Academic Evaluation) |
| **BR-10** | **Dynamic Assessment Rubric Weighting:** A course's grade composition is defined dynamically via $k$ grade components (`academic.grade_components`). The sum of all active component weights must equal exactly $100\%$ ($1.00$) prior to grade sheet submission. | FIT Academic Evaluation Framework & CHG-05 (BCNF Normalization) |
| **BR-11** | **Semester and Cumulative GPA Computation:** Semester GPA is calculated as the credit-weighted average of course grades to two decimal places. Cumulative GPA considers only the highest earned grade for retaken courses according to faculty credit regulations. | HCMUS Undergraduate Academic Regulations & Credit Training Policy |
| **BR-12** | **Separation of Grading Duties (Four-Eyes Principle):** Lecturers hold authority only to input and submit section grade sheets (`Draft` $\rightarrow$ `Submitted`). Only authorized Academic Affairs personnel (`ADM`) hold authority to approve, publish, and finalize official grades (`Submitted` $\rightarrow$ `Finalized`). | FIT Department Grading Governance Protocol & Four-Eyes Principle |
| **BR-13** | **Immutability of Finalized Grades & Audit Trail:** A finalized grade sheet cannot be modified via standard update APIs. Any post-finalization grade amendment requires a formal rectification record (`Grade_Correction_Log`) specifying prior score, new score, justification reason, and approving officer ID. | ISO/IEC 27001 Audit Trail Standards & MoET Inspection Protocol |
| **BR-14** | **Multi-Lecturer Assignment Permissions:** In co-taught sections (`academic.section_lecturers`), co-professors and TAs may enter and update continuous component scores; however, only the designated Primary Lecturer (`is_primary = TRUE`) holds permission to lock and submit the completed grade sheet. | FIT Faculty Teaching Assignment Regulations |
| **BR-15** | **Immutable Request State Machine:** Administrative requests adhere to a unidirectional progression: `Submitted` $\rightarrow$ `Under_Review` $\rightarrow$ `Approved` / `Rejected` $\rightarrow$ `Completed`. Requests cannot be deleted from the audit ledger once submitted. | HCMUS Student Affairs (Phòng CTSV) Administrative Procedure Standard |
| **BR-16** | **Information Revision Cycle:** If an administrative officer marks a request as `Requires_Information`, the request transitions back to a student-editable state while preserving all historical audit notes and timestamps. | Phòng CTSV Administrative Petition Protocol |
| **BR-17** | **Configurable Grade Appeal Window (*Phúc khảo*):** Grade appeal applications are accepted strictly within the review/appeal window configured and opened by Academic Affairs staff (`appeal_start_date` to `appeal_end_date`) following official grade finalization. Academic Affairs retains full authority to set and adjust the appeal duration per term or department. | HCMUS Student Handbook & Examination Regulations (*Quy chế Khảo thí*) |
| **BR-18** | **Immutable Student Identifier:** The student ID (MSSV) serves as the persistent, immutable institutional key within the academic domain. An account's email may be updated, but MSSV cannot be changed. | VNU-HCM Student Identification Standards |
| **BR-19** | **Protected Academic Profile Fields:** Students are forbidden from directly modifying their core academic attributes (MSSV, enrolled program track, cohort, admitted major, cumulative credits, or academic status). | HCMUS Academic Affairs Database Integrity Policy |
| **BR-20** | **Pre-Provisioned Accounts & No Public Self-Registration:** As an institutional university portal, Portal67 strictly prohibits public self-registration. All student, lecturer, and administrative staff accounts are pre-provisioned by university administration or synchronized via institutional directories. Authentication requires institutional identifier (MSSV / MSCB) and password, or authorized institutional OAuth/SSO. | HCMUS Information Security & User Identity Administration Policy |

---

## 5.4 Acceptance Criteria

*[For each Must-have feature, write acceptance criteria in Given / When / Then form. They will become acceptance test cases in PA4.]*

---

### Feature F01: User Authentication & Access Control (`FR-01`, `FR-03`)

#### Scenario AC-F01-01: Successful Authentication and Scoped JWT Token Issuance (Happy Path)
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

#### Scenario AC-F01-02: Backend Authorization Guard Rejects Unauthorized Access (Edge Case)
```gherkin
Scenario: Student attempts to access administrative endpoint
  Given A student is authenticated with role "STUDENT" and holds a valid Bearer token
  When The student sends a GET request to "/api/affairs/requests/queue"
  Then The API gateway rejects the request with HTTP status 403 Forbidden
    And Returns an error payload: '{"code": "FORBIDDEN", "message": "Insufficient permissions"}'
    And An audit security log entry is recorded with actor ID, IP address, and timestamp
```

#### Scenario AC-F01-03: Successful Authentication via Institutional OAuth/SSO
```gherkin
Scenario: Pre-provisioned user authenticates via Institutional OAuth/SSO
  Given An institutional user account "24127003" is pre-provisioned in Portal67 with role "STUDENT"
    And The user authenticates successfully via university OAuth identity provider with email "24127003@student.hcmus.edu.vn"
  When The OAuth callback returns authorization code to "/api/auth/oauth/callback"
  Then The system verifies the token, resolves identity to student MSSV "24127003"
    And Issues an Access Token (expiry <= 15 minutes) and Refresh Token (expiry <= 7 days)
```

---

### Feature F02: Student Overview & Academic Profile (`FR-04`)

#### Scenario AC-F02-01: Student Dashboard Overview Rendered with Consolidated Academic Data
```gherkin
Scenario: Student views dashboard upon successful login
  Given Student "24127003" is authenticated and navigates to "/student/dashboard"
    And The student has 3 enrolled classes for the current semester
    And The student has 1 pending document request "REQ-2026-001"
  When The dashboard page completes loading
  Then The interface displays today's scheduled class sessions with room codes and times
    And Displays the student's cumulative GPA (e.g., "8.75 / 10.0") and total credits earned (e.g., "68")
    And Displays the active document request status badge as "Under_Review"
    And Loads recent faculty markdown announcements in the notification pane
```

---

### Feature F03: Course & Section Registration (`FR-06`, `FR-07`, `FR-08`, `FR-09`)

#### Scenario AC-F03-01: Successful Section Enrollment with Integrated Theory and Practice Sessions (Happy Path)
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
    And Both Theory and Practice sessions appear automatically on the student's personal weekly timetable
```

#### Scenario AC-F03-02: Registration Rejected When Section Capacity is Exhausted (0 Over-allocations)
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

#### Scenario AC-F03-03: Timetable Conflict Advisory Warning Displayed (ADR-004)
```gherkin
Scenario: Student selects an overlapping section without blocking enrollment
  Given The student is already enrolled in section "MTH00001_01" scheduled on "Tuesday 07:30 - 09:10"
    And The student selects section "PHY00002_03" scheduled on "Tuesday 08:20 - 10:00"
  When The student views the registration selection or schedule planner
  Then The system flags a "Schedule Conflict Advisory" for Tuesday Periods 2-3
    And The conflict is visually highlighted in amber with an explanatory warning
    And The "Submit Registration" action remains enabled and accessible
```

#### Scenario AC-F03-04: Cross-Track Cohort Isolation Enforced (ADR-008)
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

### Feature F04: Timetable Visualisation & Scheduling (`FR-13`)

#### Scenario AC-F04-01: Weekly Timetable Displays Theory and Practice Sessions with Facility Details
```gherkin
Scenario: Student inspects personalized weekly timetable for active term
  Given Student "24127003" is enrolled in section "CSC10002_01"
    And "CSC10002_01" has Theory scheduled on "Monday 07:30 - 10:00" at "CS2_C32"
    And "CSC10002_01" has Practice scheduled on "Wednesday 13:00 - 16:30" at "CS2_Lab04"
  When The student navigates to "/student/timetable" and selects week 5
  Then The timetable grid displays the Monday lecture block colored by subject
    And Displays the Wednesday lab block with distinct lab indicator icon
    And Each timetable cell indicates Course Title, Section Code, Room, Campus, and Lecturer Name
```

---

### Feature F05: Gradebook & Dynamic Evaluation (`FR-17`, `FR-18`, `FR-19`)

#### Scenario AC-F05-01: Primary Lecturer Submits Dynamic Rubric Grade Sheet (Happy Path)
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

#### Scenario AC-F05-02: Academic Affairs Finalizes Grades and Recalculates Cumulative GPA
```gherkin
Scenario: Academic Affairs approves submitted grade sheet and publishes official grades
  Given A grade sheet for section "CSC10002_01" exists in "Submitted" state
    And An officer with role "Academic Affairs" is logged into the management portal
  When The officer clicks "Approve & Finalize Grade Sheet"
  Then The grade sheet state transitions to "Finalized"
    And The official grades immediately become visible to all enrolled students via their Grade Portal
    And Cumulative GPA and completed credits for affected students are automatically recalculated
```

#### Scenario AC-F05-03: Direct Mutation of Finalized Grade Rejected via API (Edge Case)
```gherkin
Scenario: User attempts to edit a finalized grade directly via API
  Given Grade for student "24127003" in section "CSC10002_01" is in "Finalized" state
  When A PUT request is dispatched to "/api/grades/sections/CSC10002_01/students/24127003"
  Then The backend rejects the request with HTTP 403 Forbidden
    And Returns error message: "Finalized grades cannot be directly edited. Submit a Grade Rectification Request."
```

#### Scenario AC-F05-04: Teaching Assistant Restricted from Submitting Official Grade Sheet
```gherkin
Scenario: Teaching Assistant attempts to finalize or submit grade sheet
  Given User "TA001" is assigned to section "CSC10002_01" with is_primary = FALSE
  When User "TA001" dispatches a POST request to "/api/grades/sections/CSC10002_01/submit"
  Then The backend rejects the request with HTTP 403 Forbidden
    And Returns error message: "Only the designated Primary Lecturer can submit section grade sheets"
```

---

### Feature F06: Student Affairs & Document Workflows (`FR-21`, `FR-22`, `FR-23`)

#### Scenario AC-F06-01: Student Submits Administrative Document Request
```gherkin
Scenario: Student submits an administrative document request successfully
  Given Student "24127003" is authenticated and accesses the Document Request portal
  When The student submits a request of type "ENROLLMENT_VERIFICATION" with reason "Military deferment"
  Then A new request record is created with a unique tracking ID (e.g., "REQ-2026-0042")
    And Initial status is set to "Submitted"
    And The request appears in the student's "My Requests" tracking list
    And The request appears in the Student Affairs (CTSV) pending queue
```

#### Scenario AC-F06-02: Administrative Staff Reviews and Approves Document Request
```gherkin
Scenario: CTSV officer approves student document request and marks ready for pickup
  Given A document request "REQ-2026-0042" exists with status "Under_Review"
    And A CTSV officer inspects the attached military board call notice
  When The officer clicks "Approve Request" and specifies pickup date "2026-11-05"
  Then The request status transitions to "Approved"
    And An automated in-app notification is dispatched to student "24127003"
    And The audit history records approver ID and timestamp
```

---

### Feature F07: Academic Operations & Term Management (`FR-25`, `FR-26`)

#### Scenario AC-F07-01: Academic Staff Opens Track-Specific Registration Window
```gherkin
Scenario: Phòng Đào tạo opens course registration for Standard Track
  Given An officer with role "Academic Affairs" accesses "/admin/registration-windows"
  When The officer creates a registration window for semester "2026-1", program "STANDARD", cohort "2024"
    And Sets start time to "2026-11-01 08:00" and end time to "2026-11-07 17:00"
    And Sets status to "OPEN"
  Then The system activates the registration window in `registration.registration_windows`
    And Eligible Standard students are granted access to enroll in sections
    And Project-Based students attempting to register under this window are excluded
```

---

### Feature F08: Asynchronous Notifications & System Alerts (`FR-28`, `FR-29`)

#### Scenario AC-F08-01: Asynchronous Dispatch and Rendering of Grade Finalization Notification
```gherkin
Scenario: System delivers real-time notification to enrolled students when grades are published
  Given Section "CSC10002_01" grade sheet transitions to "Finalized"
  When The background notification worker consumes the "GradeFinalizedEvent" from Redis Streams
  Then An in-app alert is created in `affairs.notifications` for all enrolled students
    And When student "24127003" refreshes or opens the portal
    And A bell notification badge appears with message: "Official grades for CSC10002 have been published"
    And Clicking the notification routes directly to "/student/grades"
```

---
*End of Section 5 Requirements Specification — Team 02 (BDD)*
