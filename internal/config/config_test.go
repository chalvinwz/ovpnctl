package config

import (
	"testing"

	"github.com/spf13/viper"
)

func withProfiles(t *testing.T, profiles []map[string]string) {
	t.Helper()
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("profiles", profiles)
}

func TestGetProfileByNumberAndName(t *testing.T) {
	withProfiles(t, []map[string]string{
		{"name": "Office", "config_file": "~/vpn/office.ovpn", "username": "u", "password": "p"},
		{"name": "Home", "config_file": "~/vpn/home.ovpn", "username": "u2", "password": "p2"},
	})

	p1, err := GetProfile("1")
	if err != nil {
		t.Fatalf("GetProfile by number failed: %v", err)
	}
	if p1.Name != "Office" {
		t.Fatalf("unexpected profile: %+v", p1)
	}

	p2, err := GetProfile(" home ")
	if err != nil {
		t.Fatalf("GetProfile by name failed: %v", err)
	}
	if p2.Name != "Home" {
		t.Fatalf("unexpected profile: %+v", p2)
	}
}

func TestGetProfileInvalidProfile(t *testing.T) {
	withProfiles(t, []map[string]string{
		{"name": "Broken", "config_file": "", "username": "u", "password": "p"},
	})

	if _, err := GetProfile("1"); err == nil {
		t.Fatalf("expected error for invalid profile")
	}
	if _, err := GetProfile("broken"); err == nil {
		t.Fatalf("expected error for invalid profile by name")
	}
}

func TestGetProfileNotFoundAndOutOfRange(t *testing.T) {
	withProfiles(t, []map[string]string{
		{"name": "Office", "config_file": "~/vpn/office.ovpn", "username": "u", "password": "p"},
	})

	if _, err := GetProfile("99"); err == nil {
		t.Fatalf("expected out-of-range error")
	}
	if _, err := GetProfile("does-not-exist"); err == nil {
		t.Fatalf("expected not-found error")
	}
}

func TestGetProfileNoProfiles(t *testing.T) {
	withProfiles(t, nil)
	if _, err := GetProfile("1"); err == nil {
		t.Fatalf("expected no profiles error")
	}
}

func TestLoadDuplicateProfileNames(t *testing.T) {
	withProfiles(t, []map[string]string{
		{"name": "Office", "config_file": "~/vpn/office.ovpn", "username": "u", "password": "p"},
		{"name": "office", "config_file": "~/vpn/office2.ovpn", "username": "u2", "password": "p2"},
	})

	if _, err := Load(); err == nil {
		t.Fatalf("expected duplicate profile name error")
	}
}

func TestLoadAllowsEmptyNameEntry(t *testing.T) {
	withProfiles(t, []map[string]string{
		{"name": "", "config_file": "~/vpn/x.ovpn", "username": "u", "password": "p"},
	})
	if _, err := Load(); err != nil {
		t.Fatalf("expected load success for empty name entry, got: %v", err)
	}
}

func TestLoadUnmarshalError(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("profiles", "not-a-list")

	if _, err := Load(); err == nil {
		t.Fatalf("expected unmarshal error")
	}
}

func TestGetProfileLoadError(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("profiles", "not-a-list")

	if _, err := GetProfile("1"); err == nil {
		t.Fatalf("expected load error")
	}
}
