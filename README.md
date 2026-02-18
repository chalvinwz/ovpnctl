# ovpnctl

Simple CLI tool to manage **OpenVPN 3** connections on Linux.

- List configured VPN profiles
- Connect to a profile (prompts for OTP)
- List active sessions
- Disconnect sessions (by session-id or full path)

Designed for personal / team use with external YAML configuration — credentials are **not hardcoded** in the binary.

**Platform**: Linux only (requires official OpenVPN 3 Linux client)

## Features

- External `profiles.yaml` (supports multiple profiles)
- Username, password, private key stored outside the binary
- Connect by profile name **or** number
- OTP prompt only during connect (no full credential re-entry)
- List and disconnect active sessions with numbered shortcuts
- Clean error messages & helpful output

## Requirements

- Go 1.22 or newer (to build)
- [OpenVPN 3 Linux client](https://openvpn.net/community-downloads/) (`openvpn3` command must be in PATH)

## Installation

### Build ovpnctl

```bash
git clone https://github.com/chalvinwz/ovpnctl.git
cd ovpnctl
go build -o ovpnctl ./cmd/ovpnctl
```

Move to a directory in your PATH (optional):

```bash
sudo mv ovpnctl /usr/local/bin/
# or
mv ovpnctl ~/bin/
```

## Usage

```text
ovpnctl [flags] <command>

Commands:
  profiles    List configured VPN profiles
  connect     Start VPN session (prompts OTP)
  sessions    Show active OpenVPN 3 sessions
  disconnect  Stop a session

Flags:
  --config string   path to profiles.yaml
                    (default locations: ~/.config/ovpnctl/profiles.yaml)
```

### Configuration (profiles.yaml)

Create `~/.config/ovpnctl/profiles.yaml` (or any path — pass with `--config`)

```yaml
profiles:
  - name: office
    config_file: /etc/openvpn/office.ovpn
    username: chalvin@corp
    password: Very$ecret2026!
    private_key_pass: ""

  - name: personal-fast
    config_file: ~/vpn/fast.ovpn
    username: chalvin
    password: pass1234$$
    private_key_pass: ""
```

Example commands:

```bash
# List profiles
ovpnctl profiles

# Connect (by name or number)
ovpnctl connect office
ovpnctl connect 1

# Show active sessions
ovpnctl sessions

# Disconnect (by number from sessions list, or full path)
ovpnctl disconnect 2
ovpnctl disconnect /net/openvpn/v3/sessions/8bca2e2ds11e4s478csa988s969281af2804
```

## Development

```bash
# Run directly without building
go run ./cmd/ovpnctl -h

# Add new subcommand (example)
# → create internal/cmd/newcmd.go with NewNewcmdCmd() function
# → add rootCmd.AddCommand(newNewcmdCmd()) in root.go
```

## Security Notes

- Never commit real `profiles.yaml` to git
- Use `.gitignore` to exclude it
- Consider using secret management (pass, 1Password CLI, etc.) for production/team use
