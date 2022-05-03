package secrets

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

const gitDirectory = "./.git/"
const GitConfig = gitDirectory + "config"

func RemoveProvider() error {
	gitSecretsSection := "secrets.providers"
	cmd := exec.Command("git", "config", "--unset", gitSecretsSection)

	return cmd.Run()
}

func InitProvider(paths []string, excludes []string) error {
	RemoveProvider()

	cmd := exec.Command("git", "secrets", "--install", "-f")
	err := cmd.Run()
	if err != nil {
		return err
	}

	builder := strings.Builder{}
	builder.WriteString("sealbro-cli secrets show")
	for _, path := range paths {
		builder.WriteString(" --path ")
		builder.WriteString(path)
	}

	for _, exclude := range excludes {
		builder.WriteString(" --exclude ")
		builder.WriteString(exclude)
	}

	cmd = exec.Command("git", "secrets", "--add-provider", "--", builder.String())

	return cmd.Run()
}

func checkGitDirectoryOrThrow() {
	if _, err := os.Stat(gitDirectory); os.IsNotExist(err) {
		log.Fatalln(errors.New("it isn't root directory (.git not found)"))
	}
}
