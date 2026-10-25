-- Demo tuition records covering paid, partial (with a voided payment), credit,
-- and unpaid balances. Requires 00001_demo_users.sql and 00002_demo_academics.sql.
-- +goose Up
INSERT INTO tuition_charges (student_id, term_id, description, amount_vnd, created_by, created_at)
SELECT u.id, t.id, c.description, c.amount_vnd, staff.id, c.created_at::timestamptz
FROM (VALUES
    ('student1@portal67.test', '2026-1', 'Tuition (11 credits)', 8250000, '2026-01-20 09:00+07'),
    ('student1@portal67.test', '2026-2', 'Tuition (11 credits)', 8250000, '2026-08-25 09:00+07'),
    ('student1@portal67.test', '2026-2', 'Facility fee', 500000, '2026-08-25 09:00+07'),
    ('student2@portal67.test', '2026-1', 'Tuition (8 credits)', 6000000, '2026-01-20 09:00+07'),
    ('student2@portal67.test', '2026-2', 'Tuition (4 credits)', 3000000, '2026-08-25 09:00+07')
) AS c (student_email, term_code, description, amount_vnd, created_at)
JOIN users u ON u.email = c.student_email
JOIN terms t ON t.code = c.term_code
JOIN users staff ON staff.email = 'staff1@portal67.test'
WHERE NOT EXISTS (
    SELECT 1 FROM tuition_charges x
    WHERE x.student_id = u.id AND x.term_id = t.id AND x.description = c.description
);

-- The voided payment was voided one hour after it was recorded.
INSERT INTO tuition_payments (student_id, term_id, amount_vnd, status, recorded_by, updated_by, paid_at, note, updated_at)
SELECT u.id, t.id, p.amount_vnd, p.status, staff.id, staff.id, p.paid_at::timestamptz, p.note,
    p.paid_at::timestamptz + CASE WHEN p.status = 'voided' THEN interval '1 hour' ELSE interval '0' END
FROM (VALUES
    ('student1@portal67.test', '2026-1', 8250000, 'recorded', '2026-02-10 10:00+07', 'Bank transfer'),
    ('student1@portal67.test', '2026-2', 3000000, 'recorded', '2026-09-15 10:00+07', 'First installment'),
    ('student1@portal67.test', '2026-2', 1000000, 'voided', '2026-09-20 10:00+07', 'Duplicate entry'),
    ('student2@portal67.test', '2026-1', 6500000, 'recorded', '2026-02-12 14:30+07', 'Overpaid by transfer')
) AS p (student_email, term_code, amount_vnd, status, paid_at, note)
JOIN users u ON u.email = p.student_email
JOIN terms t ON t.code = p.term_code
JOIN users staff ON staff.email = 'staff1@portal67.test'
WHERE NOT EXISTS (
    SELECT 1 FROM tuition_payments x
    WHERE x.student_id = u.id AND x.term_id = t.id AND x.note = p.note
);

INSERT INTO tuition_payment_events (payment_id, actor_id, action, amount_vnd, status, paid_at, note, created_at)
SELECT p.id, p.recorded_by, e.action, p.amount_vnd, e.status, p.paid_at, p.note, e.created_at
FROM tuition_payments p
CROSS JOIN LATERAL (VALUES
    ('recorded', 'recorded', p.paid_at),
    ('voided', 'voided', p.updated_at)
) AS e (action, status, created_at)
WHERE p.student_id IN (SELECT id FROM users WHERE email IN ('student1@portal67.test', 'student2@portal67.test'))
    AND (e.action = 'recorded' OR p.status = 'voided')
    AND NOT EXISTS (SELECT 1 FROM tuition_payment_events x WHERE x.payment_id = p.id);

-- +goose Down
DELETE FROM tuition_payment_events WHERE payment_id IN (
    SELECT p.id FROM tuition_payments p
    JOIN users u ON u.id = p.student_id
    WHERE u.email IN ('student1@portal67.test', 'student2@portal67.test')
);
DELETE FROM tuition_payments WHERE student_id IN (
    SELECT id FROM users WHERE email IN ('student1@portal67.test', 'student2@portal67.test')
);
DELETE FROM tuition_charges WHERE student_id IN (
    SELECT id FROM users WHERE email IN ('student1@portal67.test', 'student2@portal67.test')
);
