package post

import (
	"testing"
	"time"
)

func TestNewID(t *testing.T) {
	t.Parallel()

	loc := time.FixedZone("CET", 1*3600)
	base := time.Date(2026, 4, 9, 14, 30, 22, 0, loc)

	tests := []struct {
		name   string
		tm     time.Time
		suffix int
		want   string
	}{
		{
			name:   "utc no suffix",
			tm:     time.Date(2026, 4, 9, 14, 30, 22, 0, time.UTC),
			suffix: 0,
			want:   "2026-04-09-143022",
		},
		{
			name:   "non utc converts to utc",
			tm:     base,
			suffix: 0,
			want:   "2026-04-09-133022",
		},
		{
			name:   "suffix one",
			tm:     time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
			suffix: 1,
			want:   "2026-01-02-030405-1",
		},
		{
			name:   "suffix large",
			tm:     time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
			suffix: 42,
			want:   "2026-01-02-030405-42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewID(tt.tm, tt.suffix)
			if got != tt.want {
				t.Fatalf("NewID(%v, %d) = %q; want %q", tt.tm, tt.suffix, got, tt.want)
			}
		})
	}
}
