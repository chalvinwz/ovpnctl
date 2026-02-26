package openvpn3

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/chalvinwz/ovpnctl/internal/config"
)

type Session struct {
	Path   string
	Config string
}

var (
	execCommand           = exec.Command
	lookPath              = exec.LookPath
	stdout      io.Writer = os.Stdout
	stderr      io.Writer = os.Stderr
)

func RequireBinary() error {
	bin := binaryName()
	if _, err := lookPath(bin); err != nil {
		return fmt.Errorf("%s command not found in PATH (tip: set OVPN3_BIN to full binary path)", bin)
	}
	return nil
}

func StartSession(p *config.Profile, otp string) error {
	cmd := execCommand(binaryName(), "session-start", "--config", p.ExpandedConfigFile())

	input := p.Username + "\n" +
		p.Password + "\n" +
		otp + "\n" +
		p.PrivateKeyPass + "\n"

	cmd.Stdin = bytes.NewBufferString(input)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("openvpn3 session-start failed: %w", err)
	}
	return nil
}

func ListSessions() ([]Session, error) {
	out, err := execCommand(binaryName(), "sessions-list").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("sessions-list failed: %w\n%s", err, out)
	}
	return parseSessions(string(out)), nil
}

func Disconnect(path string) error {
	cmd := execCommand(binaryName(), "session-manage", "--session-path", path, "--disconnect")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("disconnect failed: %w\noutput: %s", err, out)
	}
	return nil
}

func PrintSessions() error {
	sessions, err := ListSessions()
	if err != nil {
		fmt.Fprintf(stderr, "warning: %v\n", err)
		return err
	}
	if len(sessions) == 0 {
		fmt.Println("No active sessions.")
		return nil
	}
	for i, s := range sessions {
		fmt.Printf("  %2d. %s\n", i+1, s.Path)
		if s.Config != "" {
			fmt.Printf("     → %s\n", s.Config)
		}
	}
	return nil
}

func parseSessions(raw string) []Session {
	var sessions []Session
	lines := strings.Split(raw, "\n")
	var current *Session

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Session path:") || strings.HasPrefix(line, "Path:") {
			if current != nil {
				sessions = append(sessions, *current)
			}
			current = &Session{}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				current.Path = strings.TrimSpace(parts[1])
			}
			continue
		}

		if current != nil && (strings.Contains(line, "Configuration") || strings.Contains(line, "config")) {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) > 1 {
				current.Config = strings.TrimSpace(parts[1])
			}
		}
	}

	if current != nil {
		sessions = append(sessions, *current)
	}

	return sessions
}

func LooksLikeSessionPath(s string) bool {
	return strings.HasPrefix(s, "/net/openvpn/v3/sessions/")
}

func binaryName() string {
	if v := strings.TrimSpace(os.Getenv("OVPN3_BIN")); v != "" {
		return v
	}
	return "openvpn3"
}
