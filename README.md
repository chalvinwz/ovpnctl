# ovpnctl

Lightweight CLI to manage **OpenVPN 3** sessions on Linux.

- List configured VPN profiles
- Connect to a profile (with OTP prompt)
- List active sessions
- Disconnect by session number or full session path
- Show installed CLI version

---

## Requirements

- Linux
- `openvpn3` CLI installed and available in `PATH`
- Go 1.22+ (only if building from source)

If `openvpn3` is not in your default `PATH`, you can set:

```bash
export OVPN3_BIN=/full/path/to/openvpn3
```

---

## Install

### Option A: Download from release (recommended)

```bash
curl -L -o ovpnctl https://github.com/chalvinwz/ovpnctl/releases/download/v0.1.0/ovpnctl-linux-amd64
chmod +x ovpnctl
sudo install -m 0755 ovpnctl /usr/local/bin/ovpnctl
```

### Option B: Build from source

```bash
git clone https://github.com/chalvinwz/ovpnctl.git
cd ovpnctl
go build -o ovpnctl ./cmd/ovpnctl
sudo install -m 0755 ovpnctl /usr/local/bin/ovpnctl
```

---

## Configuration

Create `~/.config/ovpnctl/profiles.yaml` (or pass custom file via `--config`).

You can start from:

```bash
cp examples/profiles.example.yaml ~/.config/ovpnctl/profiles.yaml
```

Example:

```yaml
profiles:
  - name: office
    config_file: /etc/openvpn/office.ovpn
    username: your.username
    password: your.password
    private_key_pass: ""

  - name: personal-fast
    config_file: ~/vpn/fast.ovpn
    username: your.username
    password: your.password
    private_key_pass: ""
```

### Profile notes

- `name`, `config_file`, `username`, `password` are required.
- `private_key_pass` is optional and can be empty or omitted.
- Profile names must be unique (case-insensitive).
- `config_file` supports `~/...` expansion.

---

## Usage

```bash
ovpnctl profiles
ovpnctl connect office
ovpnctl connect 1
ovpnctl sessions
ovpnctl disconnect 2
ovpnctl disconnect /net/openvpn/v3/sessions/<session-id>
ovpnctl version
```

Flags:

```text
--config string   path to profiles.yaml
                  (default search: ~/.config/ovpnctl/profiles.yaml, then ./profiles.yaml)
```

---

## Development

```bash
make build
make fmt
```

---

## Security

- Never commit real credentials.
- Keep local config files out of git.
- For team/production usage, consider secret injection (env/secret manager) instead of plaintext YAML.
