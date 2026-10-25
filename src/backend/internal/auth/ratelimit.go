package auth

import (
	"sync"
	"time"
)

// limiter is an in-memory fixed-window counter. Each key may be hit limit
// times per window; the window starts at the key's first hit. State is lost on
// restart and is not shared between processes, which is enough for a single
// API instance.
type limiter struct {
	limit  int
	window time.Duration
	now    func() time.Time

	mu        sync.Mutex
	hits      map[string]window
	lastSweep time.Time
}

type window struct {
	count   int
	resetAt time.Time
}

func newLimiter(limit int, period time.Duration) *limiter {
	return &limiter{limit: limit, window: period, now: time.Now, hits: make(map[string]window)}
}

// allow counts a hit for key. When the key is over its limit it returns false
// and how long until the window resets; rejected hits are not counted.
func (l *limiter) allow(key string) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	l.sweep(now)
	w, ok := l.hits[key]
	if !ok || !now.Before(w.resetAt) {
		w = window{resetAt: now.Add(l.window)}
	}
	if w.count >= l.limit {
		return false, w.resetAt.Sub(now)
	}
	w.count++
	l.hits[key] = w
	return true, 0
}

// reset forgets key, for example after a successful login.
func (l *limiter) reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.hits, key)
}

// sweep drops expired windows at most once per window so the map does not
// grow without bound. The caller must hold l.mu.
func (l *limiter) sweep(now time.Time) {
	if now.Sub(l.lastSweep) < l.window {
		return
	}
	l.lastSweep = now
	for key, w := range l.hits {
		if !now.Before(w.resetAt) {
			delete(l.hits, key)
		}
	}
}
