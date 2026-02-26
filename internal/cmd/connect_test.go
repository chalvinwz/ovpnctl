package cmd

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/chalvinwz/ovpnctl/internal/config"
)

func TestReadOTP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid otp", input: "123456\n", want: "123456"},
		{name: "trim spaces", input: "  654321  \n", want: "654321"},
		{name: "empty otp", input: "   \n", wantErr: true},
		{name: "no input", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readOTP(strings.NewReader(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("readOTP() = %q, want %q", got, tt.want)
			}
		})
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

	r4, w4, _ := os.Pipe()
	_, _ = w4.WriteString("\n")
	_ = w4.Close()
	os.Stdin = r4
	startSessionCmd = func(*config.Profile, string) error { return nil }
	if err := connectCmd().RunE(nil, []string{"office"}); err == nil {
		t.Fatalf("expected empty otp error")
	}
}
