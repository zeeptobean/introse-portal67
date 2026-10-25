# AI Usage Log - Master Template (Mẫu Chuẩn Chung)

> **Hướng dẫn:** Copy file này vào thư mục cá nhân của bạn (ví dụ: `docs/ai-usage/24127003-vtmhieu/W03.md`) vào đầu mỗi tuần, cập nhật nội dung trong suốt tuần làm việc và hoàn thành trước hạn chót báo cáo tuần (**Thứ Tư hàng tuần lúc 23:59**).

---

# AI Usage Log - Tuần [X]

* **Họ và tên:** [Họ và tên sinh viên]
* **Mã số sinh viên (MSSV):** [24127xxx]
* **Vai trò trong tuần:** [Ví dụ: System Architect / Developer / BA / Tester / PM]
* **Giai đoạn dự án (Phase Gate):** [Ví dụ: PA1 - Proposal / PA2 - Requirements / PA3 - Design / PA5 - Construction]
* **Khoảng thời gian:** [Từ ngày DD/MM/YYYY đến DD/MM/YYYY]
* **Mục tiêu cá nhân trong tuần:** 
  - [Mục tiêu 1: ví dụ - Hoàn thiện Use Case Specification cho F01]
  - [Mục tiêu 2: ví dụ - Thiết kế ERD Database PostgreSQL cho phân hệ Học phần]

---

## 1. Nhật Ký Chi Tiết Các Phiên Làm Việc Với AI (AI Sessions)

### Phiên 1: [Tiêu đề tác vụ - Ví dụ: Soạn thảo Use Case Specification F01]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [Claude 3.5 Sonnet / ChatGPT-4o / GitHub Copilot / Cursor / Antigravity / Gemini]
* **Mục đích sử dụng:** [Tạo khung Use Case, gợi ý luồng ngoại lệ, giải thích thuật ngữ, sinh code, refactor...]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt chính của bạn vào đây. Nên giữ lại các ngữ cảnh, ràng buộc quan trọng bạn đã đưa cho AI]
  Ví dụ: "Đóng vai trò Business Analyst cho hệ thống Portal đại học Portal67. Viết đặc tả Use Case chi tiết cho 'Đăng ký học phần', bao gồm Main Flow và Exception Flows khi hết slot hoặc trùng lịch học..."
  ```
* **Tóm tắt kết quả AI sinh ra (AI Output Summary):**
  - [AI đưa ra bản nháp bảng Use Case gồm 6 bước Actor - System]
  - [Gợi ý 4 luồng ngoại lệ: trùng lịch, vượt quá sĩ số, chưa đạt môn tiên quyết, lỗi mạng]
* **Đánh giá & Chỉnh sửa của người (Human Review & Modifications - BẮT BUỘC):**
  - **Điểm AI làm tốt:** [Cấu trúc bảng chuẩn UML/RUP, nhận diện đủ các trường hợp biên cơ bản]
  - **Điểm AI làm sai / Chưa phù hợp với Portal67:** [AI tự thêm bước thanh toán học phí qua cổng VNPay ngay khi đăng ký môn -> Sai với đặc tả Portal67 vốn chỉ lưu trạng thái công nợ và dùng Mock Data]
  - **Chỉnh sửa thực tế đã làm:** [Xóa bỏ bước thanh toán trực tiếp; điều chỉnh mã lỗi theo quy ước khoa CNTT HCMUS; bổ sung điều kiện sinh viên đang trong học kỳ hợp lệ]
* **Kiểm chứng & Đảm bảo chất lượng (Verification Steps):**
  - [x] Đã đối chiếu với mục tiêu giai đoạn trong `AGENT.md`.
  - [x] Đã kiểm tra chéo (Peer review) với thành viên: [Tên thành viên phản biện].
  - [ ] Đã chạy test tự động / biên dịch thử (nếu là code).

---

### Phiên 2: [Tiêu đề tác vụ - Ví dụ: Sinh code Golang REST Handler]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [GitHub Copilot / Cursor / Claude]
* **Mục đích sử dụng:** [Sinh boilerplate handler đăng ký học phần và xử lý transaction database]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt hoặc mô tả bối cảnh inline completion]
  ```
* **Tóm tắt kết quả AI sinh ra:**
  - [Code handler nhận JSON payload, gọi repository kiểm tra seat limit và insert bảng enrollments]
* **Đánh giá & Chỉnh sửa của người (Human Review & Modifications):**
  - **Điểm AI làm tốt:** [Cú pháp Go chuẩn, error handling rõ ràng]
  - **Điểm hạn chế / Lỗ hổng:** [AI không sử dụng database transaction `tx.Begin()` dẫn đến nguy cơ race condition nếu 2 sinh viên bấm đăng ký cùng lúc cho slot cuối cùng]
  - **Chỉnh sửa thực tế:** [Viết lại logic bọc trong transaction với mức cô lập Serializable/SELECT FOR UPDATE, thêm log có cấu trúc]
* **Kiểm chứng (Verification Steps):**
  - [x] Viết unit test kiểm tra concurrency với 10 goroutines.
  - [x] Chạy lệnh `go test ./...` đạt 100% pass.

*(Nếu có thêm các phiên làm việc khác trong tuần, hãy nhân bản mục `### Phiên X` ở trên)*

---

## 2. Liên Kết Minh Chứng Đối Chứng (Traceability & Artifacts)

| Loại tài liệu / Code | Đường dẫn file trong Repository | Commit Hash / PR liên quan |
| :--- | :--- | :--- |
| **Yêu cầu / Phân tích** | [`docs/requirements/...`](#) | `feat(requirements): add use case spec` |
| **Thiết kế / Kiến trúc** | [`docs/design/...`](#) | `docs(design): update c4 diagram` |
| **Mã nguồn (Source Code)** | [`src/...`](#) | `feat(...): implement service handler` |
| **Kiểm thử (Test Cases)** | [`docs/test/...` hoặc `src/..._test.go`](#) | `test(...): add unit tests` |

---

## 3. Tự Đánh Giá & Bài Học Kinh Nghiệm (Self-Reflection)

* **Ước tính mức độ hỗ trợ của AI trong khối lượng công việc tuần:** [Ví dụ: 30% AI - 70% Con người]
* **Thời gian tiết kiệm được:** [Ví dụ: ~3 giờ so với việc tự viết boilerplate từ đầu]
* **Nhận định về độ tin cậy của AI tuần này:** [Ví dụ: Khá tốt về mặt sinh cấu trúc và code mẫu, nhưng logic nghiệp vụ đặc thù của HCMUS cần con người chỉnh sửa kỹ lưỡng]
* **Kinh nghiệm / Cải tiến cách prompt cho tuần sau:** [Ví dụ: Cần cung cấp trước context từ file `AGENT.md` và database schema để AI không tự ý bịa các trường dữ liệu ngoài phạm vi đồ án]
