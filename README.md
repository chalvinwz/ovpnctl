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
- Connect by profile name (**case-insensitive**) or number
- `~` expansion for profile `config_file`
- Per-profile validation with clear invalid-profile output
- OpenVPN binary preflight check (`openvpn3` must exist in PATH)
- OTP prompt only during connect (no full credential re-entry)
- List/disconnect active sessions with numbered shortcuts
- Clean error messages and helpful output

## Requirements

- Go 1.22+ (to build)
- [OpenVPN 3 Linux client](https://openvpn.net/community-downloads/) (`openvpn3` command must be in PATH)

## Installation

```bash
git clone https://github.com/chalvinwz/ovpnctl.git
cd ovpnctl
go build -o ovpnctl ./cmd/ovpnctl
```

Optional:

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
                    (default search: ~/.config/ovpnctl/profiles.yaml, then ./profiles.yaml)
```

## Configuration (`profiles.yaml`)

Create `~/.config/ovpnctl/profiles.yaml` (or pass custom path via `--config`).
You can start by copying `examples/profiles.example.yaml`.

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

### Notes

- Profile names must be unique (case-insensitive).
- Required profile fields: `name`, `config_file`, `username`, `password`.
- `connect` requires non-empty OTP input.

## Examples

```bash
# List profiles
ovpnctl profiles

# Connect (name or number)
ovpnctl connect office
ovpnctl connect 1

# Active sessions
ovpnctl sessions

# Disconnect (number from sessions list, or full path)
ovpnctl disconnect 2
ovpnctl disconnect /net/openvpn/v3/sessions/8bca2e2ds11e4s478csa988s969281af2804
```

## Development

```bash
# Run directly
go run ./cmd/ovpnctl -h

# Common tasks
make build
make test
make cover
make fmt
```

## Project structure

```text
cmd/ovpnctl/         CLI entrypoint
internal/cmd/        Cobra commands + command orchestration
internal/config/     Config model, loading, and validation
internal/openvpn3/   OpenVPN3 process integration/parsing
examples/            Example configuration files
```

## Security notes

- Never commit real `profiles.yaml` to git.
- Use `.gitignore` to exclude local config.
- For team/prod use, consider loading secrets from a secret manager instead of plain YAML.
