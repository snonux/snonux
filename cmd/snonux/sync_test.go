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
	unsetEnv(t, "SNONUX_SYNC_TARGETS")
	unsetEnv(t, "SNONUX_SYNC_REMOTE_DIR")

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
	t.Setenv("SNONUX_SYNC_TARGETS", "h1, h2 ,h3")
	unsetEnv(t, "SNONUX_SYNC_REMOTE_DIR")
	cfg := &config.Config{}
	resolveSyncConfig(cfg)
	want := []string{"h1", "h2", "h3"}
	if !reflect.DeepEqual(cfg.SyncTargets, want) {
		t.Fatalf("got targets %v, want %v", cfg.SyncTargets, want)
	}
}

func TestResolveSyncConfig_envRemoteDir(t *testing.T) {
	unsetEnv(t, "SNONUX_SYNC_TARGETS")
	t.Setenv("SNONUX_SYNC_REMOTE_DIR", "/custom/path/")
	cfg := &config.Config{}
	resolveSyncConfig(cfg)
	if cfg.SyncRemoteDir != "/custom/path/" {
		t.Fatalf("got remote dir %q", cfg.SyncRemoteDir)
	}
}

func TestResolveSyncConfig_flagsOverrideEnv(t *testing.T) {
	t.Setenv("SNONUX_SYNC_TARGETS", "from-env")
	t.Setenv("SNONUX_SYNC_REMOTE_DIR", "/env/dir/")

	cfg := &config.Config{
		SyncTargets:   []string{"from-flag"},
		SyncRemoteDir: "/flag/dir/",
	}
	resolveSyncConfig(cfg)

	if !reflect.DeepEqual(cfg.SyncTargets, []string{"from-flag"}) {
		t.Fatalf("flag targets should override env: got %v", cfg.SyncTargets)
	}
	if cfg.SyncRemoteDir != "/flag/dir/" {
		t.Fatalf("flag dir should override env: got %q", cfg.SyncRemoteDir)
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	orig, hadOrig := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}
	t.Cleanup(func() {
		if hadOrig {
			_ = os.Setenv(key, orig)
		} else {
			_ = os.Unsetenv(key)
		}
	})
}
