package config

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config: %w", err)
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

	// numeric index
	if n, err := strconv.Atoi(id); err == nil {
		if n >= 1 && n <= len(cfg.Profiles) {
			return &cfg.Profiles[n-1], nil
		}
	}

	// name match
	for _, p := range cfg.Profiles {
		if p.Name == id {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("profile not found: %q", id)
}
