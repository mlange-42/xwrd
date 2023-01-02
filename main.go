package main

import (
	"fmt"
	"os"

	"github.com/mlange-42/xwrd/cli"
)

const version = "0.1.0-dev"

func main() {
	if err := cli.RootCommand(version).Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}
