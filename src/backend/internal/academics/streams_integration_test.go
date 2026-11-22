package academics

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events/redisbroker"
)

func TestRedisConsumerDeduplicatesReplay(t *testing.T) {
	redisURL := os.Getenv("PORTAL67_TEST_REDIS_URL")
	if redisURL == "" {
		t.Skip("set PORTAL67_TEST_REDIS_URL to test Redis Streams")
	}
	ctx, pool := registrationTestDB(t)
	prefix := "p67:test:" + uuid.NewString()
	broker, err := redisbroker.Open(redisURL, prefix)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = broker.Close() })
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		t.Fatal("invalid test Redis URL")
	}
	reader := redis.NewClient(options)
	t.Cleanup(func() {
		if err := reader.Del(ctx, prefix+":ev:registration", prefix+":dlq:academic").Err(); err != nil {
			t.Errorf("remove test streams: %v", err)
		}
		_ = reader.Close()
	})
	sectionID, users := seedRegistrationSection(t, ctx, pool, 1)
	if _, err := NewService(NewStore(pool)).Register(ctx, users[0], sectionID); err != nil {
		t.Fatal(err)
	}
	relay := events.NewRelay(pool, broker)
	for range 2 {
		if _, err := relay.Flush(ctx); err != nil {
			t.Fatal(err)
		}
		if _, err := pool.Exec(ctx, "UPDATE outbox SET published_at = NULL"); err != nil {
			t.Fatal(err)
		}
	}
	calls := 0
	consumer := events.Consumer{
		Broker: broker, Stream: "registration", Group: "academic", Worker: "test",
		Apply: events.Transactional(pool, "registration", func(_ context.Context, _ pgx.Tx, event events.Event) error {
			if event.Type != "registration.enrollment.created" {
				t.Errorf("unexpected event type %s", event.Type)
			}
			calls++
			return nil
		}),
	}
	if err := broker.CreateGroup(ctx, consumer.Stream, consumer.Group); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if err := consumer.Step(ctx); err != nil {
			t.Fatal(err)
		}
	}
	if calls != 1 {
		t.Fatalf("projection applied %d times", calls)
	}
	if err := broker.Publish(ctx, []events.Message{{Stream: consumer.Stream, Type: "test.invalid", Payload: "invalid"}}); err != nil {
		t.Fatal(err)
	}
	delivery, err := broker.Read(ctx, consumer.Stream, consumer.Group, consumer.Worker)
	if err != nil {
		t.Fatal(err)
	}
	if delivery.Type != "test.invalid" {
		t.Fatalf("unexpected delivery %q", delivery.Type)
	}
	if err := broker.DeadLetter(ctx, delivery, consumer.Group, "invalid payload"); err != nil {
		t.Fatal(err)
	}
	entries, err := reader.XRange(ctx, prefix+":dlq:academic", "-", "+").Result()
	if err != nil || len(entries) != 1 {
		t.Fatalf("dead-letter entries = %d, err = %v", len(entries), err)
	}
	if retried, err := broker.RetryDeadLetter(ctx, consumer.Stream, consumer.Group, entries[0].ID); err != nil || !retried {
		t.Fatalf("retry DLQ entry = %t, err = %v", retried, err)
	}
	retried, err := broker.Read(ctx, consumer.Stream, consumer.Group, consumer.Worker)
	if err != nil || retried.Type != delivery.Type || retried.Payload != delivery.Payload {
		t.Fatalf("retried delivery differs: %v", err)
	}
}
