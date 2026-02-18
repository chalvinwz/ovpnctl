package cmd

import (
	"fmt"

	"github.com/chalvinwz/ovpnctl/internal/openvpn3"
	"github.com/spf13/cobra"
)

func sessionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "sessions",
		Aliases: []string{"ls-active", "active"},
		Short:   "List active OpenVPN sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			sessions, err := openvpn3.ListSessions()
			if err != nil {
				return err
			}

			if len(sessions) == 0 {
				fmt.Println("No active sessions.")
				return nil
			}

			fmt.Println("Active sessions:")
			for i, s := range sessions {
				fmt.Printf("  %2d. %s\n", i+1, s.Path)
				if s.Config != "" {
					fmt.Printf("     Config: %s\n", s.Config)
				}
			}
			return nil
		},
	}
}
