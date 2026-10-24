-- Demo accounts. Every account's password is "portal67".
-- Apply with: goose -dir db/seeds -no-versioning up
-- +goose Up
INSERT INTO users (email, password_hash, institutional_id, full_name, date_of_birth, role)
VALUES
    ('student1@portal67.test', '$argon2id$v=19$m=65536,t=1,p=20$7Q2+UCMMTnRoU8iL/1rs5A$pmL/C19I/192haTGERtuJtlJfTzPcb0fyMrD96abRtk',
        '24120001', 'Nguyen Van An', '2006-03-14', 'student'),
    ('student2@portal67.test', '$argon2id$v=19$m=65536,t=1,p=20$7Q2+UCMMTnRoU8iL/1rs5A$pmL/C19I/192haTGERtuJtlJfTzPcb0fyMrD96abRtk',
        '24120002', 'Tran Thi Binh', '2006-07-02', 'student'),
    ('lecturer1@portal67.test', '$argon2id$v=19$m=65536,t=1,p=20$7Q2+UCMMTnRoU8iL/1rs5A$pmL/C19I/192haTGERtuJtlJfTzPcb0fyMrD96abRtk',
        'GV0001', 'Le Hoang Cuong', '1985-11-20', 'lecturer'),
    ('lecturer2@portal67.test', '$argon2id$v=19$m=65536,t=1,p=20$7Q2+UCMMTnRoU8iL/1rs5A$pmL/C19I/192haTGERtuJtlJfTzPcb0fyMrD96abRtk',
        'GV0002', 'Pham Minh Dung', '1979-05-08', 'lecturer'),
    ('staff1@portal67.test', '$argon2id$v=19$m=65536,t=1,p=20$7Q2+UCMMTnRoU8iL/1rs5A$pmL/C19I/192haTGERtuJtlJfTzPcb0fyMrD96abRtk',
        'NV0001', 'Vo Thi Em', '1990-01-25', 'staff')
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM users WHERE email IN (
    'student1@portal67.test', 'student2@portal67.test',
    'lecturer1@portal67.test', 'lecturer2@portal67.test', 'staff1@portal67.test'
);
