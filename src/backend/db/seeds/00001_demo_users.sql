-- Demo accounts. Every account's password is "portal67".
-- Apply with: goose -dir db/seeds -no-versioning up
-- +goose Up
INSERT INTO users (email, password_hash, institutional_id, full_name, date_of_birth, role)
VALUES
    ('student1@portal67.test', '$argon2id$v=19$m=19456,t=2,p=1$0TTtpJa/3OFfnMMNKo2kIA$sYElyVCoNXVDN0r3j0mOr35tRwZ1DauK51LAJO+zCe8',
        '24120001', 'Nguyen Van An', '2006-03-14', 'student'),
    ('student2@portal67.test', '$argon2id$v=19$m=19456,t=2,p=1$0TTtpJa/3OFfnMMNKo2kIA$sYElyVCoNXVDN0r3j0mOr35tRwZ1DauK51LAJO+zCe8',
        '24120002', 'Tran Thi Binh', '2006-07-02', 'student'),
    ('lecturer1@portal67.test', '$argon2id$v=19$m=19456,t=2,p=1$0TTtpJa/3OFfnMMNKo2kIA$sYElyVCoNXVDN0r3j0mOr35tRwZ1DauK51LAJO+zCe8',
        'GV0001', 'Le Hoang Cuong', '1985-11-20', 'lecturer'),
    ('lecturer2@portal67.test', '$argon2id$v=19$m=19456,t=2,p=1$0TTtpJa/3OFfnMMNKo2kIA$sYElyVCoNXVDN0r3j0mOr35tRwZ1DauK51LAJO+zCe8',
        'GV0002', 'Pham Minh Dung', '1979-05-08', 'lecturer'),
    ('staff1@portal67.test', '$argon2id$v=19$m=19456,t=2,p=1$0TTtpJa/3OFfnMMNKo2kIA$sYElyVCoNXVDN0r3j0mOr35tRwZ1DauK51LAJO+zCe8',
        'NV0001', 'Vo Thi Em', '1990-01-25', 'staff')
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM users WHERE email IN (
    'student1@portal67.test', 'student2@portal67.test',
    'lecturer1@portal67.test', 'lecturer2@portal67.test', 'staff1@portal67.test'
);
