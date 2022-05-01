package secrets

import (
	"errors"
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

const gitDirectory = "./.git/"
const appDirectory = gitDirectory + "sealbro/"
const secretKeyPath = appDirectory + "key.secret"
const encryptedSecrets = appDirectory + "encrypted.secret"

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
	if _, err := os.Stat(gitDirectory); os.IsNotExist(err) {
		return errors.New("it isn't root directory (.git not found)")
	}

	if _, err := os.Stat(appDirectory); os.IsNotExist(err) {
		err = os.Mkdir(appDirectory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(secretKeyPath); err == nil {
		// todo add expire check
		secretKey, err := os.ReadFile(secretKeyPath)
		if err != nil {
			os.Remove(secretKeyPath)
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
	create.Close()

	return err
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
	create.Close()

	return err
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
	return os.RemoveAll(appDirectory)
}
