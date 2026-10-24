-- Demo document requests: one pending, one approved, one rejected.
-- Requires 00001_demo_users.sql.
-- +goose Up
INSERT INTO document_requests (student_id, request_type, purpose, status, submitted_at)
SELECT u.id, r.request_type, r.purpose, r.status, r.submitted_at::timestamptz
FROM (VALUES
    ('student1@portal67.test', 'enrollment_confirmation', 'Apply for a student bus pass', 'pending', '2026-10-20 09:15+07'),
    ('student1@portal67.test', 'transcript', 'Scholarship application', 'approved', '2026-10-10 14:00+07'),
    ('student2@portal67.test', 'other', 'Confirmation of exam absence', 'rejected', '2026-10-12 08:30+07')
) AS r (student_email, request_type, purpose, status, submitted_at)
JOIN users u ON u.email = r.student_email
WHERE NOT EXISTS (
    SELECT 1 FROM document_requests d WHERE d.student_id = u.id AND d.purpose = r.purpose
);

INSERT INTO document_request_events (document_request_id, actor_id, from_status, to_status, comment, created_at)
SELECT d.id, d.student_id, NULL, 'pending', NULL, d.submitted_at
FROM document_requests d
JOIN users u ON u.id = d.student_id
WHERE u.email IN ('student1@portal67.test', 'student2@portal67.test')
ON CONFLICT DO NOTHING;

INSERT INTO document_request_events (document_request_id, actor_id, from_status, to_status, comment, created_at)
SELECT d.id, staff.id, 'pending', d.status, e.comment, d.submitted_at + interval '2 days'
FROM (VALUES
    ('Scholarship application', 'Transcript is ready for pickup at the faculty office.'),
    ('Confirmation of exam absence', 'Please attach a medical certificate and submit a new request.')
) AS e (purpose, comment)
JOIN document_requests d ON d.purpose = e.purpose AND d.status <> 'pending'
JOIN users staff ON staff.email = 'staff1@portal67.test'
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM document_request_events WHERE document_request_id IN (
    SELECT d.id FROM document_requests d JOIN users u ON u.id = d.student_id
    WHERE u.email IN ('student1@portal67.test', 'student2@portal67.test')
);
DELETE FROM document_requests WHERE student_id IN (
    SELECT id FROM users WHERE email IN ('student1@portal67.test', 'student2@portal67.test')
);
