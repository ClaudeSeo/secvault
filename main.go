package main

import (
	"github.com/claudeseo/secvault/cli"
)

func main() {
	cli.Initialize(&cli.Info{
		Version: "0.0.1",
	})
}
