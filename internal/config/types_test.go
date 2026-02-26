package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfileValidate(t *testing.T) {
	valid := Profile{
		Name:       "office",
		ConfigFile: "~/vpn/office.ovpn",
		Username:   "user",
		Password:   "pass",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid profile, got error: %v", err)
	}

	tests := []struct {
		name    string
		profile Profile
	}{
		{"missing name", Profile{ConfigFile: "x.ovpn", Username: "u", Password: "p"}},
		{"missing config", Profile{Name: "x", Username: "u", Password: "p"}},
		{"missing username", Profile{Name: "x", ConfigFile: "x.ovpn", Password: "p"}},
		{"missing password", Profile{Name: "x", ConfigFile: "x.ovpn", Username: "u"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.profile.Validate(); err == nil {
				t.Fatalf("expected validation error")
			}
		})
	}
}

func TestExpandedConfigFile(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		t.Skip("home dir unavailable")
	}

	p0 := Profile{ConfigFile: ""}
	if got := p0.ExpandedConfigFile(); got != "" {
		t.Fatalf("expected empty path, got %q", got)
	}

	p := Profile{ConfigFile: "~/vpn/office.ovpn"}
	got := p.ExpandedConfigFile()
	want := filepath.Join(home, "vpn", "office.ovpn")
	if got != want {
		t.Fatalf("ExpandedConfigFile() = %q, want %q", got, want)
	}

	p2 := Profile{ConfigFile: "/etc/openvpn/office.ovpn"}
	if got := p2.ExpandedConfigFile(); got != "/etc/openvpn/office.ovpn" {
		t.Fatalf("expected unchanged absolute path, got %q", got)
	}

	p3 := Profile{ConfigFile: "~"}
	if got := p3.ExpandedConfigFile(); got != home {
		t.Fatalf("expected home path, got %q", got)
	}
}
