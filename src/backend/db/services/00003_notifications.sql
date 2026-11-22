-- +goose Up
CREATE TABLE notification.notifications (
    event_id UUID NOT NULL,
    recipient_id BIGINT NOT NULL REFERENCES notification.users(id),
    kind TEXT NOT NULL,
    subject_id TEXT NOT NULL,
    message TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (event_id, recipient_id)
);
CREATE INDEX notifications_recent_idx
    ON notification.notifications (recipient_id, occurred_at DESC, event_id DESC);
GRANT SELECT, INSERT, UPDATE, DELETE ON notification.notifications TO portal67_notification;

-- +goose Down
DROP TABLE notification.notifications;
