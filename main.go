package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sealbro-cli/commands/secrets"
)

func main() {
	command := &secrets.SecretCommand{}

	encryptFlags := []cli.Flag{
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
						Usage: "get / encrypt / decrypt and compare secrets",
						Flags: encryptFlags,
						Action: func(c *cli.Context) error {
							passphrase := c.String("passphrase")
							paths := c.StringSlice("path")
							excludes := c.StringSlice("exclude")

							return command.Verify(passphrase, paths, excludes)
						},
					},
					{
						Name:  "init",
						Usage: "initialize git secrets --add-provider",
						Flags: encryptFlags,
						Action: func(c *cli.Context) error {
							paths := c.StringSlice("path")
							excludes := c.StringSlice("exclude")

							return secrets.InitProvider(paths, excludes)
						},
					},
					{
						Name: "remove",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "provider",
								Usage: "remove vault provider from " + secrets.GitConfig,
							},
							&cli.BoolFlag{
								Name:  "cache",
								Usage: "clean and remove vault provider cache from " + secrets.AppDirectory,
							},
						},
						Usage: "clean secrets",
						Action: func(c *cli.Context) error {
							isRemoveCache := c.Bool("cache")
							if isRemoveCache {
								return secrets.CleanSecretsCacheAll()
							}

							isRemoveProvider := c.Bool("provider")
							if isRemoveProvider {
								return secrets.RemoveProvider()
							}

							return nil
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
