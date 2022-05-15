package main

import (
	"github.com/sealbro/sealbro-cli/commands/secrets"
	"github.com/urfave/cli/v2"
	"log"
	"os"
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
						Name:  "copy",
						Usage: "copy recursive Vault KV secrets",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "from",
								Usage:    "KV name from copy",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "to",
								Usage:    "KV name to copy",
								Required: true,
							},
							&cli.BoolFlag{
								Name:  "clean-from",
								Usage: "drop KV from after copy",
							},
							&cli.BoolFlag{
								Name:  "clean-to",
								Usage: "drop KV to before copy",
							},
						},
						Action: func(c *cli.Context) error {
							from := c.String("from")
							to := c.String("to")
							cleanFrom := c.Bool("clean-from")
							cleanTo := c.Bool("clean-to")

							return command.Copy(from, to, cleanFrom, cleanTo)
						},
					},
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
						Usage: "initialize git secrets --add-provider --global",
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
								Usage: "remove vault provider from git config --global",
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
