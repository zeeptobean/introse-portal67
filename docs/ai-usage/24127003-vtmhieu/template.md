# AI Usage Log - Tuần [X]

* **Họ và tên:** Vũ Trần Minh Hiếu
* **Mã số sinh viên (MSSV):** 24127003
* **Vai trò chính:** System Architect / UI Designer
* **Vai trò phụ:** Developer
* **Giai đoạn dự án (Phase Gate):** [PA1 - Proposal / PA2 - Requirements / PA3 - System Design / PA4 - Test Plan / PA5 - Construction]
* **Khoảng thời gian:** [Từ ngày DD/MM/YYYY đến DD/MM/YYYY]
* **Mục tiêu cá nhân trong tuần:**
  - [Mục tiêu 1: ví dụ - Thiết kế sơ đồ kiến trúc C4 Container Diagram & Component Diagram cho phân hệ Đăng ký môn]
  - [Mục tiêu 2: ví dụ - Xây dựng lược đồ cơ sở dữ liệu quan hệ ERD PostgreSQL (các bảng courses, sections, enrollments, users)]
  - [Mục tiêu 3: ví dụ - Đặc tả chuẩn REST API response envelope và Sequence Diagram cho xác thực SSO/RBAC]

---

## 1. Nhật Ký Chi Tiết Các Phiên Làm Việc Với AI (AI Sessions)

### Phiên 1: [Tác vụ Kiến trúc Hệ thống / Thiết kế CSDL - Ví dụ: Thiết kế ERD & Tối ưu khóa/Index cho PostgreSQL]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [Claude 3.5 Sonnet / ChatGPT-4o / Antigravity]
* **Mục đích sử dụng:** [Tạo mã Mermaid cho ERD, gợi ý ràng buộc khóa ngoại (FK), chuẩn hóa 3NF và chiến lược đánh index cho bảng có lượng truy vấn cao]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt kiến trúc/database vào đây. Ví dụ: "Đóng vai trò Software Architect cho hệ thống Portal67 (Golang + PostgreSQL). Hãy thiết kế lược đồ ERD định dạng Mermaid cho module đăng ký môn học gồm: students, courses, course_sections, semesters, enrollments. Đảm bảo toàn vẹn dữ liệu, chống đăng ký trùng lặp và tối ưu truy vấn thời khóa biểu..."]
  ```
* **Tóm tắt kết quả AI sinh ra:**
  - [AI sinh mã Mermaid biểu diễn 5 thực thể và quan hệ 1-N, N-N]
  - [Gợi ý các trường: UUID làm primary key, enum cho status, và bảng trung gian `enrollments` với composite unique key]
* **Đánh giá & Chỉnh sửa của cá nhân (Human Review & Modifications):**
  - **Điểm AI làm tốt:** [Cú pháp Mermaid chuẩn xác, quan hệ giữa các bảng chặt chẽ, đề xuất enum status hợp lý]
  - **Điểm hạn chế / Rủi ro thiết kế:** [AI dùng kiểu `TEXT` tự do cho ca học và phòng học; chưa có trường đánh dấu số lượng sinh viên hiện tại `current_enrolled` để tối ưu query đếm chỗ; chưa có cơ chế Soft Delete]
  - **Chỉnh sửa thực tế đã làm:** [Bổ sung cột `max_capacity` và `enrolled_count` có CHECK constraint `enrolled_count <= max_capacity`; chuẩn hóa định dạng thời gian tiết học theo hệ thống tín chỉ HCMUS; bổ sung trường `deleted_at` cho soft-delete]
* **Kiểm chứng (Verification Steps):**
  - [x] Render thử sơ đồ Mermaid trên GitHub Markdown preview và Mermaid Live Editor.
  - [x] Thảo luận với Lead Developer (Nguyễn Tuấn Hùng) để chốt schema trước khi viết migration.

---

### Phiên 2: [Tác vụ Thiết kế Giao diện / Frontend Architecture - Ví dụ: Cấu trúc Component SvelteKit & Design Tokens]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [GitHub Copilot / Claude]
* **Mục đích sử dụng:** [Gợi ý phân rã layout SvelteKit, reactive store/state quản lý giỏ môn học, cấu hình Tailwind CSS]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt vào đây]
  ```
* **Tóm tắt kết quả AI sinh ra:**
  - [...]
* **Đánh giá & Chỉnh sửa của cá nhân:**
  - **Điểm AI làm tốt:** [...]
  - **Điểm hạn chế / Cần sửa:** [...]
  - **Chỉnh sửa thực tế:** [...]
* **Kiểm chứng:**
  - [ ] Khảo sát khả năng responsive trên các viewport theo cam kết đồ án.

---

## 2. Liên Kết Minh Chứng Đối Chứng (Traceability & Artifacts)

| Loại tài liệu / Thiết kế | Đường dẫn file trong Repository | Commit Hash / PR liên quan |
| :--- | :--- | :--- |
| **Sơ đồ kiến trúc (C4/ERD)** | [`docs/design/...`](#) | `docs(design): add C4 container diagram` |
| **Đặc tả API / Schemas** | [`docs/design/api-specs.md` hoặc code](#) | PR #... |

---

## 3. Tự Đánh Giá & Bài Học Kinh Nghiệm (Self-Reflection)

* **Ước tính mức độ hỗ trợ của AI trong tuần:** [Ví dụ: 30% AI - 70% Thiết kế và thẩm định kiến trúc]
* **Thời gian tiết kiệm được:** [Ví dụ: ~3 giờ]
* **Bài học rút ra:** [Ví dụ: Khi yêu cầu AI vẽ Mermaid diagram, cần chỉ định rõ loại diagram (`erDiagram` hoặc `flowchart LR`) và yêu cầu không dùng các ký tự đặc biệt trong nhãn để tránh lỗi cú pháp parse]
