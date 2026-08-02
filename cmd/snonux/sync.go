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

const wireGuardJumpHost = "rex@fishfinger.buetow.org:2"

// defaultSyncTargets are the built-in mirror hosts used when no configuration overrides them.
var defaultSyncTargets = []string{
	"pi0.lan.buetow.org",
	"pi1.lan.buetow.org",
}

const defaultSyncRemoteDir = "/var/www/html/snonux.foo/"

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
// port 22. It runs only if every target or its WireGuard fallback answers ICMP
// ping (Linux iputils: ping -c 1 -W …).
// The ctx parameter is accepted for cancellation propagation; it is wired into
// exec.CommandContext for the rsync subprocesses.
func syncOutput(ctx context.Context, cfg *config.Config) error {
	resolveSyncConfig(cfg)

	sshUser := os.Getenv(envSyncUser)
	if sshUser == "" {
		u, err := user.Current()
		if err != nil {
			return fmt.Errorf("sync user: %w (set %s)", err, envSyncUser)
		}
		sshUser = u.Username
	}

	syncTargets := make([]string, 0, len(cfg.SyncTargets))
	for _, host := range cfg.SyncTargets {
		target, ok := reachableSyncTarget(host, hostPingable, func(fallback string) bool {
			return hostReachableViaWireGuard(fallback, sshUser)
		})
		if !ok {
			fallback := wireGuardFallback(host)
			if fallback != "" {
				log.Printf("sync skipped: neither %q nor WireGuard fallback %q via %q is reachable (all mirror hosts must be reachable)", host, fallback, wireGuardJumpHost)
			} else {
				log.Printf("sync skipped: %q not pingable (all mirror hosts must be reachable)", host)
			}
			return nil
		}
		if target != host {
			log.Printf("sync target %q not pingable; using WireGuard fallback %q via %q", host, target, wireGuardJumpHost)
		}
		syncTargets = append(syncTargets, target)
	}

	absOut, err := filepath.Abs(cfg.OutputDir)
	if err != nil {
		return fmt.Errorf("sync output dir: %w", err)
	}
	src := filepath.Clean(absOut) + string(filepath.Separator)

	for _, host := range syncTargets {
		dest := fmt.Sprintf("%s@%s:%s", sshUser, host, cfg.SyncRemoteDir)
		log.Printf("rsync %s -> %s", src, dest)
		// --chmod overrides the locally-generated (mode 600) output permissions:
		// the remote webserver runs as its own unprivileged user (e.g. bozohttpd's
		// _httpd), not as the SSH login user, so published files must be
		// world-readable regardless of local perms.
		ssh := "ssh -p 22 -o BatchMode=yes -o ConnectTimeout=15"
		if wireGuardFallbackHost(host) {
			ssh += " -o HostKeyAlias=" + wireGuardHostKeyAlias(host) + " -J " + wireGuardJumpHost
		}
		cmd := exec.CommandContext(ctx, "rsync", "-az", "--chmod=D755,F644", "-e", ssh, src, dest)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("rsync to %s: %w", host, err)
		}
	}
	return nil
}

func reachableSyncTarget(host string, pingable, fallbackReachable func(string) bool) (string, bool) {
	if pingable(host) {
		return host, true
	}

	fallback := wireGuardFallback(host)
	if fallback != "" && fallbackReachable(fallback) {
		return fallback, true
	}

	return "", false
}

func wireGuardFallback(host string) string {
	switch host {
	case "pi0.lan.buetow.org":
		return "pi0.wg0"
	case "pi1.lan.buetow.org":
		return "pi1.wg0"
	default:
		return ""
	}
}

func wireGuardFallbackHost(host string) bool {
	return host == "pi0.wg0" || host == "pi1.wg0"
}

func wireGuardHostKeyAlias(host string) string {
	switch host {
	case "pi0.wg0":
		return "pi0.lan.buetow.org"
	case "pi1.wg0":
		return "pi1.lan.buetow.org"
	default:
		return host
	}
}

func hostReachableViaWireGuard(host, sshUser string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx,
		"ssh", "-p", "22", "-o", "BatchMode=yes", "-o", "ConnectTimeout=5",
		"-o", "HostKeyAlias="+wireGuardHostKeyAlias(host),
		"-J", wireGuardJumpHost, sshUser+"@"+host, "true",
	)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
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
