package cmd

import (
	"errors"
	"testing"

	"github.com/chalvinwz/ovpnctl/internal/config"
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
