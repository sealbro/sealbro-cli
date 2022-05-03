package secrets

import (
	"errors"
	vault "github.com/hashicorp/vault/api"
	"log"
	"os"
	"sync"
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
	config := vault.DefaultConfig()
	config.Timeout = time.Second * 10

	if len(os.Getenv("VAULT_ADDR")) == 0 || len(os.Getenv("VAULT_TOKEN")) == 0 {
		log.Fatalf("VAULT_ADDR or VAULT_TOKEN have not been defined")
	}

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

	wgVaultRequests := sync.WaitGroup{}
	secretChan := make(chan string)
	uniqSecretsChan := make(chan []string)
	defer close(uniqSecretsChan)

	go collectSecrets(&wgVaultRequests, secretChan, uniqSecretsChan)

	var err error

	wgVaultRequests.Add(len(p.paths))
	for _, path := range p.paths {
		p.getRecursiveUniqSecrets(&wgVaultRequests, secretChan, path)
	}
	wgVaultRequests.Wait()
	close(secretChan)

	return <-uniqSecretsChan, err
}

func collectSecrets(wgVaultRequests *sync.WaitGroup, secretChan <-chan string, uniqSecretsChan chan<- []string) {
	uniqSecrets := make(map[string]struct{})

	for {
		secret, ok := <-secretChan
		if !ok {
			break
		}

		uniqSecrets[secret] = struct{}{}
		wgVaultRequests.Done()
	}

	var secretsList []string

	for value := range uniqSecrets {
		secretsList = append(secretsList, value)
	}

	uniqSecretsChan <- secretsList
}

func (p *VaultSecretProvider) getRecursiveUniqSecrets(wgVaultRequests *sync.WaitGroup, secretChan chan<- string, path string) {
	defer wgVaultRequests.Done()

	if path[len(path)-1] != '/' {
		read, err := p.client.Logical().Read(path)
		if err != nil {
			log.Fatalf("cann't read secret from vault: %v", err)
		}

		for key, value := range read.Data {
			if _, ok := p.excludes[key]; ok {
				continue
			}

			wgVaultRequests.Add(1)
			secretChan <- value.(string)
		}

		return
	}

	read, err := p.client.Logical().List(path)
	if err != nil {
		log.Fatalf("cann't read list from vault: %v", err)
	}

	if v, ok := read.Data["keys"]; ok {
		keys := v.([]interface{})

		wgVaultRequests.Add(len(keys))

		for _, key := range keys {
			go p.getRecursiveUniqSecrets(wgVaultRequests, secretChan, path+key.(string))
		}
	}
}
