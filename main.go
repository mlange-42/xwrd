package main

import (
	"fmt"
	"os"

	"github.com/mlange-42/xwrd/cli"
	"github.com/mlange-42/xwrd/util"
)

const version = "0.1.0-dev"

func main() {
	util.EnsureDirs()

	if err := cli.RootCommand(version).Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}
