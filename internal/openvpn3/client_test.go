package openvpn3

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/chalvinwz/ovpnctl/internal/config"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	if os.Getenv("OVPN_HELPER_ERR") == "1" {
		fmt.Fprint(os.Stdout, os.Getenv("OVPN_HELPER_OUT"))
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, os.Getenv("OVPN_HELPER_OUT"))
	os.Exit(0)
}

func helperCmd(out string, fail bool) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cmdArgs := []string{"-test.run=TestHelperProcess", "--", name}
		cmdArgs = append(cmdArgs, args...)
		cmd := exec.Command(os.Args[0], cmdArgs...)
		cmd.Env = append(os.Environ(),
			"GO_WANT_HELPER_PROCESS=1",
			"OVPN_HELPER_OUT="+out,
		)
		if fail {
			cmd.Env = append(cmd.Env, "OVPN_HELPER_ERR=1")
		}
		return cmd
	}
}

func TestParseSessions(t *testing.T) {
	raw := `
Session path: /net/openvpn/v3/sessions/abc
Configuration name: office

Session path: /net/openvpn/v3/sessions/def
Configuration name: home
`

	sessions := parseSessions(raw)
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
	if sessions[0].Path != "/net/openvpn/v3/sessions/abc" || sessions[0].Config != "office" {
		t.Fatalf("unexpected session[0]: %+v", sessions[0])
	}
	if sessions[1].Path != "/net/openvpn/v3/sessions/def" || sessions[1].Config != "home" {
		t.Fatalf("unexpected session[1]: %+v", sessions[1])
	}
}

func TestParseSessions_PathLabelVariant(t *testing.T) {
	raw := `
Path: /net/openvpn/v3/sessions/xyz
config: corp
`

	sessions := parseSessions(raw)
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Path != "/net/openvpn/v3/sessions/xyz" || sessions[0].Config != "corp" {
		t.Fatalf("unexpected session: %+v", sessions[0])
	}
}

func TestLooksLikeSessionPath(t *testing.T) {
	if !LooksLikeSessionPath("/net/openvpn/v3/sessions/abc") {
		t.Fatalf("expected valid session path")
	}
	if LooksLikeSessionPath("abc") {
		t.Fatalf("expected invalid session path")
	}
}

func TestRequireBinary(t *testing.T) {
	orig := lookPath
	defer func() { lookPath = orig }()

	lookPath = func(file string) (string, error) { return "", errors.New("nope") }
	if err := RequireBinary(); err == nil {
		t.Fatalf("expected error when binary missing")
	}

	lookPath = func(file string) (string, error) { return "/usr/bin/openvpn3", nil }
	if err := RequireBinary(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStartSession(t *testing.T) {
	origExec := execCommand
	defer func() { execCommand = origExec }()

	execCommand = helperCmd("", false)
	p := &config.Profile{Name: "n", ConfigFile: "~/x.ovpn", Username: "u", Password: "p"}
	if err := StartSession(p, "123456"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	execCommand = helperCmd("boom", true)
	if err := StartSession(p, "123456"); err == nil {
		t.Fatalf("expected start error")
	}
}

func TestListSessions(t *testing.T) {
	origExec := execCommand
	defer func() { execCommand = origExec }()

	execCommand = helperCmd("Session path: /net/openvpn/v3/sessions/abc\nConfiguration name: office\n", false)
	s, err := ListSessions()
	if err != nil || len(s) != 1 {
		t.Fatalf("expected one session, got %v err=%v", s, err)
	}

	execCommand = helperCmd("bad", true)
	if _, err := ListSessions(); err == nil {
		t.Fatalf("expected list error")
	}
}

func TestDisconnect(t *testing.T) {
	origExec := execCommand
	defer func() { execCommand = origExec }()

	execCommand = helperCmd("", false)
	if err := Disconnect("/net/openvpn/v3/sessions/abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	execCommand = helperCmd("oops", true)
	if err := Disconnect("/net/openvpn/v3/sessions/abc"); err == nil {
		t.Fatalf("expected disconnect error")
	}
}

func TestPrintSessions(t *testing.T) {
	origExec, origOut, origErr := execCommand, stdout, stderr
	defer func() { execCommand, stdout, stderr = origExec, origOut, origErr }()

	var outBuf, errBuf bytes.Buffer
	stdout, stderr = &outBuf, &errBuf

	execCommand = helperCmd("", false)
	if err := PrintSessions(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	execCommand = helperCmd("Session path: /net/openvpn/v3/sessions/abc\nConfiguration name: office\n", false)
	if err := PrintSessions(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	execCommand = helperCmd("err", true)
	if err := PrintSessions(); err == nil {
		t.Fatalf("expected print error")
	}
}
