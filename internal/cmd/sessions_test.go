package cmd

import (
	"errors"
	"testing"

	"github.com/chalvinwz/ovpnctl/internal/openvpn3"
)

func TestSessionsCmd(t *testing.T) {
	origReq, origList := requireBinaryCmd, listSessionsCmd
	defer func() { requireBinaryCmd, listSessionsCmd = origReq, origList }()

	requireBinaryCmd = func() error { return nil }
	listSessionsCmd = func() ([]openvpn3.Session, error) {
		return []openvpn3.Session{{Path: "p1", Config: "/etc/openvpn/office.ovpn"}}, nil
	}
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
