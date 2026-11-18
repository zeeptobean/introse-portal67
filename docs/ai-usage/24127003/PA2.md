# AI Usage Log — PA2: Requirements Specification

* **Member Name:** Vũ Trần Minh Hiếu  
* **Student ID:** 24127003  
* **Role:** System Architect / UI Designer  
* **Phase Gate:** PA2 — Requirements Specification  
* **Date:** 2026-10-28  
* **Team:** Team 02 — BDD (NMCNPM, HK1 2026–2027)  

---

## Session Details

### Task: Requirements Quality Standardization (FR, NFR, Domain Rules, Acceptance Criteria, RTM)
* **Tool Used:** Antigravity (Gemini 3.8 Flash High)  
* **Objective:** Formulate a formal Software Requirements Specification (SRS) for PA2 adhering strictly to NMCNPM rubric requirements:
  1. Functional Requirements (FR) numbered, grouped, and prioritised (MoSCoW).
  2. Non-Functional Requirements (NFR) strictly measurable with quantitative metrics and thresholds (ISO/IEC 25010).
  3. Formal Domain Rules (DR) modeling HCMUS academic policies and separation of duties.
  4. Behavior-Driven Acceptance Criteria (AC) formulated in Gherkin (Given-When-Then) syntax.
  5. Bidirectional Requirements Traceability Matrix (RTM).

* **Input Prompt:**
  ```text
  đọc file docs/design/Portal67_Solution_Architecture.md và giúp tôi làm phần này: 
  Requirements quality: FR numbered, grouped, prioritised; NFR measurable; domain rules; acceptance criteria
  giải thích cần phải làm gì và làm như nào
  ```

* **AI Output Summary:**
  - Analyzed gaps in existing documentation and outlined the requirements engineering methodology.
  - Authored [`docs/requirements/PA2_Requirements_Specification.md`](../../requirements/PA2_Requirements_Specification.md) comprising:
    - **28 Functional Requirements (`FR-AUTH-01` to `FR-NOTIF-02`)** partitioned across 8 subsystems (`AUTH`, `STU`, `REG`, `TT`, `GRD`, `REQ`, `ADM`, `NOTIF`), prioritized using MoSCoW.
    - **13 Quantified Non-Functional Requirements** with specific thresholds: API read latency ($p95 \le 300\text{ms}$), peak registration latency ($p95 \le 600\text{ms}$), throughput ($\ge 100\text{ RPS}$), 0 over-subscription race conditions, Argon2id password hashing, $15\text{m}$ JWT token lifetime, backend unit test coverage $\ge 70\%$, Docker cold-start $\le 60\text{s}$.
    - **14 Domain Rules (`DR-REG`, `DR-GRD`, `DR-REQ`, `DR-ID`)** covering HCMUS registration windows, advisory timetable overlap policy, 10-point and 4-point GPA calculations, four-eyes grade sheet submission/finalization, and request state transitions.
    - **BDD Acceptance Scenarios** covering both happy paths and edge cases (class section full, timetable overlap warning, grade direct mutation rejection).
    - **Traceability Matrix (RTM)** establishing forward/backward links between requirements, domain rules, NFRs, and acceptance tests.

* **Human Review & Modifications (Critical):**
  - **Strengths:** Accurately captured university-specific domain nuances (such as non-blocking timetable conflict warnings reflecting HCMUS student practice, and strict separation between lecturer grade entry vs. academic affairs approval).
  - **Limitations / Errors Detected:** The initial draft lacked an explicit verification test condition for concurrent seat reservation race conditions.
  - **Modifications Applied:** Explicitly defined `NFR-SCAL-02` requiring 0 over-allocations when 20 concurrent threads compete for 1 final seat; validated that all requirements map to the MVP scope (Weeks 8–10).

* **Verification Steps:**
  - [x] Verified file created at [`docs/requirements/PA2_Requirements_Specification.md`](../../requirements/PA2_Requirements_Specification.md).
  - [x] Cross-checked against Portal67 Solution Architecture baseline and AGENT.md guidelines.
  - [x] Commit reference / PR link: `docs(requirements): create PA2 SRS baseline with FR, NFR, DR, and AC`


---

### Task 2: Requirements Synchronization with PA3 Architectural Design & Solution Architecture v3.0
* **Date:** 2026-10-31  
* **Tool Used:** Antigravity (Gemini 3.8 Flash High)  
* **Objective:** Synchronize the PA2 Requirements Baseline with PA3 System Design Specification v3.0 and Portal67 Solution Architecture v3.0, incorporating key architectural decisions:
  1. HCMUS Dual-Track Academic Model (Standard 2 semesters vs. Project-Based 3 trimesters).
  2. Theory-Practice Section Hierarchy with atomic dual-seat reservation.
  3. High-concurrency `registration-service` engine utilizing in-memory Redis Lua scripts.
  4. Dynamic assessment rubrics achieving strict BCNF normalization.
  5. Multi-lecturer co-teaching and Teaching Assistant (TA) role permissions.
  6. Grade appeal (*Phúc khảo*) state machine.
  7. Requirements Evolution Log (`CHG-01` to `CHG-07`) and expanded bidirectional RTM.

* **Input Prompt:**
  ```text
  đọc toàn bộ những thay đổi trong docs/design/PA3_System_Design_Specification.md và docs/design/Portal67_Solution_Architecture.md và giúp tôi cập nhật cho docs/requirements/PA2_Requirements_Specification.md
  ```

* **AI Output Summary:**
  - Upgraded [`docs/requirements/PA2_Requirements_Specification.md`](../../requirements/PA2_Requirements_Specification.md) to **Document Version 3.0**:
    - **Actor Model Refinement (Section 1.2):** Explicitly separated Student profiles (`STANDARD` vs `PROJECT_BASED`), Lecturer roles (Primary Lecturer vs TA), and Administrative divisions (Phòng Đào tạo vs VP Đề án vs CTSV).
    - **Catalog Expansion (Section 2):** Refined `FR-REG-01`, `FR-REG-02`, `FR-REG-03`, `FR-GRD-03`, `FR-ADM-01`, `FR-ADM-02`, and added `FR-REG-07` (Cross-Registration Petitions).
    - **NFR Tightening (Section 3):** Added `NFR-USAB-03` (WCAG 2.1 AA Accessibility checklist) and `NFR-ARCH-01` (Microservice Fault Isolation), with sharpened registration write latency ($p95 \le 300\text{ms}$, Redis Lua $\le 10\text{ms}$).
    - **Domain Rules (Section 4):** Added `DR-REG-06` (Cohort Isolation), `DR-REG-07` (Integrated Theory-Practice Enrollment), `DR-REG-08` (Atomic Lua Reservation), and `DR-GRD-06` (Multi-Lecturer Submission Guard).
    - **BDD Scenarios (Section 5):** Added acceptance tests for Integrated Theory-Practice enrollment (`AC-REG-01`), Cross-track enrollment blocking (`AC-REG-04`), TA submission rejection (`AC-GRD-04`), and Configurable Grade Appeal (`AC-REQ-02`).
    - **Bidirectional RTM (Section 6):** Mapped all requirements to target microservices, API routes, database schemas, and NFRs.
    - **Requirements Evolution Log (Section 7):** Formally documented `CHG-01` through `CHG-07`.
  - Authored standalone report deliverable [`docs/requirements/PA2_Section_5_Requirements.md`](../../requirements/PA2_Section_5_Requirements.md) strictly formatted to the official HCMUS SE Dept. SRS Template:
    - **5.1 Functional Requirements:** All statements formatted as *"The system shall ..."*, partitioned into features `F01`–`F08`, MoSCoW prioritized, and attributed to institutional sources.
    - **5.2 Non-Functional Requirements:** Verifiable conditions with concrete quantitative thresholds and verification methods.
    - **5.3 Domain Requirements and Business Rules:** Formatted into `BR-01`–`BR-19` with statutory and policy citations.
    - **5.4 Acceptance Criteria:** BDD Given/When/Then scenarios for all 8 Must-have features (`F01`–`F08`) prepared for PA4 test suites.

* **Verification Steps:**
  - [x] Verified file updated at [`docs/requirements/PA2_Requirements_Specification.md`](../../requirements/PA2_Requirements_Specification.md).
  - [x] Verified standalone deliverable created at [`docs/requirements/PA2_Section_5_Requirements.md`](../../requirements/PA2_Section_5_Requirements.md).
  - [x] Verified complete alignment with `docs/design/PA3_System_Design_Specification.md` and `docs/design/Portal67_Solution_Architecture.md`.
