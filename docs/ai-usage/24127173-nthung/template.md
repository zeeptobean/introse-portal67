# AI Usage Log - Tuần [X]

* **Họ và tên:** Nguyễn Tuấn Hùng
* **Mã số sinh viên (MSSV):** 24127173
* **Vai trò chính:** Lead Developer
* **Vai trò phụ:** QA / Tester
* **Giai đoạn dự án (Phase Gate):** [PA3 - System Design / PA4 - Test Plan / PA5 - Construction & Testing]
* **Khoảng thời gian:** [Từ ngày DD/MM/YYYY đến DD/MM/YYYY]
* **Mục tiêu cá nhân trong tuần:**
  - [Mục tiêu 1: ví dụ - Xây dựng tầng Repository và Service cho API Đăng ký môn học (Golang)]
  - [Mục tiêu 2: ví dụ - Viết Unit Tests với `testify` và xử lý chống Race Condition khi sinh viên đăng ký đồng thời]
  - [Mục tiêu 3: ví dụ - Chuẩn hóa envelope response JSON theo chuẩn dự án]

---

## 1. Nhật Ký Chi Tiết Các Phiên Làm Việc Với AI (AI Sessions)

### Phiên 1: [Tác vụ Backend Golang - Ví dụ: Xây dựng Service Handler & Xử lý Transaction CSDL]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [GitHub Copilot / Cursor / Claude 3.5 Sonnet / Antigravity]
* **Mục đích sử dụng:** [Sinh cấu trúc nghiệp vụ sạch (Clean Layered Architecture), quản lý lỗi idiomatic Go và transaction rollback]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt coding Go vào đây. Ví dụ: "Viết hàm EnrollCourseService trong Go nhận vào studentID và sectionID. Logic cần thực hiện: 1. Kiểm tra section còn slot không; 2. Kiểm tra sinh viên có bị trùng giờ học không; 3. Thêm record vào bảng enrollments và tăng current_enrolled lên 1. Tất cả phải nằm trong một DB Transaction với mức khóa thích hợp..."]
  ```
* **Tóm tắt kết quả AI sinh ra:**
  - [AI sinh function Go với context `context.Context`, sử dụng `sqlx` hoặc GORM transaction `tx.Begin()`]
  - [Tạo custom error struct để map ra HTTP status codes tương ứng]
* **Đánh giá & Chỉnh sửa của cá nhân (Human Review & Modifications):**
  - **Điểm AI làm tốt:** [Cú pháp Go chuẩn idiomatic, xử lý kiểm tra `err != nil` ở mọi bước, code dễ đọc]
  - **Điểm hạn chế / Bug tiềm ẩn:** [AI chỉ dùng `SELECT` thông thường trước khi `UPDATE`, gây ra race condition nếu có 2 request đồng thời; thiếu lock hàng `SELECT ... FOR UPDATE` trong transaction]
  - **Chỉnh sửa thực tế đã làm:** [Thêm câu lệnh SQL `FOR UPDATE` khi đọc dữ liệu section; chuẩn hóa mã lỗi trả về theo chuẩn JSON quy định tại mục 8.1 trong `AGENT.md`]
* **Kiểm chứng (Verification Steps):**
  - [x] Chạy kiểm tra tĩnh với `golangci-lint run`.
  - [x] Viết benchmark / concurrent test với goroutine để xác nhận không bị vượt sĩ số.

---

### Phiên 2: [Tác vụ Kiểm thử Đơn vị (Unit Testing) - Ví dụ: Viết Table-Driven Tests với Mock DB]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [GitHub Copilot / ChatGPT]
* **Mục đích sử dụng:** [Sinh ma trận test case kiểm tra table-driven testing cho hàm tính GPA hoặc kiểm tra trùng lịch]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt viết unit test vào đây]
  ```
* **Tóm tắt kết quả AI sinh ra:**
  - [...]
* **Đánh giá & Chỉnh sửa của cá nhân:**
  - **Điểm AI làm tốt:** [...]
  - **Điểm hạn chế / Lỗ hổng:** [...]
  - **Chỉnh sửa thực tế:** [...]
* **Kiểm chứng:**
  - [ ] Chạy `go test -v -cover ./...` đạt độ bao phủ code (coverage) mong đợi.

---

## 2. Liên Kết Minh Chứng Đối Chứng (Traceability & Artifacts)

| Loại mã nguồn / Test | Đường dẫn file trong Repository | Commit Hash / PR liên quan |
| :--- | :--- | :--- |
| **Backend Go API** | [`src/backend/internal/...`](#) | `feat(api): implement course enrollment service` |
| **Unit Test Go** | [`src/backend/internal/..._test.go`](#) | `test(enrollment): add concurrent tests` |

---

## 3. Tự Đánh Giá & Bài Học Kinh Nghiệm (Self-Reflection)

* **Ước tính mức độ hỗ trợ của AI trong tuần:** [Ví dụ: 40% AI - 60% Con người trực tiếp viết & debug]
* **Thời gian tiết kiệm được:** [Ví dụ: ~4 giờ viết boilerplate và mock struct]
* **Bài học rút ra:** [Ví dụ: AI rất giỏi sinh code nhanh, nhưng với các bài toán chịu tải cao và concurrency trong Go, bắt buộc lập trình viên phải nắm sâu về cơ chế khóa dữ liệu (database lock/mutex) để rà soát lại]
