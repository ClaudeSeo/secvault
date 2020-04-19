package cli

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type Info struct {
	AppName     string
	Description string
	Version     string
	AuthorName  string
	AuthorEmail string
}

func Initialize(info *Info) {
	var secretName string
	var file string
	var fileType string

	app := cli.NewApp()
	app.Name = info.AppName
	app.Description = info.Description
	app.Version = info.Version
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  info.AuthorName,
			Email: info.AuthorEmail,
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
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "secret-name",
					Usage:       "Secrest Name",
					Destination: &secretName,
				},
				&cli.StringFlag{
					Name:        "type",
					Usage:       "Output File Type (support: dotenv, kubernetes, json)",
					Destination: &fileType,
				},
			},
			Action: func(ctx *cli.Context) error {
				Get(secretName, fileType)
				return nil
			},
		},
		&cli.Command{
			Name:  "put",
			Usage: "Put environment variable to Secrets Manager",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "secret-name",
					Usage:       "Secrest Name",
					Destination: &secretName,
				},
				&cli.StringFlag{
					Name:        "file",
					Usage:       "File Path",
					Destination: &file,
				},
				&cli.StringFlag{
					Name:        "type",
					Usage:       "Input File Type (support: dotenv, json)",
					Destination: &fileType,
				},
			},
			Action: func(ctx *cli.Context) error {
				Put(secretName, file, fileType)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
