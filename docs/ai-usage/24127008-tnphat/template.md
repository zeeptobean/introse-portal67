# AI Usage Log - Tuần [X]

* **Họ và tên:** Trương Nhật Phát
* **Mã số sinh viên (MSSV):** 24127008
* **Vai trò chính:** QA / Tester
* **Vai trò phụ:** PM / Business Analyst
* **Giai đoạn dự án (Phase Gate):** [PA2 - Requirements / PA4 - Test Strategy & Planning / PA5 - Testing & Verification]
* **Khoảng thời gian:** [Từ ngày DD/MM/YYYY đến DD/MM/YYYY]
* **Mục tiêu cá nhân trong tuần:**
  - [Mục tiêu 1: ví dụ - Xây dựng Master Test Plan cho giai đoạn PA4]
  - [Mục tiêu 2: ví dụ - Thiết kế ma trận kiểm thử (Test Case Matrix) bằng kỹ thuật Phân vùng tương đương (EP) & Phân tích giá trị biên (BVA) cho tính năng F01 và F03]
  - [Mục tiêu 3: ví dụ - Viết kịch bản kiểm thử API tự động và rà soát bug hệ thống]

---

## 1. Nhật Ký Chi Tiết Các Phiên Làm Việc Với AI (AI Sessions)

### Phiên 1: [Tác vụ Thiết kế Test Case / Ma trận Kiểm thử - Ví dụ: Test cases cho Tính năng Đăng ký môn học]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [ChatGPT-4o / Claude 3.5 Sonnet / Antigravity]
* **Mục đích sử dụng:** [Gợi ý các ca kiểm thử biên (Boundary Cases), ca kiểm thử tiêu cực (Negative Cases) và dữ liệu test giả lập (Mock Test Data)]
* **Prompt đưa vào (Input Prompt):**
  ```text
  [Dán prompt QA/Testing vào đây. Ví dụ: "Đóng vai trò Senior QA Engineer, hãy thiết kế bảng Test Cases chi tiết cho tính năng 'Đăng ký học phần' của Portal67. Áp dụng kỹ thuật Boundary Value Analysis và Equivalence Partitioning cho các trường: sĩ số lớp (1..50), số tín chỉ tối đa học kỳ (14..24 tín chỉ), trùng lịch tiết 1-3 thứ Hai. Cung cấp Test Case ID, Mục đích, Pre-conditions, Input Data, Expected Results..."]
  ```
* **Tóm tắt kết quả AI sinh ra:**
  - [AI sinh bảng gồm 12 Test Cases bao gồm cả Positive và Negative cases]
  - [Phân loại đầy đủ mức độ nghiêm trọng (Severity/Priority)]
* **Đánh giá & Chỉnh sửa của cá nhân (Human Review & Modifications):**
  - **Điểm AI làm tốt:** [Bao phủ tốt các trường hợp biên của số tín chỉ (13, 14, 24, 25 tín chỉ); định dạng bảng test case chuẩn quốc tế IEEE 829]
  - **Điểm thiếu sót / Cần bổ sung:** [AI chưa bao quát được trường hợp sinh viên đăng ký môn học có nhiều ca thực hành khác nhau cùng thuộc 1 mã môn lý thuyết; thiếu kịch bản test timeout khi server quá tải]
  - **Chỉnh sửa thực tế đã làm:** [Bổ sung thêm 4 test cases cho trường hợp lớp Lý thuyết + Thực hành phụ thuộc; thêm kịch bản kiểm thử tải đồng thời (concurrent registration); chuẩn hóa mã test case theo tiền tố `TC_F01_xxx`]
* **Kiểm chứng (Verification Steps):**
  - [x] Rà soát đối chiếu lại với tài liệu SRS của Business Analyst (Lê Minh Đức).
  - [x] Lưu bảng Test Cases vào thư mục `docs/test/`.

---

### Phiên 2: [Tác vụ Tự động hóa Kiểm thử / Báo cáo Bug - Ví dụ: Kịch bản API Testing hoặc E2E]
* **Ngày thực hiện:** [DD/MM/YYYY]
* **Công cụ sử dụng:** [GitHub Copilot / Claude]
* **Mục đích sử dụng:** [Gợi ý test script kiểm thử API (Postman/Go test) hoặc format bug report]
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
  - [ ] Đã chạy thử bộ test script và kiểm tra pass/fail rate.

---

## 2. Liên Kết Minh Chứng Đối Chứng (Traceability & Artifacts)

| Loại tài liệu / Kịch bản test | Đường dẫn file trong Repository | Commit Hash / PR liên quan |
| :--- | :--- | :--- |
| **Kế hoạch kiểm thử (Test Plan)** | [`docs/test/...`](#) | `docs(test): add test plan PA4` |
| **Ma trận Test Cases** | [`docs/test/test-cases-f01.md`](#) | Test Matrix Docs |
| **Báo cáo Bug / Issue** | GitHub Issues #... | Bug report link |

---

## 3. Tự Đánh Giá & Bài Học Kinh Nghiệm (Self-Reflection)

* **Ước tính mức độ hỗ trợ của AI trong tuần:** [Ví dụ: 35% AI - 65% Phân tích nghiệp vụ kiểm thử và xác minh kết quả]
* **Thời gian tiết kiệm được:** [Ví dụ: ~3 giờ sinh dữ liệu mẫu và định dạng bảng biểu]
* **Bài học rút ra:** [Ví dụ: AI tạo test case tổng quát rất nhanh nhưng thường bỏ sót các nghiệp vụ đặc thù của trường như mã môn tương đương hoặc liên kết lớp lý thuyết - bài tập. QA cần dựa vào kinh nghiệm thực tế để bổ sung các ca kiểm thử hóc búa]
