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
		Short: "Disconnect by session number or full session path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireBinaryCmd(); err != nil {
				return err
			}

			path, err := resolveDisconnectPath(args[0])
			if err != nil {
				return err
			}

			if err := disconnectCmdExec(path); err != nil {
				return err
			}

			fmt.Printf("Disconnected %s\n", path)
			fmt.Println("\nRemaining sessions:")
			_ = printSessionsCmd() // best-effort
			return nil
		},
	}
}

func resolveDisconnectPath(target string) (string, error) {
	if openvpn3.LooksLikeSessionPath(target) {
		return target, nil
	}

	sessions, err := listSessionsCmd()
	if err != nil {
		return "", err
	}

	n, err := strconv.Atoi(target)
	if err != nil {
		return "", fmt.Errorf("argument must be a number or full session path")
	}
	if n < 1 || n > len(sessions) {
		return "", fmt.Errorf("invalid session number (use 'ovpnctl sessions' to list)")
	}
	return sessions[n-1].Path, nil
}
