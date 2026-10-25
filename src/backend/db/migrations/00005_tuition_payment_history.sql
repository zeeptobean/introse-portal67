-- +goose Up
-- updated_by is the staff member behind the latest change to the payment.
ALTER TABLE tuition_payments ADD COLUMN updated_by BIGINT REFERENCES users (id);
UPDATE tuition_payments SET updated_by = recorded_by;
ALTER TABLE tuition_payments ALTER COLUMN updated_by SET NOT NULL;

-- Append-only history. Each row is the payment as it was right after the
-- change made by actor_id. 'snapshot' rows hold the state of payments that
-- existed before history was tracked; their earlier edits are unknown.
CREATE TABLE tuition_payment_events (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    payment_id BIGINT NOT NULL REFERENCES tuition_payments (id),
    actor_id BIGINT NOT NULL REFERENCES users (id),
    action TEXT NOT NULL CHECK (action IN ('snapshot', 'recorded', 'updated', 'voided')),
    amount_vnd BIGINT NOT NULL CHECK (amount_vnd > 0),
    status TEXT NOT NULL CHECK (status IN ('recorded', 'voided')),
    paid_at TIMESTAMPTZ NOT NULL,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (action <> 'voided' OR status = 'voided')
);

CREATE INDEX tuition_payment_events_history_idx
    ON tuition_payment_events (payment_id, created_at, id);

INSERT INTO tuition_payment_events (payment_id, actor_id, action, amount_vnd, status, paid_at, note, created_at)
SELECT id, updated_by, 'snapshot', amount_vnd, status, paid_at, note, updated_at
FROM tuition_payments
ORDER BY id;

-- +goose Down
DROP TABLE tuition_payment_events;
ALTER TABLE tuition_payments DROP COLUMN updated_by;
