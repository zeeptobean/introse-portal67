package academics

import (
	"errors"
	"testing"
	"time"
)

func TestCheckRegistrationWindow(t *testing.T) {
	opens := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	closes := time.Date(2026, 9, 21, 0, 0, 0, 0, time.UTC)
	sec := lockedSection{RegistrationOpensAt: opens, RegistrationClosesAt: closes}
	tests := []struct {
		now  time.Time
		want error
	}{
		{opens.Add(-time.Second), ErrRegistrationNotOpen},
		{opens, nil},
		{closes.Add(-time.Second), nil},
		{closes, ErrRegistrationClosed},
	}
	for _, tt := range tests {
		if got := checkRegistrationWindow(sec, tt.now); !errors.Is(got, tt.want) {
			t.Errorf("checkRegistrationWindow at %v = %v, want %v", tt.now, got, tt.want)
		}
	}
}
