package secrets

import (
	"errors"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type SecretProvider interface {
	GetAllUniqSecrets(paths []string) ([]string, error)
	DeepCopyKV(from, to string) error
}

type VaultSecretProvider struct {
	client   *vault.Client
	excludes map[string]struct{}
}

type secretKeyValue struct {
	path, key, value string
}

func MakeVaultSecretProvider(excludes []string) SecretProvider {
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

	return &VaultSecretProvider{client: client, excludes: mapExcludes}
}

func (p *VaultSecretProvider) DeepCopyKV(from, to string) error {
	secretKeyValues, err := p.getSecrets([]string{from})
	if err != nil {
		return err
	}

	for _, keyValue := range secretKeyValues {
		log.Println(keyValue.path, keyValue.key, keyValue.value)
	}

	return err
}

func (p *VaultSecretProvider) GetAllUniqSecrets(paths []string) ([]string, error) {
	secretKeyValues, err := p.getSecrets(paths)
	if err != nil {
		return nil, err
	}
	uniqSecrets := make(map[string]struct{})

	for _, secret := range secretKeyValues {
		uniqSecrets[secret.value] = struct{}{}
	}

	var secretsList []string
	for key := range uniqSecrets {
		secretsList = append(secretsList, key)
	}

	return secretsList, err
}

func (p *VaultSecretProvider) getSecrets(paths []string) ([]secretKeyValue, error) {
	if len(paths) == 0 {
		return nil, errors.New("arg path empty")
	}

	wgVaultRequests := sync.WaitGroup{}
	secretChan := make(chan secretKeyValue)
	uniqSecretsChan := make(chan []secretKeyValue)
	defer close(uniqSecretsChan)

	go collectSecretKeyValues(&wgVaultRequests, secretChan, uniqSecretsChan)

	var err error

	wgVaultRequests.Add(len(paths))
	for _, path := range paths {
		p.getRecursiveUniqSecrets(&wgVaultRequests, secretChan, path)
	}
	wgVaultRequests.Wait()
	close(secretChan)

	secretKeyValues := <-uniqSecretsChan

	return secretKeyValues, err
}

func collectSecretKeyValues(wgVaultRequests *sync.WaitGroup, secretChan <-chan secretKeyValue, uniqSecretsChan chan<- []secretKeyValue) {
	var secretsList []secretKeyValue

	for {
		secret, ok := <-secretChan
		if !ok {
			break
		}

		secretsList = append(secretsList, secret)
		wgVaultRequests.Done()
	}

	uniqSecretsChan <- secretsList
}

func (p *VaultSecretProvider) getRecursiveUniqSecrets(wgVaultRequests *sync.WaitGroup, secretChan chan<- secretKeyValue, path string) {
	defer wgVaultRequests.Done()

	if path[len(path)-1] != '/' {

		// kv2
		correctedPath := strings.Replace(path, "metadata/", "data/", 1)

		read, err := p.client.Logical().Read(correctedPath)
		if err != nil {
			log.Fatalf("cann't read secret from vault: %v", err)
		}

		var data map[string]interface{}
		if correctedPath != path {
			data = read.Data["data"].(map[string]interface{})
		} else {
			data = read.Data
		}

		for key, value := range data {
			if _, ok := p.excludes[key]; ok {
				continue
			}

			wgVaultRequests.Add(1)
			secretChan <- secretKeyValue{
				path:  path,
				key:   key,
				value: fmt.Sprintf("%v", value),
			}
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
