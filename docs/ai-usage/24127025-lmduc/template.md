# AI Usage Log - Tuần [X]

* **Họ và tên:** Lê Minh Đức
* **Mã số sinh viên (MSSV):** 24127025
* **Vai trò chính:** Business Analyst (BA)
* **Vai trò phụ:** UI/UX Designer
* **Giai đoạn dự án (Phase Gate):** [PA1 - Proposal / PA2 - Requirements / PA3 - Design / PA4 - Test Plan / PA5 - Construction]
* **Khoảng thời gian:** [Từ ngày DD/MM/YYYY đến DD/MM/YYYY]
* **Mục tiêu cá nhân trong tuần:**
  - [Mục tiêu 1: ví dụ - Hoàn thiện tài liệu SRS cho nhóm tính năng F01 (Đăng ký học phần) và F06 (Gửi yêu cầu giấy tờ)]
  - [Mục tiêu 2: ví dụ - Viết chi tiết Use Case Specifications và xác định Non-Functional Requirements]
  - [Mục tiêu 3: ví dụ - Phác thảo Wireframe màn hình Đăng ký học phần]

---

## 1. Nhật Ký Chi Tiết Các Phiên Làm Việc Với AI (AI Sessions)

### Phiên 1: [Tác vụ Nghiệp vụ & Yêu cầu - Ví dụ: Xây dựng Use Case Specification & Luồng ngoại lệ]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [Claude 3.5 Sonnet / ChatGPT-4o / Antigravity]
* **Mục đích sử dụng:** [Gợi ý kịch bản tương tác người dùng - hệ thống, bắt trọn các tình huống lỗi nghiệp vụ]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt phân tích nghiệp vụ vào đây. Ví dụ: "Đóng vai trò Business Analyst, hãy viết Use Case chi tiết cho tính năng 'Đăng ký học phần' của hệ thống Portal đại học. Xác định Pre-conditions, Post-conditions, Main Flow và ít nhất 4 Exception Flows (vượt sĩ số, trùng lịch, nợ môn tiên quyết, hết hạn đăng ký)..."]
  ```
* **Tóm tắt kết quả AI sinh ra:**
  - [AI sinh bảng Use Case Specification chuẩn cấu trúc IEEE 830]
  - [Liệt kê các bước thao tác giữa Sinh viên và Hệ thống cùng các thông báo lỗi tương ứng]
* **Đánh giá & Chỉnh sửa của cá nhân (Human Review & Modifications):**
  - **Điểm AI làm tốt:** [Mô tả luồng chính rõ ràng, cách đặt tên Actor và Trigger chuẩn xác]
  - **Điểm thiếu sót / Cần tinh chỉnh:** [AI đề xuất luồng cho phép sinh viên đăng ký 'danh sách chờ' (waitlist) - tính năng này ngoài phạm vi MoSCoW MVP của Portal67; chưa xét đến quy chế quy đổi học phần của HCMUS]
  - **Chỉnh sửa thực tế đã làm:** [Lược bỏ tính năng waitlist để bám sát MVP; chuẩn hóa bảng mã lỗi theo format JSON quy định trong AGENT.md; bổ sung bảng phân quyền RBAC]
* **Kiểm chứng (Verification Steps):**
  - [x] Rà soát lại với tài liệu PRD / SRS tổng thể.
  - [x] Thảo luận với System Architect (Vũ Trần Minh Hiếu) để đảm bảo tính khả thi về mặt kiến trúc dữ liệu.

---

### Phiên 2: [Tác vụ Thiết kế Giao diện (UI/UX) - Ví dụ: Ý tưởng bố cục Wireframe & Trải nghiệm người dùng]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [Claude / ChatGPT / UI-UX tools]
* **Mục đích sử dụng:** [Gợi ý phân cấp thông tin (Information Architecture), luồng thao tác giảm thiểu click cho sinh viên khi tra cứu thời khóa biểu]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt UI/UX vào đây]
  ```
* **Tóm tắt kết quả AI sinh ra:**
  - [...]
* **Đánh giá & Chỉnh sửa của cá nhân:**
  - **Điểm AI làm tốt:** [...]
  - **Điểm hạn chế / Lỗ hổng:** [...]
  - **Chỉnh sửa thực tế:** [...]
* **Kiểm chứng:**
  - [ ] Đã chuyển đổi ý tưởng thành wireframe trên Figma / tài liệu `docs/design/`.

---

## 2. Liên Kết Minh Chứng Đối Chứng (Traceability & Artifacts)

| Loại tài liệu | Đường dẫn file trong Repository | Commit Hash / PR liên quan |
| :--- | :--- | :--- |
| **Yêu cầu (SRS/Use Case)** | [`docs/requirements/...`](#) | `docs(req): add use case spec` |
| **Thiết kế Wireframe/UI** | [`docs/design/...`](#) | Wireframe diagrams |

---

## 3. Tự Đánh Giá & Bài Học Kinh Nghiệm (Self-Reflection)

* **Ước tính mức độ hỗ trợ của AI trong tuần:** [Ví dụ: 35% AI - 65% Phân tích và biên tập thực tế]
* **Thời gian tiết kiệm được:** [Ví dụ: ~3.5 giờ]
* **Bài học rút ra:** [Ví dụ: AI có xu hướng 'thổi phồng' tính năng (over-feature). Khi prompt cần đóng khung nghiêm ngặt danh mục MoSCoW trong AGENT.md để không bị loãng phạm vi đồ án]
