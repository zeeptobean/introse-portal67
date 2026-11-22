package redisbroker

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
)

func (b *Broker) streamKey(stream string) string { return b.prefix + ":ev:" + stream }

func (b *Broker) CreateGroup(ctx context.Context, stream, group string) error {
	err := b.client.XGroupCreateMkStream(ctx, b.streamKey(stream), group, "0").Err()
	if err != nil && strings.HasPrefix(err.Error(), "BUSYGROUP ") {
		return nil
	}
	return err
}

func (b *Broker) Read(ctx context.Context, stream, group, worker string) (events.Delivery, error) {
	const reclaimAfter = 5 * time.Second
	key := b.streamKey(stream)
	pending, err := b.client.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: key, Group: group, Start: "-", End: "+", Count: 1,
	}).Result()
	if err != nil {
		return events.Delivery{}, err
	}
	if len(pending) > 0 {
		entry := pending[0]
		if entry.Idle < reclaimAfter {
			select {
			case <-ctx.Done():
				return events.Delivery{}, ctx.Err()
			case <-time.After(reclaimAfter - entry.Idle):
			}
			return events.Delivery{}, nil
		}
		messages, _, err := b.client.XAutoClaim(ctx, &redis.XAutoClaimArgs{
			Stream: key, Group: group, Consumer: worker, Start: entry.ID, MinIdle: reclaimAfter, Count: 1,
		}).Result()
		if err != nil || len(messages) == 0 {
			return events.Delivery{}, err
		}
		return delivery(stream, messages[0], entry.RetryCount+1), nil
	}
	streams, err := b.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group: group, Consumer: worker, Streams: []string{key, ">"}, Count: 1, Block: time.Second,
	}).Result()
	if errors.Is(err, redis.Nil) {
		return events.Delivery{}, nil
	}
	if err != nil {
		return events.Delivery{}, err
	}
	if len(streams) == 0 || len(streams[0].Messages) == 0 {
		return events.Delivery{}, nil
	}
	return delivery(stream, streams[0].Messages[0], 1), nil
}

func delivery(stream string, message redis.XMessage, attempts int64) events.Delivery {
	kind, _ := message.Values["type"].(string)
	payload, _ := message.Values["payload"].(string)
	return events.Delivery{
		Message: events.Message{Stream: stream, Type: kind, Payload: payload},
		ID:      message.ID, Attempts: attempts,
	}
}

func (b *Broker) Ack(ctx context.Context, delivery events.Delivery, group string) error {
	return b.client.XAck(ctx, b.streamKey(delivery.Stream), group, delivery.ID).Err()
}

var sendToDeadLetter = redis.NewScript(`
if #redis.call('XPENDING', KEYS[1], ARGV[1], ARGV[2], ARGV[2], 1) == 0 then return 0 end
redis.call('XADD', KEYS[2], '*', 'stream', ARGV[3], 'entry_id', ARGV[2],
           'type', ARGV[4], 'payload', ARGV[5], 'error', ARGV[6])
return redis.call('XACK', KEYS[1], ARGV[1], ARGV[2])
`)

func (b *Broker) DeadLetter(ctx context.Context, delivery events.Delivery, group, reason string) error {
	return sendToDeadLetter.Run(ctx, b.client, []string{b.streamKey(delivery.Stream), b.prefix + ":dlq:" + group},
		group, delivery.ID, delivery.Stream, delivery.Type, delivery.Payload, reason).Err()
}

var retryDeadLetter = redis.NewScript(`
local entries = redis.call('XRANGE', KEYS[1], ARGV[1], ARGV[1])
if #entries == 0 then return 0 end
local fields = {}
for i = 1, #entries[1][2], 2 do fields[entries[1][2][i]] = entries[1][2][i+1] end
if fields.stream ~= ARGV[2] then return redis.error_reply('entry belongs to another stream') end
redis.call('XADD', KEYS[2], '*', 'type', fields.type, 'payload', fields.payload)
return redis.call('XDEL', KEYS[1], ARGV[1])
`)

func (b *Broker) RetryDeadLetter(ctx context.Context, stream, group, entryID string) (bool, error) {
	count, err := retryDeadLetter.Run(ctx, b.client, []string{b.prefix + ":dlq:" + group, b.streamKey(stream)}, entryID, stream).Int()
	return count == 1, err
}
