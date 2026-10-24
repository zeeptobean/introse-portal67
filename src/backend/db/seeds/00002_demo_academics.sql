-- Demo academic records: a finished term with grades and an open term with
-- available, full, and clashing sections. Requires 00001_demo_users.sql.
-- +goose Up
INSERT INTO grading_scales (name) VALUES ('Sample 10-point scale') ON CONFLICT DO NOTHING;

INSERT INTO grade_bands (grading_scale_id, lower_bound, letter_grade, grade_points)
SELECT gs.id, b.lower_bound, b.letter_grade, b.grade_points
FROM grading_scales gs
CROSS JOIN (VALUES
    (0.00, 'F', 0.0),
    (4.00, 'D', 1.0),
    (5.00, 'C', 2.0),
    (6.50, 'B', 3.0),
    (8.00, 'A', 3.5),
    (9.00, 'A+', 4.0)
) AS b (lower_bound, letter_grade, grade_points)
WHERE gs.name = 'Sample 10-point scale'
ON CONFLICT DO NOTHING;

INSERT INTO terms (code, name, starts_on, ends_on, grading_scale_id)
SELECT t.code, t.name, t.starts_on::date, t.ends_on::date, gs.id
FROM grading_scales gs
CROSS JOIN (VALUES
    ('2026-1', 'Semester 1, 2026', '2026-02-02', '2026-06-28'),
    ('2026-2', 'Semester 2, 2026', '2026-09-07', '2027-01-17')
) AS t (code, name, starts_on, ends_on)
WHERE gs.name = 'Sample 10-point scale'
ON CONFLICT DO NOTHING;

INSERT INTO courses (code, name, credits, kind) VALUES
    ('CSC10001', 'Introduction to Programming', 4, 'course'),
    ('CSC10002', 'Programming Techniques', 4, 'course'),
    ('CSC10003', 'Object-Oriented Programming', 4, 'course'),
    ('CSC13002', 'Introduction to Software Engineering', 4, 'course'),
    ('MTH00003', 'Calculus 1B', 3, 'course'),
    ('CSC14005', 'Introduction to Machine Learning', 4, 'specialized_topic'),
    ('CSC15001', 'Special Topic: Blockchain Applications', 3, 'specialized_topic')
ON CONFLICT DO NOTHING;

-- CSC13002 section 02 has one seat and is full; CSC10003 01 clashes with CSC13002 01.
INSERT INTO sections (course_id, term_id, lecturer_id, section_code, capacity)
SELECT c.id, t.id, u.id, s.section_code, s.capacity
FROM (VALUES
    ('CSC10001', '2026-1', 'lecturer1@portal67.test', '01', 40),
    ('CSC10002', '2026-1', 'lecturer2@portal67.test', '01', 40),
    ('MTH00003', '2026-1', 'lecturer2@portal67.test', '01', 60),
    ('CSC13002', '2026-2', 'lecturer1@portal67.test', '01', 40),
    ('CSC13002', '2026-2', 'lecturer2@portal67.test', '02', 1),
    ('CSC10003', '2026-2', 'lecturer1@portal67.test', '01', 40),
    ('MTH00003', '2026-2', 'lecturer2@portal67.test', '01', 60),
    ('CSC14005', '2026-2', 'lecturer2@portal67.test', '01', 30),
    ('CSC15001', '2026-2', 'lecturer1@portal67.test', '01', 25)
) AS s (course_code, term_code, lecturer_email, section_code, capacity)
JOIN courses c ON c.code = s.course_code
JOIN terms t ON t.code = s.term_code
JOIN users u ON u.email = s.lecturer_email
ON CONFLICT DO NOTHING;

INSERT INTO class_meetings (section_id, weekday, starts_at, ends_at, room)
SELECT sec.id, m.weekday, m.starts_at::time, m.ends_at::time, m.room
FROM (VALUES
    ('CSC10001', '2026-1', '01', 2, '07:30', '11:30', 'E101'),
    ('CSC10002', '2026-1', '01', 4, '13:30', '17:30', 'E202'),
    ('MTH00003', '2026-1', '01', 5, '07:30', '10:00', 'F301'),
    ('CSC13002', '2026-2', '01', 1, '07:30', '11:30', 'F102'),
    ('CSC13002', '2026-2', '02', 3, '13:30', '17:30', 'F103'),
    ('CSC10003', '2026-2', '01', 1, '09:30', '11:30', 'E203'),
    ('CSC10003', '2026-2', '01', 4, '07:30', '09:30', 'E203'),
    ('MTH00003', '2026-2', '01', 4, '13:30', '15:30', 'F301'),
    ('MTH00003', '2026-2', '01', 6, '07:30', '09:30', 'F301'),
    ('CSC14005', '2026-2', '01', 2, '07:30', '11:30', 'Lab I.41'),
    ('CSC15001', '2026-2', '01', 6, '13:30', '16:30', 'E305')
) AS m (course_code, term_code, section_code, weekday, starts_at, ends_at, room)
JOIN courses c ON c.code = m.course_code
JOIN terms t ON t.code = m.term_code
JOIN sections sec ON sec.course_id = c.id AND sec.term_id = t.id AND sec.section_code = m.section_code
ON CONFLICT DO NOTHING;

-- The CSC13002 and MTH00003 finals overlap to demonstrate exam conflicts.
INSERT INTO exams (section_id, label, starts_at, ends_at, room)
SELECT sec.id, x.label, x.starts_at::timestamptz, x.ends_at::timestamptz, x.room
FROM (VALUES
    ('CSC13002', '01', 'Final exam', '2027-01-05 08:00+07', '2027-01-05 10:00+07', 'F102'),
    ('MTH00003', '01', 'Final exam', '2027-01-05 09:00+07', '2027-01-05 10:30+07', 'F301'),
    ('CSC10003', '01', 'Final exam', '2027-01-07 13:30+07', '2027-01-07 15:30+07', 'E203'),
    ('CSC14005', '01', 'Project defense', '2027-01-08 07:30+07', '2027-01-08 11:30+07', 'Lab I.41')
) AS x (course_code, section_code, label, starts_at, ends_at, room)
JOIN courses c ON c.code = x.course_code
JOIN terms t ON t.code = '2026-2'
JOIN sections sec ON sec.course_id = c.id AND sec.term_id = t.id AND sec.section_code = x.section_code
WHERE NOT EXISTS (SELECT 1 FROM exams e WHERE e.section_id = sec.id AND e.label = x.label);

INSERT INTO enrollments (student_id, section_id, registered_at)
SELECT u.id, sec.id, e.registered_at::timestamptz
FROM (VALUES
    ('student1@portal67.test', 'CSC10001', '2026-1', '01', '2026-01-15 09:00+07'),
    ('student1@portal67.test', 'CSC10002', '2026-1', '01', '2026-01-15 09:05+07'),
    ('student1@portal67.test', 'MTH00003', '2026-1', '01', '2026-01-15 09:10+07'),
    ('student2@portal67.test', 'CSC10001', '2026-1', '01', '2026-01-16 10:00+07'),
    ('student2@portal67.test', 'CSC10002', '2026-1', '01', '2026-01-16 10:05+07'),
    ('student1@portal67.test', 'CSC13002', '2026-2', '01', '2026-08-20 08:00+07'),
    ('student1@portal67.test', 'MTH00003', '2026-2', '01', '2026-08-20 08:05+07'),
    ('student2@portal67.test', 'CSC13002', '2026-2', '02', '2026-08-21 08:00+07')
) AS e (student_email, course_code, term_code, section_code, registered_at)
JOIN users u ON u.email = e.student_email
JOIN courses c ON c.code = e.course_code
JOIN terms t ON t.code = e.term_code
JOIN sections sec ON sec.course_id = c.id AND sec.term_id = t.id AND sec.section_code = e.section_code
ON CONFLICT DO NOTHING;

-- student1 failed MTH00003; student2's CSC10002 enrollment stays ungraded.
INSERT INTO grades (enrollment_id, score, graded_by, updated_at)
SELECT en.id, g.score, sec.lecturer_id, '2026-07-05 10:00+07'
FROM (VALUES
    ('student1@portal67.test', 'CSC10001', 8.50),
    ('student1@portal67.test', 'CSC10002', 6.75),
    ('student1@portal67.test', 'MTH00003', 3.00),
    ('student2@portal67.test', 'CSC10001', 9.20)
) AS g (student_email, course_code, score)
JOIN users u ON u.email = g.student_email
JOIN courses c ON c.code = g.course_code
JOIN terms t ON t.code = '2026-1'
JOIN sections sec ON sec.course_id = c.id AND sec.term_id = t.id AND sec.section_code = '01'
JOIN enrollments en ON en.student_id = u.id AND en.section_id = sec.id
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM grades WHERE enrollment_id IN (
    SELECT e.id FROM enrollments e
    JOIN sections s ON s.id = e.section_id
    JOIN terms t ON t.id = s.term_id
    WHERE t.code IN ('2026-1', '2026-2')
);
DELETE FROM enrollments WHERE section_id IN (
    SELECT s.id FROM sections s JOIN terms t ON t.id = s.term_id WHERE t.code IN ('2026-1', '2026-2')
);
DELETE FROM exams WHERE section_id IN (
    SELECT s.id FROM sections s JOIN terms t ON t.id = s.term_id WHERE t.code IN ('2026-1', '2026-2')
);
DELETE FROM class_meetings WHERE section_id IN (
    SELECT s.id FROM sections s JOIN terms t ON t.id = s.term_id WHERE t.code IN ('2026-1', '2026-2')
);
DELETE FROM sections WHERE term_id IN (SELECT id FROM terms WHERE code IN ('2026-1', '2026-2'));
DELETE FROM courses WHERE code IN ('CSC10001', 'CSC10002', 'CSC10003', 'CSC13002', 'MTH00003', 'CSC14005', 'CSC15001');
DELETE FROM terms WHERE code IN ('2026-1', '2026-2');
DELETE FROM grade_bands WHERE grading_scale_id IN (SELECT id FROM grading_scales WHERE name = 'Sample 10-point scale');
DELETE FROM grading_scales WHERE name = 'Sample 10-point scale';
