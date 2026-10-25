# AI Usage Log - Tuần [X]

* **Họ và tên:** Đoàn Quốc Bảo
* **Mã số sinh viên (MSSV):** 24127326
* **Vai trò chính:** Team Lead / Project Manager (PM)
* **Vai trò phụ:** Developer
* **Giai đoạn dự án (Phase Gate):** [PA1 - Proposal / PA2 - Requirements / PA3 - Design / PA4 - Test Plan / PA5 - Construction]
* **Khoảng thời gian:** [Từ ngày DD/MM/YYYY đến DD/MM/YYYY]
* **Mục tiêu cá nhân trong tuần:**
  - [Mục tiêu 1: ví dụ - Lên kế hoạch phân rã công việc WBS & tiến độ Gantt cho PA2]
  - [Mục tiêu 2: ví dụ - Điều phối họp nhóm tuần, rà soát tiến độ nộp báo cáo status report]
  - [Mục tiêu 3: ví dụ - Hỗ trợ review pull request và quản lý GitHub Project Board]

---

## 1. Nhật Ký Chi Tiết Các Phiên Làm Việc Với AI (AI Sessions)

### Phiên 1: [Tác vụ Quản lý Dự án / Lập Kế hoạch - Ví dụ: Xây dựng Work Breakdown Structure (WBS) & Dự báo rủi ro]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [ChatGPT-4o / Claude 3.5 Sonnet / Antigravity / Gemini]
* **Mục đích sử dụng:** [Gợi ý phân rã nhiệm vụ cho milestone hiện tại, phân bổ effort giờ làm, nhận diện rủi ro chậm tiến độ]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt quản trị dự án vào đây. Ví dụ: "Đóng vai trò Project Manager theo mô hình Waterfall cho đồ án đại học. Hãy phân rã WBS cấp độ 3 cho giai đoạn Requirements Specification (PA2) gồm 5 thành viên với các vai trò PM, BA, Architect, Lead Dev, Tester..."]
  ```
* **Tóm tắt kết quả AI sinh ra:**
  - [AI liệt kê danh mục 15 công việc con kèm ước lượng thời gian hoàn thành theo giờ]
  - [Đề xuất ma trận RACI và 3 rủi ro chính về xung đột lịch học của sinh viên]
* **Đánh giá & Chỉnh sửa của cá nhân (Human Review & Modifications):**
  - **Điểm AI làm tốt:** [Cấu trúc WBS rõ ràng, đúng chuẩn quản trị PMBOK]
  - **Điểm chưa phù hợp với nhóm:** [AI ước tính thời gian quá dư dả so với cam kết 6h/tuần/thành viên của Team 02; các task chưa gắn đúng với 8 tính năng F01-F08 trong AGENT.md]
  - **Chỉnh sửa thực tế đã làm:** [Điều chỉnh lại thời lượng thực tế phù hợp với quỹ thời gian tuần; map trực tiếp task vào issue trên GitHub Projects; bổ sung điều khoản 3-strike rule của nhóm]
* **Kiểm chứng (Verification Steps):**
  - [x] Đã thống nhất phân bổ công việc với các thành viên trong cuộc họp Discord.
  - [x] Đã cập nhật lên GitHub Projects và file báo cáo tuần `docs/management/weekly/`.

---

### Phiên 2: [Tác vụ Kỹ thuật / Code / Review - Ví dụ: Rà soát Pull Request hoặc Hỗ trợ cấu hình dự án]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [GitHub Copilot / Claude / Antigravity]
* **Mục đích sử dụng:** [Hỗ trợ review code backend Go / frontend Svelte hoặc cấu hình CI/CD GitHub Actions]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt vào đây]
  ```
* **Tóm tắt kết quả AI sinh ra:**
  - [...]
* **Đánh giá & Chỉnh sửa của cá nhân:**
  - **Điểm AI làm tốt:** [...]
  - **Điểm hạn chế / Lỗ hổng:** [...]
  - **Chỉnh sửa thực tế:** [...]
* **Kiểm chứng:**
  - [ ] Đã kiểm tra commit tuân thủ Conventional Commits theo quy định nhóm.
  - [ ] Đã test chạy thử workflow CI/CD.

---

## 2. Liên Kết Minh Chứng Đối Chứng (Traceability & Artifacts)

| Loại tài liệu / Công việc | Đường dẫn file trong Repository | Link Issue / PR / Commit |
| :--- | :--- | :--- |
| **Kế hoạch & Báo cáo** | [`docs/management/weekly/...`](#) | Weekly Status Report |
| **Bảng công việc** | GitHub Projects Board | Task Issue #... |
| **Review / Đóng góp mã** | [`src/...`](#) | PR #... |

---

## 3. Tự Đánh Giá & Bài Học Kinh Nghiệm (Self-Reflection)

* **Ước tính mức độ hỗ trợ của AI trong tuần:** [Ví dụ: 25% AI - 75% Tự thực hiện và điều phối]
* **Thời gian tiết kiệm được:** [Ví dụ: ~2.5 giờ]
* **Bài học rút ra:** [Ví dụ: Khi dùng AI để phân bổ task, cần nạp kèm danh sách vai trò và cam kết giờ của từng thành viên để AI không chia việc lệch năng lực]
