package processor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindLocalImages(t *testing.T) {
	t.Parallel()

	t.Run("remote skipped", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		got := findLocalImages(`![](https://cdn.example/p.png) ![](http://x/y.jpg)`, dir)
		if len(got) != 0 {
			t.Fatalf("expected no locals, got %v", got)
		}
	})

	t.Run("missing file ignored", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		got := findLocalImages(`![](nope.png)`, dir)
		if len(got) != 0 {
			t.Fatalf("expected no locals, got %v", got)
		}
	})

	t.Run("picks existing basename", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "shot.png"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		got := findLocalImages(`![alt](shot.png)`, dir)
		if len(got) != 1 || got[0] != "shot.png" {
			t.Fatalf("got %v; want [shot.png]", got)
		}
	})

	tests := []struct {
		name    string
		md      string
		files   []string
		want    []string
		wantLen int
	}{
		{
			name:    "multiple locals order",
			md:      `![](a.png) ![](b.png)`,
			files:   []string{"a.png", "b.png"},
			wantLen: 2,
		},
		{
			name:    "alt with spaces",
			md:      `![my photo](z.gif)`,
			files:   []string{"z.gif"},
			want:    []string{"z.gif"},
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			for _, f := range tt.files {
				if err := os.WriteFile(filepath.Join(dir, f), []byte("x"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got := findLocalImages(tt.md, dir)
			if tt.want != nil {
				if len(got) != len(tt.want) {
					t.Fatalf("got %v; want %v", got, tt.want)
				}
				for i := range tt.want {
					if got[i] != tt.want[i] {
						t.Fatalf("got %v; want %v", got, tt.want)
					}
				}
				return
			}
			if len(got) != tt.wantLen {
				t.Fatalf("len(got)=%d; want %d (%v)", len(got), tt.wantLen, got)
			}
		})
	}
}
