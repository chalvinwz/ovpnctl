package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chalvinwz/ovpnctl/internal/openvpn3"
	"github.com/spf13/cobra"
)

func disconnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disconnect SESSION|PROFILE",
		Short: "Disconnect by session number, full session path, or profile name",
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

	if n, err := strconv.Atoi(target); err == nil {
		if n < 1 || n > len(sessions) {
			return "", fmt.Errorf("invalid session number (use 'ovpnctl sessions' to list)")
		}
		return sessions[n-1].Path, nil
	}

	profile, err := getProfileCmd(target)
	if err != nil {
		return "", fmt.Errorf("argument must be session number, full session path, or profile name")
	}

	configPath := strings.ToLower(strings.TrimSpace(profile.ExpandedConfigFile()))
	configBase := strings.ToLower(filepath.Base(configPath))
	for _, s := range sessions {
		sessCfg := strings.ToLower(strings.TrimSpace(s.Config))
		if sessCfg == "" {
			continue
		}
		if sessCfg == configPath || strings.Contains(sessCfg, configPath) || filepath.Base(sessCfg) == configBase {
			return s.Path, nil
		}
	}

	return "", fmt.Errorf("no active session found for profile %q", profile.Name)
}
