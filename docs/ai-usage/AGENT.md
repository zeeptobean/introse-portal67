# AGENT.md - AI Agent Operational Guidelines for AI Usage Logging

This document defines explicit instructions and protocols for all **AI Coding Assistants, LLMs, and Autonomous Agents** (including Antigravity, Claude, ChatGPT, GitHub Copilot, Cursor, Gemini) interacting with the **Portal67** repository.

---

## 1. Prime Directive: Academic Integrity & Transparency

As mandated by HCMUS Software Engineering Department policy:
* **Never obscure or conceal AI-assisted generation.**
* Whenever you (the AI Agent) generate significant architectural designs, business requirements, backend/frontend code, or test specifications, **you must assist the student in logging the interaction transparently.**

---

## 2. Team Member Mapping Reference

Use this table to map the active student (via Git config, author email, or user prompt) to their dedicated log folder:

| Student ID | Full Name | Academic Email | Role | Log Directory |
| :--- | :--- | :--- | :--- | :--- |
| **24127326** | **Đoàn Quốc Bảo** | `dqbao2432@clc.fitus.edu.vn` | PM / Lead / Dev | `docs/ai-usage/24127326-dqbao/` |
| **24127025** | **Lê Minh Đức** | `lmduc2433@clc.fitus.edu.vn` | BA / UI-UX | `docs/ai-usage/24127025-lmduc/` |
| **24127003** | **Vũ Trần Minh Hiếu** | `vtmhieu2416@clc.fitus.edu.vn` | Architect / UI / Dev | `docs/ai-usage/24127003-vtmhieu/` |
| **24127173** | **Nguyễn Tuấn Hùng** | `nthung2432@clc.fitus.edu.vn` | Lead Dev / QA | `docs/ai-usage/24127173-nthung/` |
| **24127008** | **Trương Nhật Phát** | `tnphat2418@clc.fitus.edu.vn` | QA / Tester / PM | `docs/ai-usage/24127008-tnphat/` |

---

## 3. Agent Responsibilities & Workflows

### 3.1. Automated Log Snippet Generation
Whenever you complete a significant task (e.g., implementing an API handler in Go, writing a Svelte component, drafting a Use Case specification in Markdown, or writing unit tests), provide a pre-formatted **AI Usage Log Entry** at the end of your response or offer to append it to the current week's log.

Format of the snippet to suggest to the user:
```markdown
### Phiên [X]: [Tên tác vụ vừa thực hiện]
* **Ngày thực hiện:** YYYY-MM-DD
* **Công cụ:** [Tên AI Agent / Model sử dụng]
* **Mục đích:** [Tóm tắt mục đích công việc]
* **Prompt chính:** "[Tóm tắt câu lệnh của người dùng hoặc prompt quan trọng]"
* **AI Output tóm tắt:** [Mô tả ngắn các file, struct, hoặc logic đã sinh ra]
* **Đánh giá & Chỉnh sửa của người (Human Review & Modifications):**
  - **Điểm AI làm tốt:** [Cấu trúc code/tài liệu được sinh nhanh, đúng cú pháp]
  - **Điểm hạn chế cần chú ý:** [Các điểm biên, giả định hoặc hạn chế cần con người rà soát]
  - **Chỉnh sửa thực tế:** [Những gì sinh viên đã can thiệp, chỉnh sửa để khớp với Portal67]
* **Kiểm chứng (Verification):** [Chạy test `go test`, linter `golangci-lint`, hoặc peer-review]
```

### 3.2. Proactive "Human-in-the-Loop" Prompts
To help students score high on evaluation rubrics:
* **Flag Assumptions:** When generating code or documentation, explicitly notify the student of assumptions made (e.g., *"Note: I assumed max course enrollment is 50. Please verify this constraint against the project requirements."*).
* **Highlight Edge Cases:** Mention boundary conditions that require human verification (e.g., concurrency conflicts, time overlap, authorization boundaries).

### 3.3. Safeguards & Prohibitions
* **Do NOT blindly overwrite or erase** any student's existing weekly log files (`W01.md`, `W02.md`, etc.).
* **Preserve individual reflections:** Never fabricate personal thoughts or reflections on behalf of the student; always leave those sections as editable templates or ask for their real feedback.
* **Respect Project Scope:** Always keep generated code and specs aligned with the **Mock Data** and **No Real Payment** constraints specified in the main `AGENT.md`.
