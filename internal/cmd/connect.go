package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func connectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect PROFILE",
		Short: "Start VPN session (prompts OTP only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireBinaryCmd(); err != nil {
				return err
			}

			profile, err := getProfileCmd(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Connecting to %q (%s)\n", profile.Name, profile.ExpandedConfigFile())

			fmt.Print("OTP: ")
			otp, err := readOTP(os.Stdin)
			if err != nil {
				return err
			}

			if err := startSessionCmd(profile, otp); err != nil {
				return fmt.Errorf("session-start failed: %w", err)
			}

			fmt.Println("Connection established.")
			return nil
		},
	}
}

func readOTP(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		return "", fmt.Errorf("failed to read OTP")
	}
	otp := strings.TrimSpace(scanner.Text())
	if otp == "" {
		return "", fmt.Errorf("OTP cannot be empty")
	}
	return otp, nil
}
