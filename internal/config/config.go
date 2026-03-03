package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config: %w", err)
	}

	seen := make(map[string]struct{}, len(cfg.Profiles))
	for i := range cfg.Profiles {
		cfg.Profiles[i].Name = strings.TrimSpace(cfg.Profiles[i].Name)
		cfg.Profiles[i].ConfigFile = strings.TrimSpace(cfg.Profiles[i].ConfigFile)
		cfg.Profiles[i].Username = strings.TrimSpace(cfg.Profiles[i].Username)
		cfg.Profiles[i].Password = strings.TrimSpace(cfg.Profiles[i].Password)
		cfg.Profiles[i].PrivateKeyPass = strings.TrimSpace(cfg.Profiles[i].PrivateKeyPass)

		key := strings.ToLower(cfg.Profiles[i].Name)
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			return nil, fmt.Errorf("duplicate profile name: %q", cfg.Profiles[i].Name)
		}
		seen[key] = struct{}{}
	}

	return &cfg, nil
}

func GetProfile(id string) (*Profile, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	if len(cfg.Profiles) == 0 {
		return nil, fmt.Errorf("no profiles defined in configuration")
	}

	id = strings.TrimSpace(id)

	// numeric index
	if n, err := strconv.Atoi(id); err == nil {
		if n >= 1 && n <= len(cfg.Profiles) {
			p := cfg.Profiles[n-1]
			if err := p.Validate(); err != nil {
				return nil, fmt.Errorf("invalid profile %q: %w", p.Name, err)
			}
			return &p, nil
		}
	}

	// case-insensitive name match
	for _, p := range cfg.Profiles {
		if strings.EqualFold(strings.TrimSpace(p.Name), id) {
			if err := p.Validate(); err != nil {
				return nil, fmt.Errorf("invalid profile %q: %w", p.Name, err)
			}
			return &p, nil
		}
	}

	return nil, fmt.Errorf("profile not found: %q", id)
}
