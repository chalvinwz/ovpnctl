package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfigWithExplicitPath(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	d := t.TempDir()
	p := filepath.Join(d, "profiles.yaml")
	if err := os.WriteFile(p, []byte("profiles: []\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cfgFile = p
	if err := initConfig(); err != nil {
		t.Fatalf("initConfig failed: %v", err)
	}
}

func TestInitConfigMissingDefaultAllowed(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	cfgFile = ""

	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	_ = os.Chdir(t.TempDir())

	if err := initConfig(); err != nil {
		t.Fatalf("expected missing default config to be allowed, got: %v", err)
	}
}

func TestInitConfigInvalidYaml(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	d := t.TempDir()
	p := filepath.Join(d, "profiles.yaml")
	if err := os.WriteFile(p, []byte("profiles: ["), 0o600); err != nil {
		t.Fatal(err)
	}

	cfgFile = p
	if err := initConfig(); err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestExecuteHelp(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"ovpnctl", "--help"}
	cfgFile = ""
	if err := Execute(); err != nil {
		t.Fatalf("execute help failed: %v", err)
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"ovpnctl", "unknown-command"}
	cfgFile = ""
	if err := Execute(); err == nil {
		t.Fatalf("expected execute error")
	}
}

func TestExecuteProfilesCommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	d := t.TempDir()
	p := filepath.Join(d, "profiles.yaml")
	if err := os.WriteFile(p, []byte("profiles: []\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"ovpnctl", "--config", p, "profiles"}
	cfgFile = ""
	if err := Execute(); err != nil {
		t.Fatalf("execute profiles failed: %v", err)
	}
}
