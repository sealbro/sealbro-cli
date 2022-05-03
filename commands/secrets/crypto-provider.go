package secrets

import (
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"log"
	"os"
)

type CryptoProvider interface {
	Encrypt(rawText string) error
	Decrypt() (string, error)
}

type GpgCryptoProvider struct {
	privateKey string
	passphrase []byte
}

const AppDirectory = gitDirectory + "sealbro/"
const secretKeyPath = AppDirectory + "key.secret"
const encryptedSecrets = AppDirectory + "encrypted.secret"

func MakeGpgCryptoProvider(passphrase []byte) CryptoProvider {
	p := &GpgCryptoProvider{
		passphrase: passphrase,
	}

	err := p.generateKey()
	if err != nil {
		log.Fatal(err)
	}

	return p
}

func (p *GpgCryptoProvider) generateKey() error {
	checkGitDirectoryOrThrow()

	if _, err := os.Stat(AppDirectory); os.IsNotExist(err) {
		err = os.Mkdir(AppDirectory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(secretKeyPath); err == nil {
		// todo add expire check
		secretKey, err := os.ReadFile(secretKeyPath)
		if err != nil {
			err = os.Remove(secretKeyPath)
			if err != nil {
				return err
			}
		} else {
			p.privateKey = string(secretKey)
			return nil
		}
	}

	key, err := helper.GenerateKey("git", "git@secrets.com", p.passphrase, "rsa", 4096)
	if err != nil {
		return err
	}

	p.privateKey = key

	create, err := os.Create(secretKeyPath)
	if err != nil {
		return err
	}

	_, err = create.WriteString(key)
	if err != nil {
		return err
	}

	return create.Close()
}

func (p *GpgCryptoProvider) Encrypt(rawText string) error {
	armor, err := helper.EncryptMessageArmored(p.privateKey, rawText)
	if err != nil {
		return err
	}

	create, err := os.Create(encryptedSecrets)
	if err != nil {
		return err
	}

	_, err = create.WriteString(armor)
	if err != nil {
		return err
	}

	return create.Close()
}

func (p *GpgCryptoProvider) Decrypt() (string, error) {
	encryptedSecretsFile, err := os.ReadFile(encryptedSecrets)
	if err != nil {
		return "", err
	}

	return helper.DecryptMessageArmored(p.privateKey, p.passphrase, string(encryptedSecretsFile))
}

func CleanSecretsCache() error {
	return os.RemoveAll(encryptedSecrets)
}

func CleanSecretsCacheAll() error {
	return os.RemoveAll(AppDirectory)
}
