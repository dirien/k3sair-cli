package main

import (
	"fmt"
	"github.com/k3sair/cmd/k3sair"
	"os"
)

var (
	version string
	commit  string
)

func main() {
	if err := k3sair.Execute(version, commit); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
