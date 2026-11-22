package main

import (
	"fmt"
	"os"

	"github.com/CoCaptinCrazy/secret-bitpack/internal/cli"
)

func main() {
	if err := cli.ExecuteWithIO(os.Args[1:], os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
