package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/snonux/snonux/internal/config"
)

// SNONUX_SYNC_USER overrides the SSH username for rsync (default: current login name).
const envSyncUser = "SNONUX_SYNC_USER"

// defaultSyncTargets are the built-in mirror hosts used when no configuration overrides them.
var defaultSyncTargets = []string{
	"pi0.lan.buetow.org",
	"pi1.lan.buetow.org",
}

const defaultSyncRemoteDir = "/var/www/html/snonux/"

// resolveSyncConfig populates cfg.SyncTargets and cfg.SyncRemoteDir from the
// environment if they are empty, applying sensible defaults.
func resolveSyncConfig(cfg *config.Config) {
	if len(cfg.SyncTargets) == 0 {
		if v := os.Getenv("SNONUX_SYNC_TARGETS"); v != "" {
			cfg.SyncTargets = splitAndTrim(v)
		} else {
			cfg.SyncTargets = append([]string(nil), defaultSyncTargets...)
		}
	}
	if cfg.SyncRemoteDir == "" {
		if v := os.Getenv("SNONUX_SYNC_REMOTE_DIR"); v != "" {
			cfg.SyncRemoteDir = strings.TrimSpace(v)
		} else {
			cfg.SyncRemoteDir = defaultSyncRemoteDir
		}
	}
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// syncOutput rsyncs localOutput (trailing-slash source) to each sync target over SSH
// port 22. It runs only if every target answers ICMP ping (Linux iputils: ping -c 1 -W …).
// The ctx parameter is accepted for cancellation propagation; it is wired into
// exec.CommandContext for the rsync subprocesses.
func syncOutput(ctx context.Context, cfg *config.Config) error {
	resolveSyncConfig(cfg)

	for _, host := range cfg.SyncTargets {
		if !hostPingable(host) {
			log.Printf("sync skipped: %q not pingable (all mirror hosts must be reachable)", host)
			return nil
		}
	}

	sshUser := os.Getenv(envSyncUser)
	if sshUser == "" {
		u, err := user.Current()
		if err != nil {
			return fmt.Errorf("sync user: %w (set %s)", err, envSyncUser)
		}
		sshUser = u.Username
	}

	absOut, err := filepath.Abs(cfg.OutputDir)
	if err != nil {
		return fmt.Errorf("sync output dir: %w", err)
	}
	src := filepath.Clean(absOut) + string(filepath.Separator)

	ssh := "ssh -p 22 -o BatchMode=yes -o ConnectTimeout=15"
	for _, host := range cfg.SyncTargets {
		dest := fmt.Sprintf("%s@%s:%s", sshUser, host, cfg.SyncRemoteDir)
		log.Printf("rsync %s -> %s", src, dest)
		cmd := exec.CommandContext(ctx, "rsync", "-az", "-e", ssh, src, dest)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("rsync to %s: %w", host, err)
		}
	}
	return nil
}

func hostPingable(host string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Linux iputils-ping: -c 1 one packet, -W 3 wait up to 3s for reply.
	cmd := exec.CommandContext(ctx, "ping", "-c", "1", "-W", "3", host)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}
