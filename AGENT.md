# AGENT.md - Developer & AI Agent Guide for Portal67

This document serves as the authoritative source of truth for AI Agents (Copilot, Claude, Gemini, Antigravity, etc.) and human developers working on the **Portal67** project by **Team 02 (BDD)**.

---

## 1. Project Overview & Metadata

* **Project Name:** Portal67
* **Pitch / Vision:** A modern, high-performance, and unified academic and administrative web portal that resolves key bottlenecks and enhances the student and faculty experience at HCMUS.
* **Course:** Introduction to Software Engineering
* **Department & University:** Software Engineering Department, Faculty of Information Technology, University of Science, VNU-HCM (HCMUS).
* **Academic Term:** Semester 1 (HK1), Academic Year 2026 - 2027.
* **Team Name:** Team 02 - **BDD**
* **Repository:** `https://github.com/zeeptobean/introse-portal67.git`
* **Communication Channels:** Discord (primary chat & team sync) & Messenger
* **Task Board:** GitHub Projects
* **Target Grade & Commitment:** 9.0+ score; each team member commits a minimum of 6 hours/week.

### Team Roster & Role Assignments

| Student ID | Full Name | Academic Email | Primary Role | Secondary (Backup) Role |
| :--- | :--- | :--- | :--- | :--- |
| **24127326** | **Đoàn Quốc Bảo** | `dqbao2432@clc.fitus.edu.vn` | **Team Lead / Project Manager** | Developer |
| **24127025** | **Lê Minh Đức** | `lmduc2433@clc.fitus.edu.vn` | **Business Analyst (BA)** | UI/UX Designer |
| **24127003** | **Vũ Trần Minh Hiếu** | `vtmhieu2416@clc.fitus.edu.vn` | **System Architect / UI Designer** | Developer |
| **24127173** | **Nguyễn Tuấn Hùng** | `nthung2432@clc.fitus.edu.vn` | **Lead Developer** | QA / Tester |
| **24127008** | **Trương Nhật Phát** | `tnphat2418@clc.fitus.edu.vn` | **QA / Tester** | PM / Business Analyst |

---

## 2. Background, Motivation & System Goals

### 2.1. Problem Statement
The current university portal systems and administrative workflows face noticeable pain points:
1. **Server Overload & Concurrency Bottlenecks:** During peak course registration windows, the portal often suffers severe slowdowns and downtime.
2. **Fragmented Workflows:** Academic services are scattered across different platforms (separate systems for Regular vs. High-Quality/Advanced programs, Moodle, Google Classroom, and third-party Google Forms for requests).
3. **Lack of Integrated Utility Tools:** Critical features such as GPA forecasting, intelligent timetable conflict detection, and centralized paper request tracking are missing or rely on unofficial browser extensions (e.g., KHTN GPA 2.0).

### 2.2. Solution Objectives
* Develop a unified, responsive, high-performance web portal accessible by **Students**, **Teachers (Lecturers)**, and **Administrative Staff**.
* Build a reliable working prototype within the 3-week implementation phase (Weeks 8–10).
* *Note:* The system will operate on **sample/seeded data** and does not connect to live university production databases. It will not process real money transfers or issue legally binding physical documents.

---

## 3. Scope & Key Features (MoSCoW Prioritization)

| ID | Feature | Description | Priority | Target Delivery |
| :--- | :--- | :--- | :--- | :--- |
| **F01** | **Course & Specialized Topic Registration** | Browse open courses and elective/specialized topics, inspect live seat availability and schedules, submit or drop enrollments. Validates seat limits and schedule conflicts using sample data. | **Must** | MVP (Week 8–10) |
| **F02** | **Tuition Management** | View semester tuition breakdown, billing records, remaining balances, and payment statuses. Administrative staff can update simulated payment entries (no real transactions). | **Must** | MVP (Week 8–10) |
| **F03** | **Academic Result Management** | Lecturers enter demo grades; students view grades and real-time semester and cumulative GPA estimations based on predefined rules. | **Must** | MVP (Week 8–10) |
| **F04** | **Course & Exam Scheduling** | Display weekly student class timetables and examination dates; detect and visually highlight overlapping classes or conflicting exam slots. | **Must** | MVP (Week 8–10) |
| **F05** | **Resource Management** | Long-term classroom, laboratory, and equipment reservation for faculty and staff. | **Could** | Deferred beyond MVP (only implemented if core work finishes early) |
| **F06** | **Document Request & Approval** | Students submit paper/certificate requests; administrative workers review, approve, or reject with feedback; students track real-time status and audit history. | **Must** | MVP (Week 8–10) |
| **F07** | **Role-Based Access Control (RBAC)** | Strict access differentiation across roles: Students, Lecturers, and Administrative Staff. | **Must** | MVP (Week 8–10) |
| **F08** | **Input Validation & Error Handling** | Comprehensive validation to reject invalid registrations, missing fields, or unauthorized grade updates with informative user error messaging. | **Must** | MVP (Week 8–10) |

---

## 4. System Architecture & Tech Stack

### 4.1. Preliminary C4 System Context Architecture

```
                               ┌─────────────────────────┐
                               │       Student           │
                               │ (Undergrad / Postgrad)  │
                               └───────────┬─────────────┘
                                           │ Registers courses, checks GPA,
                                           │ timetable, tuition, paper requests
                                           ▼
┌─────────────────────────┐    ┌─────────────────────────┐    ┌─────────────────────────┐
│  Administrative Worker  ├───►│        Portal67         │◄───┤         Teacher         │
│(Faculty / Dept Staff)   │    │    (HCMUS Web Portal)   │    │       (Lecturers)       │
│ Approves paper requests │    └──────┬─────┬─────┬──────┘    │ Enters demo grades,     │
│ & updates tuition status│           │     │     │           │ manages course sections │
└─────────────────────────┘           │     │     │           └─────────────────────────┘
                                      │     │     │
                 Authenticates SSO    │     │     │ Sends simulated payments
                                      ▼     │     ▼
        ┌───────────────────────────────┐   │   ┌───────────────────────────────┐
        │     Authentication Service    │   │   │        Payment Gateway        │
        │    (External SSO Provider)    │   │   │ (Simulated payment processor) │
        └───────────────────────────────┘   │   └───────────────────────────────┘
                                            │
                        Sends notifications │
                                            ▼
                                ┌───────────────────────────────┐
                                │         Email Service         │
                                │(Automated notifications dispatch)
                                └───────────────────────────────┘
```

### 4.2. Technology Stack Selection

* **Frontend:**
  * **Framework:** **Svelte / SvelteKit** — Selected for minimal runtime overhead, zero virtual DOM overhead, fine-grained reactivity, and rapid load times.
  * **Styling:** **Tailwind CSS** — Utility-first styling for responsive design across mobile, laptop, and desktop viewports.
* **Backend:**
  * **Language & Runtime:** **Golang (Go)** — Selected for outstanding execution performance, high concurrency support (goroutines/channels) under heavy registration traffic, memory efficiency, and robust standard library.
  * **API Architecture:** RESTful APIs with structured JSON payloads, modular domain packages, and clean layering (Router / Handler → Service / Business Logic → Repository / Storage).
* **Database & Caching:**
  * **Primary Database:** **PostgreSQL** — Relational database providing ACID compliance, foreign-key integrity, and optimized query execution.
  * **Cache & Session Management:** **Redis** — Fast in-memory key-value store for caching active courses, schedule lookup tables, and managing user sessions.
* **DevOps, Infrastructure & Tooling:**
  * **Containerization:** **Docker** and **Docker Compose** for reproducible local and staging environments.
  * **Web Server & Reverse Proxy:** **Nginx** for SSL termination, reverse proxy routing, and load balancing.
  * **CI/CD:** **GitHub Actions** for automated build checks, linting, and unit test execution.
  * **Testing:** Go's built-in `testing` package (with `testify`), Jest/Playwright/Vitest for frontend unit & component testing.

---

## 5. Repository Directory Layout

The repository follows a clean, standardized software engineering structure:

```text
introse-portal67/
├── .github/
│   ├── copilot-commit-instructions.md  # Enforced Conventional Commits standard
│   └── workflows/                      # GitHub Actions CI/CD workflows
├── .vscode/
│   └── settings.json                   # Shared workspace settings & commit prompts
├── docs/                               # All software engineering documentation
│   ├── ai-usage/                       # AI usage tracking logs (mandatory for NMCNPM)
│   ├── design/                         # C4 models, ERDs, class diagrams, wireframes
│   │   ├── pa1-preliminary-architecture.png
│   │   ├── pa1-schedule.png
│   │   └── pa1-wbs.png
│   ├── management/                     # Project plans, agendas, meeting minutes
│   │   └── weekly/                     # Weekly progress status reports (e.g., W1-status.docx)
│   ├── requirements/                   # PRS, SRS, use case models, user stories
│   └── test/                           # Test plans, test cases, automated test reports
├── pa/                                 # Official submission deliverables per phase gate
│   ├── PA0/                            # Project Idea & Team Contract (T02-PA0.pdf)
│   ├── PA1/                            # Project Proposal & Plan (T02-PA1.pdf)
│   ├── PA2/                            # Requirements Specification
│   ├── PA3/                            # System & Architectural Design
│   ├── PA4/                            # Test Planning & Test Strategy
│   └── PA5/                            # Implementation, Testing & Deployment
├── src/                                # Source code root
│   ├── backend/                        # Golang REST API backend
│   └── frontend/                       # Svelte / Tailwind CSS web application
├── .gitignore                          # Ignored artifacts, dependencies, and environment files
├── AGENT.md                            # Comprehensive guidance for AI agents & developers
└── README.md                           # Public repository introduction
```

---

## 6. Project Lifecycle & Phase Gates (Waterfall Model)

The project follows the classical Waterfall process model. Every phase concludes with an official deliverable (**PA**) serving as a project baseline:

```
  PA0: Team Contract & Idea (Week 2) [DONE]
          │
          ▼
  PA1: Proposal & Work Plan (Week 3) [DONE]
          │
          ▼
  PA2: Requirements Specification (Week 5) [CURRENT FOCUS]
          │
          ▼
  PA3: Architecture & System Design (Week 7)
          │
          ▼
  PA4: Test Strategy & Planning (Week 8)
          │
          ▼
  PA5: Construction, Testing & Deployment (Week 10)
          │
          ▼
  Final Presentation & Live Demo (Week 11)
```

### Detailed Phase Gates

| Phase Gate | Milestone | Deadline | Required Exit Deliverables | Status |
| :--- | :--- | :--- | :--- | :--- |
| **PA0** | Initiation & Contract | End of W2 | Project pitch, target user definition, team contract, ground rules, signatures. | **Completed** |
| **PA1** | Proposal & Project Plan | End of W3 | Problem statement, preliminary architecture, WBS, effort estimation, Gantt schedule. | **Completed** |
| **PA2** | Requirements Engineering | End of W5 | PRS & SRS, functional & non-functional requirements, use case specifications, wireframes. | **In Progress** |
| **PA3** | System Design | End of W7 | C4 architecture diagrams, DB schema (ERD), class diagrams, REST API specs, sequence diagrams. | Upcoming |
| **PA4** | Test Planning | End of W8 | Comprehensive test plan, unit/integration/system test cases, test data sets, pass/fail metrics. | Upcoming |
| **PA5** | Construction & Deploy | End of W10 | Full Golang + Svelte implementation of Must-have features, Docker deployment, test verification. | Upcoming |
| **Final** | Final Demo & Presentation | W11 | Live demonstration of all core workflows, presentation slides, final project report. | Upcoming |

---

## 7. Development Guidelines & Team Policies

### 7.1. Git Branching & Pull Request (PR) Policy
1. **Branch Protection:** The `main` branch is protected. **Direct commits or pushes to `main` are strictly prohibited.**
2. **Branch Naming Standard:**
   * Feature: `feat/<feature-name>` (e.g., `feat/course-registration`)
   * Bug Fix: `fix/<bug-description>` (e.g., `fix/gpa-calculation`)
   * Documentation: `docs/<topic>` (e.g., `docs/pa2-srs`)
   * Testing: `test/<module>` (e.g., `test/auth-service`)
3. **Pull Request Protocol:**
   * Every branch must be merged via a GitHub Pull Request.
   * Every PR requires **at least one approving code review** from another team member before merging.
4. **No Force Push:**
   * Never execute `git push --force` or `--force-with-lease` on shared collaborative branches.
   * *(Violations incur an internal fine of 20,000 VND as agreed upon in the Team Contract).*

### 7.2. Conventional Commit Guidelines
All commit messages must strictly adhere to the guidelines set in `.github/copilot-commit-instructions.md`:
* **Format:** `<type>(<scope>): <short summary>`
* **Allowed Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`.
* **Rules:**
  1. Use the imperative, present-tense mood in the summary (e.g., `add feature`, never `added feature` or `adds feature`).
  2. Do not capitalize the first letter of the short summary.
  3. Do not place a period (`.`) at the end of the commit header.
  * **Valid examples:**
    * `feat(registration): add seat limit check before enrollment`
    * `fix(gpa): correct weighted average formula for elective courses`
    * `docs(readme): update setup instructions for golang backend`
  * **Invalid examples:**
    * `Added new feature.` *(wrong tense, capitalized, has trailing dot)*
    * `WIP: fix stuff` *(invalid type)*

### 7.3. Security & Secret Management
* **Zero Secrets in Git:** Never commit secrets, credentials, database passwords, API tokens, or session private keys.
* Keep all `.env`, `.env.local`, and sensitive configuration files in `.gitignore`.
* Provide an `.env.example` template with dummy placeholder values for team setup.

### 7.4. AI Usage Transparency & Logging
* As required by HCMUS faculty policy, all AI-generated or AI-assisted contributions (code generation, refactoring, architectural documentation) must be logged in `docs/ai-usage/`.
* Log entries must specify: tool used, purpose/prompt, scope of code/docs produced, and human review/verification steps.

### 7.5. Communication & Weekly Status Cadence
* Team communications occur primarily via Discord (and Messenger). Team members are expected to respond within 24 hours on weekdays.
* Weekly meetings are held on weekends.
* **Weekly Status Report Deadline:** Every Wednesday at 23:59 (updated from PA0).
* **3-Strike Accountability Rule:**
  * Strike 1: Rescheduled deadline recorded in the weekly status report.
  * Strike 2: Reported to Teaching Assistant (TA); individual contribution factor penalized.
  * Strike 3: Escalated to TA with request to remove member from group.

---

## 8. Specific Guidelines for AI Agents

When acting as an AI assistant within this repository:

### 8.1. Backend Development (Golang)
* **Code Idioms:** Write idiomatic, clean Go code adhering to `gofmt` and `golangci-lint` conventions.
* **Error Handling:** Explicitly handle errors (`if err != nil { return err }`). Avoid unchecked panics.
* **Concurrency:** Use goroutines and channels safely; prevent data races using `sync.Mutex` or atomic operations where appropriate.
* **Package Architecture:** Structure Go code cleanly (e.g., `cmd/server/main.go`, `internal/handlers`, `internal/services`, `internal/repository`, `internal/models`).
* **API Standards:** Return uniform REST JSON envelopes with standard HTTP status codes:
  ```json
  {
    "success": true,
    "data": { ... },
    "error": null
  }
  ```
* **Testing:** Provide table-driven tests for Go business logic using the `testing` package.

### 8.2. Frontend Development (Svelte & Tailwind CSS)
* **Reactivity:** Leverage Svelte's lightweight reactivity model (`$state`, `$derived`, or stores depending on Svelte version). Keep components small, modular, and single-purpose.
* **Styling:** Use Tailwind CSS utility classes; avoid inline styles or custom ad-hoc CSS unless strictly necessary.
* **UX & Accessibility:** Ensure all forms display clear loading indicators, disabled states during API calls, accessible labels, and actionable error messages.
* **Type Safety:** Use TypeScript with Svelte for type-safe props and API response contracts.

### 8.3. Human-in-the-Loop & Verification
* Never delete existing working files, documentation, or git keeps without explicit confirmation.
* Before declaring any implementation complete, verify and test the code.
* Always preserve existing comments and documentation integrity.
