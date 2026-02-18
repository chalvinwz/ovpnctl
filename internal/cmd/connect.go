package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chalvinwz/ovpnctl/internal/config"
	"github.com/chalvinwz/ovpnctl/internal/openvpn3"
	"github.com/spf13/cobra"
)

func connectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect PROFILE",
		Short: "Start VPN session (prompts OTP only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := config.GetProfile(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Connecting to %q (%s)\n", profile.Name, profile.ConfigFile)

			fmt.Print("OTP: ")
			scanner := bufio.NewScanner(os.Stdin)
			if !scanner.Scan() {
				return fmt.Errorf("failed to read OTP")
			}
			otp := strings.TrimSpace(scanner.Text())

			if err := openvpn3.StartSession(profile, otp); err != nil {
				return fmt.Errorf("session-start failed: %w", err)
			}

			fmt.Println("Connection established.")
			return nil
		},
	}
}
