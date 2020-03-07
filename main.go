package main

import (
	"github.com/claudeseo/secvault/cli"
)

func main() {
	cli.Initialize(&cli.Info{
		AppName:     "secvault",
		Description: "simple cli tool to easily manage sensitive environment variable using AWS Secrets Manager",
		AuthorName:  "ClaudeSeo",
		AuthorEmail: "ehdaudtj@gmail.com",
		Version:     "0.0.1",
	})
}
