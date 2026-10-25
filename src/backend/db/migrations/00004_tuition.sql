-- +goose Up
-- Whole VND amounts. Charges are entered per student and term, not generated
-- from registration. Payments are simulations; no money is transferred.
CREATE TABLE tuition_charges (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users (id),
    term_id BIGINT NOT NULL REFERENCES terms (id),
    description TEXT NOT NULL CHECK (btrim(description) <> ''),
    amount_vnd BIGINT NOT NULL CHECK (amount_vnd > 0),
    created_by BIGINT NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX tuition_charges_student_term_idx ON tuition_charges (student_id, term_id);

-- A voided payment stays visible but no longer counts toward the paid total.
CREATE TABLE tuition_payments (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users (id),
    term_id BIGINT NOT NULL REFERENCES terms (id),
    amount_vnd BIGINT NOT NULL CHECK (amount_vnd > 0),
    status TEXT NOT NULL DEFAULT 'recorded' CHECK (status IN ('recorded', 'voided')),
    recorded_by BIGINT NOT NULL REFERENCES users (id),
    paid_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    note TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX tuition_payments_student_term_idx ON tuition_payments (student_id, term_id);

-- +goose Down
DROP TABLE tuition_payments;
DROP TABLE tuition_charges;
