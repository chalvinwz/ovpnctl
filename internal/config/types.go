package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Profile struct {
	Name           string `yaml:"name" mapstructure:"name"`
	ConfigFile     string `yaml:"config_file" mapstructure:"config_file"`
	Username       string `yaml:"username" mapstructure:"username"`
	Password       string `yaml:"password" mapstructure:"password"`
	PrivateKeyPass string `yaml:"private_key_pass" mapstructure:"private_key_pass"`
}

type Config struct {
	Profiles []Profile `yaml:"profiles" mapstructure:"profiles"`
}

func (p *Profile) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(p.ConfigFile) == "" {
		return fmt.Errorf("config_file is required")
	}
	if strings.TrimSpace(p.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(p.Password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func (p *Profile) ExpandedConfigFile() string {
	path := strings.TrimSpace(p.ConfigFile)
	if path == "" {
		return path
	}

	if strings.HasPrefix(path, "~/") || path == "~" {
		home, err := os.UserHomeDir()
		if err == nil && home != "" {
			if path == "~" {
				return home
			}
			return filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}

	return path
}
