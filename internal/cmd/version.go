package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show ovpnctl version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ovpnctl %s\n", Version)
		},
	}
}
