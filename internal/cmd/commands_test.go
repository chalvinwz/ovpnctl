package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/chalvinwz/ovpnctl/internal/config"
	"github.com/chalvinwz/ovpnctl/internal/openvpn3"
)

func TestProfilesCmd(t *testing.T) {
	origLoad := loadConfigCmd
	defer func() { loadConfigCmd = origLoad }()

	loadConfigCmd = func() (*config.Config, error) {
		return &config.Config{Profiles: []config.Profile{{Name: "office", ConfigFile: "~/a", Username: "u", Password: "p"}}}, nil
	}
	if err := profilesCmd().RunE(nil, nil); err != nil {
		t.Fatalf("profiles cmd failed: %v", err)
	}

	loadConfigCmd = func() (*config.Config, error) {
		return &config.Config{Profiles: []config.Profile{{Name: "broken"}}}, nil
	}
	if err := profilesCmd().RunE(nil, nil); err != nil {
		t.Fatalf("profiles cmd invalid profile listing failed: %v", err)
	}

	loadConfigCmd = func() (*config.Config, error) { return &config.Config{}, nil }
	if err := profilesCmd().RunE(nil, nil); err != nil {
		t.Fatalf("profiles cmd empty failed: %v", err)
	}

	loadConfigCmd = func() (*config.Config, error) { return nil, errors.New("x") }
	if err := profilesCmd().RunE(nil, nil); err == nil {
		t.Fatalf("expected error")
	}
}

func TestSessionsCmd(t *testing.T) {
	origReq, origList := requireBinaryCmd, listSessionsCmd
	defer func() { requireBinaryCmd, listSessionsCmd = origReq, origList }()

	requireBinaryCmd = func() error { return nil }
	listSessionsCmd = func() ([]openvpn3.Session, error) { return []openvpn3.Session{{Path: "p1", Config: "c"}}, nil }
	if err := sessionsCmd().RunE(nil, nil); err != nil {
		t.Fatalf("sessions cmd failed: %v", err)
	}

	listSessionsCmd = func() ([]openvpn3.Session, error) { return nil, nil }
	if err := sessionsCmd().RunE(nil, nil); err != nil {
		t.Fatalf("sessions cmd empty failed: %v", err)
	}

	listSessionsCmd = func() ([]openvpn3.Session, error) { return nil, errors.New("list failed") }
	if err := sessionsCmd().RunE(nil, nil); err == nil {
		t.Fatalf("expected list error")
	}

	requireBinaryCmd = func() error { return errors.New("missing") }
	if err := sessionsCmd().RunE(nil, nil); err == nil {
		t.Fatalf("expected binary error")
	}
}

func TestDisconnectCmd(t *testing.T) {
	origReq, origList := requireBinaryCmd, listSessionsCmd
	origDisc, origPrint, origGet := disconnectCmdExec, printSessionsCmd, getProfileCmd
	defer func() {
		requireBinaryCmd, listSessionsCmd = origReq, origList
		disconnectCmdExec, printSessionsCmd, getProfileCmd = origDisc, origPrint, origGet
	}()

	requireBinaryCmd = func() error { return nil }
	listSessionsCmd = func() ([]openvpn3.Session, error) {
		return []openvpn3.Session{{Path: "/net/openvpn/v3/sessions/abc", Config: "/etc/openvpn/office.ovpn"}}, nil
	}
	getProfileCmd = func(string) (*config.Profile, error) {
		return &config.Profile{Name: "office", ConfigFile: "/etc/openvpn/office.ovpn", Username: "u", Password: "p"}, nil
	}
	called := ""
	disconnectCmdExec = func(path string) error { called = path; return nil }
	printSessionsCmd = func() error { return nil }

	if err := disconnectCmd().RunE(nil, []string{"1"}); err != nil {
		t.Fatalf("disconnect cmd by number failed: %v", err)
	}
	if called == "" {
		t.Fatalf("expected disconnect to be called")
	}

	if err := disconnectCmd().RunE(nil, []string{"office"}); err != nil {
		t.Fatalf("disconnect cmd by profile failed: %v", err)
	}

	if err := disconnectCmd().RunE(nil, []string{"/net/openvpn/v3/sessions/direct"}); err != nil {
		t.Fatalf("disconnect cmd by direct path failed: %v", err)
	}

	disconnectCmdExec = func(path string) error { return errors.New("disconnect failed") }
	if err := disconnectCmd().RunE(nil, []string{"1"}); err == nil {
		t.Fatalf("expected disconnect error")
	}

	listSessionsCmd = func() ([]openvpn3.Session, error) { return nil, errors.New("list failed") }
	requireBinaryCmd = func() error { return nil }
	if err := disconnectCmd().RunE(nil, []string{"1"}); err == nil {
		t.Fatalf("expected resolve/list error")
	}

	requireBinaryCmd = func() error { return errors.New("missing") }
	if err := disconnectCmd().RunE(nil, []string{"1"}); err == nil {
		t.Fatalf("expected binary error")
	}
}

func TestConnectCmd(t *testing.T) {
	origReq, origGet, origStart := requireBinaryCmd, getProfileCmd, startSessionCmd
	defer func() { requireBinaryCmd, getProfileCmd, startSessionCmd = origReq, origGet, origStart }()

	requireBinaryCmd = func() error { return nil }
	getProfileCmd = func(string) (*config.Profile, error) {
		return &config.Profile{Name: "office", ConfigFile: "/etc/openvpn/office.ovpn", Username: "u", Password: "p"}, nil
	}
	startSessionCmd = func(*config.Profile, string) error { return nil }

	origStdin := os.Stdin
	defer func() { os.Stdin = origStdin }()
	r, w, _ := os.Pipe()
	_, _ = w.WriteString("123456\n")
	_ = w.Close()
	os.Stdin = r

	if err := connectCmd().RunE(nil, []string{"office"}); err != nil {
		t.Fatalf("connect cmd failed: %v", err)
	}

	requireBinaryCmd = func() error { return errors.New("missing") }
	if err := connectCmd().RunE(nil, []string{"office"}); err == nil {
		t.Fatalf("expected binary error")
	}

	requireBinaryCmd = func() error { return nil }
	getProfileCmd = func(string) (*config.Profile, error) { return nil, errors.New("no profile") }
	if err := connectCmd().RunE(nil, []string{"office"}); err == nil {
		t.Fatalf("expected profile error")
	}

	getProfileCmd = func(string) (*config.Profile, error) {
		return &config.Profile{Name: "office", ConfigFile: "/etc/openvpn/office.ovpn", Username: "u", Password: "p"}, nil
	}
	r2, w2, _ := os.Pipe()
	_ = w2.Close()
	os.Stdin = r2
	if err := connectCmd().RunE(nil, []string{"office"}); err == nil {
		t.Fatalf("expected otp read error")
	}

	r3, w3, _ := os.Pipe()
	_, _ = w3.WriteString("123456\n")
	_ = w3.Close()
	os.Stdin = r3
	startSessionCmd = func(*config.Profile, string) error { return errors.New("start failed") }
	if err := connectCmd().RunE(nil, []string{"office"}); err == nil {
		t.Fatalf("expected start error")
	}
}

func TestResolveDisconnectPath(t *testing.T) {
	origList, origGet := listSessionsCmd, getProfileCmd
	defer func() { listSessionsCmd, getProfileCmd = origList, origGet }()

	listSessionsCmd = func() ([]openvpn3.Session, error) {
		return []openvpn3.Session{{Path: "/net/openvpn/v3/sessions/abc", Config: "/etc/openvpn/office.ovpn"}}, nil
	}
	getProfileCmd = func(string) (*config.Profile, error) {
		return &config.Profile{Name: "office", ConfigFile: "/etc/openvpn/office.ovpn", Username: "u", Password: "p"}, nil
	}

	if p, err := resolveDisconnectPath("/net/openvpn/v3/sessions/direct"); err != nil || p == "" {
		t.Fatalf("expected direct path success, got %q err=%v", p, err)
	}
	if p, err := resolveDisconnectPath("1"); err != nil || p != "/net/openvpn/v3/sessions/abc" {
		t.Fatalf("expected numeric success, got %q err=%v", p, err)
	}
	if p, err := resolveDisconnectPath("office"); err != nil || p != "/net/openvpn/v3/sessions/abc" {
		t.Fatalf("expected profile success, got %q err=%v", p, err)
	}

	listSessionsCmd = func() ([]openvpn3.Session, error) {
		return []openvpn3.Session{{Path: "/net/openvpn/v3/sessions/def", Config: "office.ovpn"}}, nil
	}
	if p, err := resolveDisconnectPath("office"); err != nil || p != "/net/openvpn/v3/sessions/def" {
		t.Fatalf("expected profile basename success, got %q err=%v", p, err)
	}

	listSessionsCmd = func() ([]openvpn3.Session, error) {
		return []openvpn3.Session{{Path: "/net/openvpn/v3/sessions/ghi", Config: "config=/etc/openvpn/office.ovpn"}}, nil
	}
	if p, err := resolveDisconnectPath("office"); err != nil || p != "/net/openvpn/v3/sessions/ghi" {
		t.Fatalf("expected profile contains success, got %q err=%v", p, err)
	}

	listSessionsCmd = func() ([]openvpn3.Session, error) {
		return []openvpn3.Session{
			{Path: "/net/openvpn/v3/sessions/empty", Config: ""},
			{Path: "/net/openvpn/v3/sessions/jkl", Config: "/etc/openvpn/office.ovpn"},
		}, nil
	}
	if p, err := resolveDisconnectPath("office"); err != nil || p != "/net/openvpn/v3/sessions/jkl" {
		t.Fatalf("expected profile match skipping empty config, got %q err=%v", p, err)
	}

	if _, err := resolveDisconnectPath("9"); err == nil {
		t.Fatalf("expected out-of-range error")
	}

	getProfileCmd = func(string) (*config.Profile, error) {
		return nil, errors.New("missing")
	}
	if _, err := resolveDisconnectPath("unknown"); err == nil {
		t.Fatalf("expected invalid argument error")
	}

	getProfileCmd = func(string) (*config.Profile, error) {
		return &config.Profile{Name: "office", ConfigFile: "/etc/openvpn/other.ovpn", Username: "u", Password: "p"}, nil
	}
	if _, err := resolveDisconnectPath("office"); err == nil {
		t.Fatalf("expected no active session for profile error")
	}

	listSessionsCmd = func() ([]openvpn3.Session, error) { return nil, errors.New("list") }
	if _, err := resolveDisconnectPath("1"); err == nil {
		t.Fatalf("expected list error")
	}
}
