package gateway

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/authn"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/events"
)

func userVersionKey(prefix string, id int64) string {
	if prefix == "" {
		prefix = "p67"
	}
	return fmt.Sprintf("%s:sess-user-version:%d", prefix, id)
}

func (g *Gateway) versionCurrent(ctx context.Context, session authn.Session) bool {
	call, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	version, err := g.cache.Get(call, userVersionKey(g.prefix, session.UserID)).Int64()
	return errors.Is(err, redis.Nil) || (err == nil && version <= session.Version)
}

var advanceUserVersion = redis.NewScript(`
local current = redis.call('GET', KEYS[1])
if not current or tonumber(current) < tonumber(ARGV[1]) then
  redis.call('SET', KEYS[1], ARGV[1])
end
return 1
`)

// InvalidateCachedUsers advances a Redis version for each directory event.
// The update is idempotent, including if the gateway stops before acking.
func InvalidateCachedUsers(client *redis.Client, prefix string) func(context.Context, events.Event) error {
	return func(ctx context.Context, event events.Event) error {
		if event.Type != "identity.user.upserted" && event.Type != "identity.user.deactivated" {
			return nil
		}
		id, err := strconv.ParseInt(event.Aggregate.ID, 10, 64)
		if err != nil || id <= 0 || event.Aggregate.Version <= 0 || event.SchemaVersion != 1 {
			return errors.New("invalid directory event version")
		}
		return advanceUserVersion.Run(ctx, client, []string{userVersionKey(prefix, id)}, event.Aggregate.Version).Err()
	}
}
