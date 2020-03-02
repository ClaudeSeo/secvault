package cli

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type Info struct {
	Version string
}

func Initialize(info *Info) {
	var secretName string
	var outputType string

	app := cli.NewApp()
	app.Name = "secvault"
	app.Version = info.Version
	app.Usage = "simple cli tool to easily manage sensitive environment variable using AWS Secrets Manager"
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "ClaudeSeo",
			Email: "ehdaudtj@gmail.com",
		},
	}

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "secret-name",
			Usage:       "Secrest Name",
			Destination: &secretName,
		},
		&cli.StringFlag{
			Name:        "output-type",
			Usage:       "Output Type (null, dotenv, kubernetes, json)",
			Destination: &outputType,
		},
	}

	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "list",
			Usage: "Environment variable list in Secrets Manager",
			Action: func(ctx *cli.Context) error {
				List()
				return nil
			},
		},
		&cli.Command{
			Name:  "get",
			Usage: "Get environment variable stored in Secrets Manager",
			Flags: flags,
			Action: func(ctx *cli.Context) error {
				Get(secretName, outputType)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
