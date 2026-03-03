package main

import (
	"os"

	"log/slog"

	"github.com/chalvinwz/ovpnctl/internal/cmd"
)

var (
	execute = cmd.Execute
	exitFn  = os.Exit
)

func main() {
	if err := execute(); err != nil {
		slog.Error("command failed", "err", err)
		exitFn(1)
	}
}
