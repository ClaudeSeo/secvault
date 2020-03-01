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
	var secretId string
	var outputType string

	app := cli.NewApp()
	app.Name = "secvalut"
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
			Name:        "secret-id",
			Usage:       "Secrest Id",
			Destination: &secretId,
		},
		&cli.StringFlag{
			Name:        "output-type",
			Usage:       "Output Type (null, dotenv, kubernetes, json)",
			Destination: &outputType,
		},
	}

	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "get",
			Usage: "Get environment variable stored in Secrets Manager",
			Flags: flags,
			Action: func(ctx *cli.Context) error {
				Get(secretId, outputType)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
