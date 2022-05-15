package secrets

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type SecretCommand struct {
}

func (c *SecretCommand) Copy(from, to string, cleanFrom, cleanTo bool) error {
	// get secrets recursive from KV
	// check exists KV
	// copy secrets to new KV
	// drop old KV

	secretProvider := MakeVaultSecretProvider([]string{})
	err := secretProvider.DeepCopyKV(fmt.Sprintf("/%s/", from), fmt.Sprintf("/%v/", to))

	return err
}

func (c *SecretCommand) Show(passphrase string, paths []string, excludes []string) (string, error) {
	cryptoProvider := MakeGpgCryptoProvider(generatePassphrase(passphrase, paths, excludes))

	decrypt, err := cryptoProvider.Decrypt()
	if err != nil {
		CleanSecretsCache()

		secretProvider := MakeVaultSecretProvider(excludes)
		allSecrets, err := secretProvider.GetAllUniqSecrets(paths)
		if err != nil {
			return "", err
		}

		secrets := strings.Join(allSecrets, "\n")
		err = cryptoProvider.Encrypt(secrets)

		return secrets, err
	}

	return decrypt, err
}

func (c *SecretCommand) Verify(passphrase string, paths []string, excludes []string) error {
	err := CleanSecretsCacheAll()
	if err != nil {
		return err
	}
	log.Println("Clean cache")

	secretProvider := MakeVaultSecretProvider(excludes)
	allSecrets, err := secretProvider.GetAllUniqSecrets(paths)
	if err != nil {
		return err
	}
	log.Println("Got secrets")

	secrets := strings.Join(allSecrets, "\n")

	cryptoProvider := MakeGpgCryptoProvider(generatePassphrase(passphrase, paths, excludes))
	err = cryptoProvider.Encrypt(secrets)
	if err != nil {
		return err
	}
	log.Println("Encrypted secrets")

	decrypt, err := cryptoProvider.Decrypt()
	if err != nil {
		return err
	}
	log.Println("Decrypted secrets")

	if secrets != decrypt {
		return errors.New("raw secrets not equal decrypt secrets")
	}

	err = CleanSecretsCacheAll()
	if err != nil {
		return err
	}
	log.Println("Clean cache")

	log.Println("Verify success!")

	return nil
}

func generatePassphrase(userPassphrase string, paths []string, excludes []string) []byte {
	// args hash for expire after changes
	argsPass := strings.Join(append(excludes, paths...), "|")
	hash := sha1.New()
	hash.Write([]byte(argsPass))
	sha := hash.Sum(nil)

	// for expire every month
	prefix := time.Now().Format("2006-01")

	return []byte(fmt.Sprintf("%v-%v-%x", prefix, userPassphrase, sha))
}

func (c *SecretCommand) Clean() {

}

func (c *SecretCommand) Remove() {

}

func (c *SecretCommand) Init() {

}
