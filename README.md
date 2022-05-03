# sealbro-cli

Helpful CLI makes your daily programming routine simpler and safer

`go install github.com/sealbro/sealbro-cli@latest`

## sealbro-cli secrets -h

Stdout vault provider for [git-secrets](https://github.com/awslabs/git-secrets)

- set environment variables
  - `SEALBRO_PASSPHRASE` - (optional) passphrase for encryption key
  - `VAULT_ADDR` - vault address
  - `VAULT_TOKEN` - vault token
- add into [git-secrets](https://github.com/awslabs/git-secrets) sealbro-cli vault provider
  - `git secrets --add-provider -- sealbro-cli secrets show --path /devops/ --exclude ansible_user --exclude url --exclude username --exclude private --exclude public`
    - `--path` - vault kv
    - `--exclude` - vault secret key name to exclude from output

## Useful advices  

- export variables from `.env` file
  - `export $(grep -v '^#' .env | xargs -0)`
- uninstall module after `go install`
  - `rm -rf ~/go/bin/sealbro-cli`
- remove from git config
  - `git config --unset secrets.providers`