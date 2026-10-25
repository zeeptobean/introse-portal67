# Hệ Thống Nhật Ký Sử Dụng AI (AI Usage Logs) - Portal67

Thư mục này lưu trữ toàn bộ nhật ký sử dụng công cụ AI (Generative AI & LLMs) của từng thành viên trong **Nhóm 02 - BDD (Portal67)** thuộc môn **Nhập môn Công nghệ Phần mềm (NMCNPM)** - Khoa CNTT, Trường Đại học Khoa học Tự nhiên, ĐHQG-HCM (HCMUS).

---

## 1. Cấu Trúc Thư Mục

Để dễ dàng quản lý cá nhân và thuận tiện cho Giảng viên / Trợ giảng (TA) tra cứu khi chấm điểm tiến độ định kỳ:

```text
docs/ai-usage/
├── README.md                      # Hướng dẫn chung về quy định và workflow AI Log
├── AGENT.md                       # Bản chỉ dẫn dành riêng cho AI Agents hỗ trợ dự án
├── TEMPLATE.md                    # Template chuẩn tổng quát cho mọi vai trò
├── 24127326-dqbao/               # Thư mục của Đoàn Quốc Bảo (Team Lead / PM / Dev)
│   ├── template.md                # Template cá nhân hóa theo vai trò PM/Dev
│   └── W02.md, W03.md...          # Nhật ký theo từng tuần
├── 24127025-lmduc/               # Thư mục của Lê Minh Đức (BA / UI-UX Designer)
│   ├── template.md                # Template cá nhân hóa theo vai trò BA/UI-UX
│   └── W02.md, W03.md...
├── 24127003-vtmhieu/              # Thư mục của Vũ Trần Minh Hiếu (System Architect / UI / Dev)
│   ├── template.md                # Template cá nhân hóa theo vai trò Architect/Dev
│   └── W02.md, W03.md...
├── 24127173-nthung/               # Thư mục của Nguyễn Tuấn Hùng (Lead Developer / Tester)
│   ├── template.md                # Template cá nhân hóa theo vai trò Lead Dev/QA
│   └── W02.md, W03.md...
└── 24127008-tnphat/               # Thư mục của Trương Nhật Phát (QA / Tester / BA)
    ├── template.md                # Template cá nhân hóa theo vai trò QA/Tester
    └── W02.md, W03.md...
```

---

## 2. Quy Định & Chu Kỳ Nộp Báo Cáo

1. **Chu kỳ bắt buộc:** Mỗi thành viên tự tạo và hoàn thiện **1 file log mỗi tuần** (ví dụ: `W02.md`, `W03.md`...).
2. **Hạn chót cập nhật:** **23:59 Thứ Tư hàng tuần** (đồng bộ cùng lịch nộp Báo cáo Tiến độ Tuần - Weekly Status Report theo mục 7.5 trong `AGENT.md`).
3. **Tính trung thực & minh bạch:**
   - Tuyệt đối không che giấu việc sử dụng AI. Việc sử dụng AI hợp lý để tối ưu năng suất được khuyến khích.
   - Điểm số đánh giá cao khả năng **Human-in-the-loop**: Bạn phát hiện lỗi của AI ở đâu, kiểm chứng như thế nào và tối ưu ra sao để khớp với yêu cầu của Portal67.
4. **Cách đặt tên file log tuần:**
   - Trong thư mục cá nhân: đặt tên `W01.md`, `W02.md`, `W03.md`, ... (chữ hoa `W` kèm số tuần có 2 chữ số).

---

## 3. Quy Trình 3 Bước Mỗi Tuần Dành Cho Thành Viên

1. **Đầu tuần:** Vào thư mục của mình, copy file `template.md` và đổi tên thành `W<XX>.md` (ví dụ `W03.md`). Điền mục tiêu tuần.
2. **Trong tuần (Ghi chép song song):** Mỗi khi dùng AI cho các tác vụ quan trọng (phân tích nghiệp vụ, vẽ sơ đồ, viết code, sinh test case), hãy dành 1-2 phút copy prompt và ghi chú ngắn gọn: AI làm đúng gì, làm sai gì, bạn đã sửa gì.
3. **Cuối tuần (Trước 23:59 Thứ Tư):** Rà soát lại, gắn link commit/PR đối chứng và tự đánh giá mức độ đóng góp trước khi commit lên branch cá nhân và tạo PR vào `main`.
