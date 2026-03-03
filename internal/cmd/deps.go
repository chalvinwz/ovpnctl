package cmd

import (
	"github.com/chalvinwz/ovpnctl/internal/config"
	"github.com/chalvinwz/ovpnctl/internal/openvpn3"
)

var (
	requireBinaryCmd  = openvpn3.RequireBinary
	listSessionsCmd   = openvpn3.ListSessions
	disconnectCmdExec = openvpn3.Disconnect
	printSessionsCmd  = openvpn3.PrintSessions
	getProfileCmd     = config.GetProfile
	loadConfigCmd     = config.Load
	startSessionCmd   = openvpn3.StartSession
)
