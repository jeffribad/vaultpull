# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files with role-based filtering

---

## Installation

```bash
go install github.com/youruser/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/youruser/vaultpull/releases).

---

## Usage

Authenticate with your Vault instance and pull secrets into a `.env` file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-vault-token"

vaultpull --path secret/myapp --role developer --output .env
```

This will fetch all secrets at the given path permitted for the `developer` role and write them to `.env`.

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to sync | *(required)* |
| `--role` | Role used to filter accessible secrets | `default` |
| `--output` | Output file path | `.env` |
| `--overwrite` | Overwrite existing file if present | `false` |

### Example Output

```env
DB_HOST=postgres.internal
DB_PASSWORD=s3cr3t
API_KEY=abc123
```

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance (v1.12+)
- A valid Vault token or AppRole credentials

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 youruser