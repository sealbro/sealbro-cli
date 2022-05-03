# sealbro-cli

Helpful CLI makes your daily programming routine simpler and safer

`go install -v github.com/sealbro/sealbro-cli@latest`

## sealbro-cli secrets -h

Stdout vault provider for [git-secrets](https://github.com/awslabs/git-secrets)

### How use?

- set environment variables
  - `SEALBRO_PASSPHRASE` - (optional) passphrase for encryption key
  - `VAULT_ADDR` - vault address
  - `VAULT_TOKEN` - vault token
- install and setup [git-secrets](https://github.com/awslabs/git-secrets)
  - install on your macos/windows/linux 
    - macos `brew install git-secrets`
  - install git secrets for repository
    - `cd <your repository path>`
    - `git secrets --install -f`
  - add `sealbro-cli` vault provider into [git-secrets](https://github.com/awslabs/git-secrets)
    - `git secrets --add-provider -- sealbro-cli secrets show --path /devops/ --exclude ansible_user --exclude url --exclude username --exclude private --exclude public`
      - `--path` - vault kv
      - `--exclude` - vault secret key name to exclude from output
- or little easier after install `git-secrets`
  -  use `secrets init` (clean prev providers / install git-secrets for repository / add new providers global)
  - `sealbro-cli secrets init --path /devops/ --exclude ansible_user --exclude url --exclude username --exclude private --exclude public`

## Useful advices  

- export variables from `.env` file
  - `export $(grep -v '^#' .env | xargs -0)`
- macos set global environments zsh
  - `echo "export SEALBRO_PASSPHRASE=<passphrase>" > ~/.zshenv`
  - `echo "export VAULT_ADDR=<address>" > ~/.zshenv`
  - `echo "export VAULT_TOKEN=<token>" > ~/.zshenv`
  - `echo "export GOPATH=$HOME/go" > ~/.zshenv`
  - `echo "export PATH="$GOPATH/bin:$PATH" > ~/.zshenv`
- uninstall module after `go install`
  - `rm -rf $GOPATH/bin/sealbro-cli`
- remove from git config
  - `git config [--global] --unset secrets.providers`
