package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func profilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "profiles",
		Aliases: []string{"ls"},
		Short:   "List configured OpenVPN profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfigCmd()
			if err != nil {
				return err
			}

			if len(cfg.Profiles) == 0 {
				fmt.Println("No profiles configured.")
				return nil
			}

			fmt.Println("Configured profiles:")
			for i, p := range cfg.Profiles {
				if err := p.Validate(); err != nil {
					fmt.Printf("  %2d. %-20s  INVALID (%v)\n", i+1, p.Name, err)
					continue
				}
				fmt.Printf("  %2d. %-20s  %s\n", i+1, p.Name, p.ExpandedConfigFile())
			}
			return nil
		},
	}
}
