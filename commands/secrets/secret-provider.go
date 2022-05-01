package secrets

import (
	"errors"
	vault "github.com/hashicorp/vault/api"
	"log"
	"time"
)

type SecretProvider interface {
	GetAllSecrets() ([]string, error)
}

type VaultSecretProvider struct {
	client   *vault.Client
	paths    []string
	excludes map[string]struct{}
}

func MakeVaultSecretProvider(paths, excludes []string) SecretProvider {

	// todo check VAULT_ADDR and VAULT_TOKEN environments

	config := vault.DefaultConfig()
	config.Timeout = time.Second * 10

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	mapExcludes := make(map[string]struct{})
	for _, exclude := range excludes {
		mapExcludes[exclude] = struct{}{}
	}

	return &VaultSecretProvider{client: client, paths: paths, excludes: mapExcludes}
}

func (p *VaultSecretProvider) GetAllSecrets() ([]string, error) {
	if len(p.paths) == 0 {
		return nil, errors.New("arg path empty")
	}

	uniqSecrets := make(map[string]struct{})

	var err error
	for _, path := range p.paths {
		errOne := p.getRecursiveUniqSecrets(uniqSecrets, path)

		if err == nil && errOne != nil {
			err = errOne
		}
	}

	var secretsList []string

	for value, _ := range uniqSecrets {
		secretsList = append(secretsList, value)
	}

	return secretsList, err
}

func (p *VaultSecretProvider) getRecursiveUniqSecrets(uniqSecrets map[string]struct{}, path string) error {
	if path[len(path)-1] != '/' {
		read, err := p.client.Logical().Read(path)
		if err != nil {
			return err
		}

		for key, value := range read.Data {
			if _, ok := p.excludes[key]; ok {
				continue
			}

			uniqSecrets[value.(string)] = struct{}{}
		}

		return nil
	}

	read, err := p.client.Logical().List(path)
	if err != nil {
		return err
	}

	if v, ok := read.Data["keys"]; ok {
		keys := v.([]interface{})
		for _, key := range keys {
			p.getRecursiveUniqSecrets(uniqSecrets, path+key.(string))
		}
	}

	return nil
}
