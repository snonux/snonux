package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
)

// SNONUX_SYNC_USER overrides the SSH username for rsync (default: current login name).
const envSyncUser = "SNONUX_SYNC_USER"

var syncTargets = []string{
	"pi0.lan.buetow.org",
	"pi1.lan.buetow.org",
}

const syncRemoteDir = "/var/www/html/snonux/"

// syncOutput rsyncs localOutput (trailing-slash source) to each sync target over SSH
// port 22. It runs only if every target answers ICMP ping (Linux iputils: ping -c 1 -W …).
func syncOutput(localOutput string) error {
	for _, host := range syncTargets {
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

	absOut, err := filepath.Abs(localOutput)
	if err != nil {
		return fmt.Errorf("sync output dir: %w", err)
	}
	src := filepath.Clean(absOut) + string(filepath.Separator)

	ssh := "ssh -p 22 -o BatchMode=yes -o ConnectTimeout=15"
	for _, host := range syncTargets {
		dest := fmt.Sprintf("%s@%s:%s", sshUser, host, syncRemoteDir)
		log.Printf("rsync %s -> %s", src, dest)
		cmd := exec.Command("rsync", "-az", "-e", ssh, src, dest)
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
