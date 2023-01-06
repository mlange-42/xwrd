package main

import (
	"fmt"
	"os"

	"github.com/mlange-42/xwrd/cli"
	"github.com/mlange-42/xwrd/core"
	"github.com/mlange-42/xwrd/util"
)

const version = "0.1.2"

func main() {
	util.EnsureDirs()

	config, err := core.LoadConfig()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	if err := cli.RootCommand(&config, version).Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}
