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
