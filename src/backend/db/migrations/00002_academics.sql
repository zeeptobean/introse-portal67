-- +goose Up
-- A scale is frozen once referenced by graded enrollments. New policy, new scale.
CREATE TABLE grading_scales (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CHECK (btrim(name) <> '')
);

-- The band with the greatest lower_bound <= the raw score applies. Each usable
-- scale needs a 0.00 band; the seed and service are responsible for that.
CREATE TABLE grade_bands (
    grading_scale_id BIGINT NOT NULL REFERENCES grading_scales (id),
    lower_bound NUMERIC(4,2) NOT NULL CHECK (lower_bound BETWEEN 0 AND 10),
    letter_grade TEXT NOT NULL CHECK (btrim(letter_grade) <> ''),
    grade_points NUMERIC(3,2) NOT NULL CHECK (grade_points BETWEEN 0 AND 4),
    PRIMARY KEY (grading_scale_id, lower_bound)
);

CREATE TABLE terms (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code TEXT NOT NULL UNIQUE CHECK (btrim(code) <> ''),
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    starts_on DATE NOT NULL,
    ends_on DATE NOT NULL,
    grading_scale_id BIGINT NOT NULL REFERENCES grading_scales (id),
    CHECK (ends_on >= starts_on)
);

CREATE TABLE courses (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code TEXT NOT NULL UNIQUE CHECK (btrim(code) <> ''),
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    credits SMALLINT NOT NULL CHECK (credits > 0),
    kind TEXT NOT NULL DEFAULT 'course' CHECK (kind IN ('course', 'specialized_topic'))
);

CREATE TABLE sections (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES courses (id),
    term_id BIGINT NOT NULL REFERENCES terms (id),
    lecturer_id BIGINT NOT NULL REFERENCES users (id),
    section_code TEXT NOT NULL CHECK (btrim(section_code) <> ''),
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    UNIQUE (term_id, course_id, section_code)
);

CREATE INDEX sections_course_idx ON sections (course_id);
CREATE INDEX sections_lecturer_idx ON sections (lecturer_id);

-- Weekly local campus times (Asia/Ho_Chi_Minh) repeated within the term dates.
-- Monday = 1, Sunday = 7.
CREATE TABLE class_meetings (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    section_id BIGINT NOT NULL REFERENCES sections (id),
    weekday SMALLINT NOT NULL CHECK (weekday BETWEEN 1 AND 7),
    starts_at TIME NOT NULL,
    ends_at TIME NOT NULL,
    room TEXT,
    CHECK (ends_at > starts_at),
    UNIQUE (section_id, weekday, starts_at, ends_at)
);

CREATE TABLE enrollments (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users (id),
    section_id BIGINT NOT NULL REFERENCES sections (id),
    status TEXT NOT NULL DEFAULT 'enrolled' CHECK (status IN ('enrolled', 'cancelled')),
    registered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    cancelled_at TIMESTAMPTZ,
    UNIQUE (student_id, section_id),
    CHECK (
        (status = 'enrolled' AND cancelled_at IS NULL)
        OR (status = 'cancelled' AND cancelled_at IS NOT NULL AND cancelled_at >= registered_at)
    )
);

CREATE INDEX enrollments_active_section_idx ON enrollments (section_id) WHERE status = 'enrolled';

CREATE TABLE exams (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    section_id BIGINT NOT NULL REFERENCES sections (id),
    label TEXT NOT NULL CHECK (btrim(label) <> ''),
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    room TEXT,
    CHECK (ends_at > starts_at)
);

CREATE INDEX exams_section_idx ON exams (section_id);

-- No grade row means ungraded. A score of zero is a real grade.
CREATE TABLE grades (
    enrollment_id BIGINT PRIMARY KEY REFERENCES enrollments (id),
    score NUMERIC(4,2) NOT NULL CHECK (score BETWEEN 0 AND 10),
    graded_by BIGINT NOT NULL REFERENCES users (id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE grades;
DROP TABLE exams;
DROP TABLE enrollments;
DROP TABLE class_meetings;
DROP TABLE sections;
DROP TABLE courses;
DROP TABLE terms;
DROP TABLE grade_bands;
DROP TABLE grading_scales;
