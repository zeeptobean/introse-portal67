package auth

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	now := time.Date(2026, 10, 25, 9, 0, 0, 0, time.UTC)
	l := newLimiter(2, time.Minute)
	l.now = func() time.Time { return now }

	for i := range 2 {
		if ok, _ := l.allow("a"); !ok {
			t.Fatalf("hit %d rejected, want allowed", i+1)
		}
	}
	ok, retry := l.allow("a")
	if ok || retry != time.Minute {
		t.Fatalf("third hit = (%v, %v), want (false, 1m)", ok, retry)
	}
	if ok, _ := l.allow("b"); !ok {
		t.Fatal("other key rejected, want allowed")
	}

	now = now.Add(30 * time.Second)
	if ok, retry := l.allow("a"); ok || retry != 30*time.Second {
		t.Fatalf("hit after 30s = (%v, %v), want (false, 30s)", ok, retry)
	}

	now = now.Add(30 * time.Second)
	if ok, _ := l.allow("a"); !ok {
		t.Fatal("hit after window rejected, want allowed")
	}

	l.allow("a")
	l.reset("a")
	if ok, _ := l.allow("a"); !ok {
		t.Fatal("hit after reset rejected, want allowed")
	}
}

func TestLimiterSweepsExpiredKeys(t *testing.T) {
	now := time.Date(2026, 10, 25, 9, 0, 0, 0, time.UTC)
	l := newLimiter(1, time.Minute)
	l.now = func() time.Time { return now }
	l.allow("a")
	l.allow("b")
	now = now.Add(2 * time.Minute)
	l.allow("c")
	if len(l.hits) != 1 {
		t.Fatalf("len(hits) = %d after sweep, want 1", len(l.hits))
	}
}
