package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sealbro-cli/commands/secrets"
)

func main() {
	// * secrets remove
	// * secrets init --path --excludes

	command := &secrets.SecretCommand{}

	encryptFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "vault-addr",
			EnvVars: []string{"VAULT_ADDR"},
		},
		&cli.StringFlag{
			Name:    "vault-token",
			EnvVars: []string{"VAULT_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "passphrase",
			EnvVars: []string{"SEALBRO_PASSPHRASE"},
		},
		&cli.StringSliceFlag{
			Name:     "path",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name: "exclude",
		},
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "secrets",
				Usage: "commands for processing secret values",
				Action: func(c *cli.Context) error {
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:  "show",
						Usage: "print secrets to stdout",
						Flags: encryptFlags,
						Action: func(c *cli.Context) error {
							passphrase := c.String("passphrase")
							paths := c.StringSlice("path")
							excludes := c.StringSlice("exclude")

							show, err := command.Show(passphrase, paths, excludes)
							if err != nil {
								return err
							}

							_, err = os.Stdout.WriteString(show)

							return err
						},
					},
					{
						Name:  "verify",
						Usage: "get, encrypt, decrypt and compare secrets",
						Flags: encryptFlags,
						Action: func(c *cli.Context) error {
							passphrase := c.String("passphrase")
							paths := c.StringSlice("path")
							excludes := c.StringSlice("exclude")

							return command.Verify(passphrase, paths, excludes)
						},
					},
					{
						Name:  "clean",
						Usage: "clean secrets cache",
						Action: func(c *cli.Context) error {
							return secrets.CleanSecretsCacheAll()
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
