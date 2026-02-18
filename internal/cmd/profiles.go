package cmd

import (
	"fmt"

	"github.com/chalvinwz/ovpnctl/internal/config"
	"github.com/spf13/cobra"
)

func profilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "profiles",
		Aliases: []string{"ls", "profiles"},
		Short:   "List configured OpenVPN profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if len(cfg.Profiles) == 0 {
				fmt.Println("No profiles configured.")
				return nil
			}

			fmt.Println("Configured profiles:")
			for i, p := range cfg.Profiles {
				fmt.Printf("  %2d. %-20s  %s\n", i+1, p.Name, p.ConfigFile)
			}
			return nil
		},
	}
}
