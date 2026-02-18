package cmd

import (
	"fmt"
	"strconv"

	"github.com/chalvinwz/ovpnctl/internal/openvpn3"
	"github.com/spf13/cobra"
)

func disconnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disconnect SESSION",
		Short: "Disconnect session (number from 'sessions' or full path)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]

			var path string
			if openvpn3.LooksLikeSessionPath(target) {
				path = target
			} else {
				n, err := strconv.Atoi(target)
				if err != nil {
					return fmt.Errorf("argument must be a number or full session path")
				}

				sessions, err := openvpn3.ListSessions()
				if err != nil {
					return err
				}

				if n < 1 || n > len(sessions) {
					return fmt.Errorf("invalid session number (use 'ovpnctl sessions' to list)")
				}

				path = sessions[n-1].Path
			}

			if err := openvpn3.Disconnect(path); err != nil {
				return err
			}

			fmt.Printf("Disconnected %s\n", path)

			fmt.Println("\nRemaining sessions:")
			_ = openvpn3.PrintSessions() // best-effort
			return nil
		},
	}
}
