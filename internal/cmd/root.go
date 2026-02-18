package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "ovpnctl",
		Short: "OpenVPN 3 profile & session manager",
		Long: `ovpnctl manages OpenVPN 3 connections using external YAML profiles.
Supports listing profiles, connecting (OTP prompt only), listing sessions and disconnecting.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to profiles.yaml (default: ~/.config/ovpnctl/profiles.yaml or ./profiles.yaml)")

	rootCmd.AddCommand(
		profilesCmd(),
		connectCmd(),
		sessionsCmd(),
		disconnectCmd(),
	)

	return rootCmd.Execute()
}

func initConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Try standard locations
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(filepath.Join(home, ".config", "ovpnctl"))
		}
		viper.AddConfigPath(".")
		viper.SetConfigName("profiles")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("cannot read config: %w", err)
		}
		// missing file is allowed — subcommands will fail gracefully
	}
	return nil
}
