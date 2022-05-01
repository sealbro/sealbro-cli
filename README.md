# sealbro-cli
CLI with helpful methods for me

## Vault provider for git-secrets

- set environment variables
  - `SEALBRO_PASSPHRASE` - (optional) passphrase for encryption key
  - `VAULT_ADDR` - vault address
  - `VAULT_TOKEN` - vault token
- command for export variables from `.env` file
  - `export $(grep -v '^#' .env | xargs -0)`
- install
  - `go install`
- [git-secrets](https://github.com/awslabs/git-secrets) add vault provider
  - `git secrets --add-provider -- sealbro-cli secrets show --path /devops/ --exclude ansible_user --exclude url --exclude username --exclude private --exclude public`
    - `--path` - vault kv
    - `--exclude` - vault secret key name to exclude from output
- `sealbro-cli secrets -h` - more information about commands
- uninstall
  - `rm -rf ~/go/bin/sealbro-cli`