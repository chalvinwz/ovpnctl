package cmd

import (
	"errors"
	"testing"

	"github.com/chalvinwz/ovpnctl/internal/openvpn3"
)

func TestDisconnectCmd(t *testing.T) {
	origReq, origList := requireBinaryCmd, listSessionsCmd
	origDisc, origPrint := disconnectCmdExec, printSessionsCmd
	defer func() {
		requireBinaryCmd, listSessionsCmd = origReq, origList
		disconnectCmdExec, printSessionsCmd = origDisc, origPrint
	}()

	requireBinaryCmd = func() error { return nil }
	listSessionsCmd = func() ([]openvpn3.Session, error) {
		return []openvpn3.Session{{Path: "/net/openvpn/v3/sessions/abc", Config: "/etc/openvpn/office.ovpn"}}, nil
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

func TestResolveDisconnectPath(t *testing.T) {
	origList := listSessionsCmd
	defer func() { listSessionsCmd = origList }()

	listSessionsCmd = func() ([]openvpn3.Session, error) {
		return []openvpn3.Session{{Path: "/net/openvpn/v3/sessions/abc", Config: "/etc/openvpn/office.ovpn"}}, nil
	}

	if p, err := resolveDisconnectPath("/net/openvpn/v3/sessions/direct"); err != nil || p == "" {
		t.Fatalf("expected direct path success, got %q err=%v", p, err)
	}
	if p, err := resolveDisconnectPath("1"); err != nil || p != "/net/openvpn/v3/sessions/abc" {
		t.Fatalf("expected numeric success, got %q err=%v", p, err)
	}

	if _, err := resolveDisconnectPath("bad"); err == nil {
		t.Fatalf("expected invalid argument error")
	}

	if _, err := resolveDisconnectPath("9"); err == nil {
		t.Fatalf("expected out-of-range error")
	}

	listSessionsCmd = func() ([]openvpn3.Session, error) { return nil, errors.New("list") }
	if _, err := resolveDisconnectPath("1"); err == nil {
		t.Fatalf("expected list error")
	}
}
