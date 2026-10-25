package dbx

import "testing"

func TestContainsPattern(t *testing.T) {
	tests := []struct {
		query, want string
	}{
		{"", "%%"},
		{"  CSC  ", "%CSC%"},
		{"100%", `%100\%%`},
		{"a_b", `%a\_b%`},
		{`C:\path`, `%C:\\path%`},
		{`\%`, `%\\\%%`},
	}
	for _, tt := range tests {
		if got := ContainsPattern(tt.query); got != tt.want {
			t.Errorf("ContainsPattern(%q) = %q, want %q", tt.query, got, tt.want)
		}
	}
}
