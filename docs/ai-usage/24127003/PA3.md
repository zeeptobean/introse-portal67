# AI Usage Log — PA3: System & Architectural Design

* **Member Name:** Vũ Trần Minh Hiếu  
* **Student ID:** 24127003  
* **Role:** System Architect / UI Designer  
* **Phase Gate:** PA3 — System & Architectural Design  
* **Date:** 2026-10-28  
* **Team:** Team 02 — BDD (NMCNPM, HK1 2026–2027)  

---

## Session Details

### Task: PA3 Architectural & System Design Specification (100% Rubric Compliance)
* **Tool Used:** Antigravity (Gemini 3.8 Flash High)  
* **Objective:** Analyze the official PA3 grading rubric provided via image and construct a complete, audit-proof architectural design document achieving the Excellent level (90%–100%) across all 5 criteria:
  1. Architecture (30%): Justified style, C4 Context, Container, Component, and Deployment diagrams, $\ge 3$ ADRs.
  2. Object-Oriented Design (30%): Consistent Class, Sequence, and State diagrams; Class specifications; principles/patterns justified.
  3. Data Design (20%): Correct ERD, table specifications with integrity constraints, normalization (1NF–3NF/BCNF), indexes.
  4. UI/UX Design (10%): Screen flow, specifications of key screens, accessibility checklist (WCAG 2.1 AA).
  5. Consistency and Traceability (10%): Realization of PA2 requirements with bidirectional matrix and requirements change log.

* **Input Prompt:**
  ```text
  tạo các file markdown cho cái trên, đồng thời xem, phân tích và tạo luôn file markdown cho ảnh này sau khi phân tích
  [Uploaded Image: Official PA3 Grading Rubric with 5 criteria totaling 100%]
  ```

* **AI Output Summary:**
  - Performed an architectural gap analysis against existing project files.
  - Authored [`docs/design/PA3_System_Design_Specification.md`](../../design/PA3_System_Design_Specification.md) detailing:
    - **Architecture (30%):** Style evaluation matrix comparing Monolith vs. Pragmatic Microservices vs. Distributed Microservices; C4 Level 1 (Context) and Level 2 (Container); C4 Level 3 (Component) for Academic Service; Docker Compose physical deployment topology; and summary of 6 ADRs.
    - **OO Design (30%):** Comprehensive Domain Class Diagram; State Machine diagrams for Grade Sheet and Request Lifecycles; Concurrency-safe Registration Sequence Diagram (`SELECT FOR UPDATE`); Go domain service and repository interface definitions; and justification of SOLID principles, Repository Pattern, and Asynchronous Event Dispatch.
    - **Data Design (20%):** Relational Crow's Foot ERD; full table specifications for `academic.sections` and `academic.grade_entries` with constraints (`CHECK`, `UNIQUE`, defaults); normalization mathematical proof (1NF $\rightarrow$ 3NF/BCNF) with intentional denormalization justification; and SQL indexing strategy.
    - **UI/UX Design (10%):** Complete Mermaid Screen Flow diagram; UX component specifications for Course Registration and Lecturer Grade Entry; and WCAG 2.1 AA compliance checklist.
    - **Consistency & Traceability (10%):** End-to-end matrix mapping all PA2 requirements to PA3 microservices, API routes, database tables, and class operations; plus Requirements Evolution Log (`CHG-01` to `CHG-04`).

* **Human Review & Modifications (Critical):**
  - **Strengths:** Fully utilized GitHub-flavored Markdown and valid Mermaid diagrams (avoiding unsupported flowchart features or invalid syntax); seamlessly linked PA2 requirements to concrete database tables and Go interfaces.
  - **Limitations / Errors Detected:** 
    - Multi-tenancy and database separation between microservices required explicit representation in SQL table specifications.
    - C4 Level 1 (System Context) initially lacked detailed actor separation (missing explicit CTSV staff persona) and external systems (missing Payment Gateway VietQR, Object Storage MinIO), as well as tabular specifications and boundary definitions required for full rubric compliance.
    - Initial grade design (`grade_entries`) used hard-coded columns (`attendance`, `midterm`, `coursework`, `final_exam`), which caused non-superkey functional dependency violations ($\{components\} \rightarrow overall\_score$) against strict BCNF, while lacking flexibility for courses with non-standard grading rubrics.
  - **Modifications Applied:** 
    - Ensured all SQL table definitions explicitly declare schema namespaces (`auth`, `academic`, `affairs`) corresponding to microservices database boundaries.
    - Upgraded Section 1.2 with an expanded C4 Level 1 Context diagram (4 user actors aligned strictly with PA2 baseline: Student, Lecturer, Academic & Administrative Staff / Giáo vụ, and System Administrator; and 4 external systems), System Actors & User Personas Matrix, External System Integration Specification table, System Scope Boundary (In-Scope vs. Out-of-Scope), and Architectural Context Principles.
    - Fixed Mermaid C4 syntax compatibility: replaced unsupported `Enterprise_System(...)` with standard `System(...)` macro in C4Context; wrapped `Container_Boundary(...)` with `{ ... }` curly braces (replacing invalid `Boundary_End()`) in C4Component; and corrected deployment header and node declarations to `C4Deployment` and `Deployment_Node(...)`.
    - Localized 100% of the specification document into English to maintain consistent international academic reporting standards.
    - Refactored Data Design to **100% BCNF compliance**: decoupled static grade sheets into `academic.grade_components` (dynamic assessment criteria, weights, descriptions) and `academic.student_grades` (per-student score entries tied to `academic.enrollments`); extracted multi-meeting weekly timetables into `academic.section_schedules`; updated domain class diagrams, Go interfaces, ERD, and traceability matrix (`CHG-05`).
    - **Decoupled Course Registration into Standalone Microservice (`registration-service`):** Addressed portal crashing vulnerabilities during high-concurrency peak registration windows (e.g. 8:00 AM) by separating course registration from `academic-service` into an independent high-concurrency engine (Port 8085). Refactored C4 Container, C4 Component (Registration Service), C4 Deployment, High-Concurrency Sequence Diagram, Go interfaces, and RTM to utilize sub-2ms in-memory Redis Lua atomic seat reservations and asynchronous write-behind persistence queues (`CHG-06`, `ADR-007`).
    - **Added Grade Appeal & Score Rectification State Machine:** Enhanced Section 2.2 with a dedicated state machine for Grade Appeals (`Phúc khảo` — `FR-REQ-04`, `FR-GRD-06`, `DR-REQ-03`, `DR-GRD-05`), formalizing the statutory 7-day window, departmental re-assessment, staff approval workflow, immutable `Grade_Correction_Log` audit trail, and automatic GPA recalculation.

* **Verification Steps:**
  - [x] Verified file created at [`docs/design/PA3_System_Design_Specification.md`](../../design/PA3_System_Design_Specification.md).
  - [x] Confirmed all 5 criteria from the rubric image are thoroughly documented.
  - [x] Commit reference / PR link: `docs(design): create comprehensive PA3 architecture and system design specification`
