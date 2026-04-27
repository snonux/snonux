package main

import (
	"os"
	"reflect"
	"testing"

	"codeberg.org/snonux/snonux/internal/config"
)

func TestSplitAndTrim(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{" a , b ", []string{"a", "b"}},
		{"a,,c", []string{"a", "c"}},
		{"", []string{}},
		{",,", []string{}},
	}

	for _, c := range cases {
		got := splitAndTrim(c.in)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("splitAndTrim(%q) = %v; want %v", c.in, got, c.want)
		}
	}
}

func TestResolveSyncConfig_defaults(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{}
	resolveSyncConfig(cfg)
	want := []string{"pi0.lan.buetow.org", "pi1.lan.buetow.org"}
	if !reflect.DeepEqual(cfg.SyncTargets, want) {
		t.Fatalf("got targets %v, want %v", cfg.SyncTargets, want)
	}
	if cfg.SyncRemoteDir != "/var/www/html/snonux/" {
		t.Fatalf("got remote dir %q", cfg.SyncRemoteDir)
	}
}

func TestResolveSyncConfig_envTargets(t *testing.T) {
	t.Parallel()

	orig := os.Getenv("SNONUX_SYNC_TARGETS")
	defer os.Setenv("SNONUX_SYNC_TARGETS", orig)

	os.Setenv("SNONUX_SYNC_TARGETS", "h1, h2 ,h3")
	cfg := &config.Config{}
	resolveSyncConfig(cfg)
	want := []string{"h1", "h2", "h3"}
	if !reflect.DeepEqual(cfg.SyncTargets, want) {
		t.Fatalf("got targets %v, want %v", cfg.SyncTargets, want)
	}
}

func TestResolveSyncConfig_envRemoteDir(t *testing.T) {
	t.Parallel()

	orig := os.Getenv("SNONUX_SYNC_REMOTE_DIR")
	defer os.Setenv("SNONUX_SYNC_REMOTE_DIR", orig)

	os.Setenv("SNONUX_SYNC_REMOTE_DIR", "/custom/path/")
	cfg := &config.Config{}
	resolveSyncConfig(cfg)
	if cfg.SyncRemoteDir != "/custom/path/" {
		t.Fatalf("got remote dir %q", cfg.SyncRemoteDir)
	}
}

func TestResolveSyncConfig_flagsOverrideEnv(t *testing.T) {
	t.Parallel()

	origTargets := os.Getenv("SNONUX_SYNC_TARGETS")
	origDir := os.Getenv("SNONUX_SYNC_REMOTE_DIR")
	defer func() {
		os.Setenv("SNONUX_SYNC_TARGETS", origTargets)
		os.Setenv("SNONUX_SYNC_REMOTE_DIR", origDir)
	}()

	os.Setenv("SNONUX_SYNC_TARGETS", "from-env")
	os.Setenv("SNONUX_SYNC_REMOTE_DIR", "/env/dir/")

	cfg := &config.Config{
		SyncTargets:     []string{"from-flag"},
		SyncRemoteDir:   "/flag/dir/",
	}
	resolveSyncConfig(cfg)

	if !reflect.DeepEqual(cfg.SyncTargets, []string{"from-flag"}) {
		t.Fatalf("flag targets should override env: got %v", cfg.SyncTargets)
	}
	if cfg.SyncRemoteDir != "/flag/dir/" {
		t.Fatalf("flag dir should override env: got %q", cfg.SyncRemoteDir)
	}
}
