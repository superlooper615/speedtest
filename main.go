package main

import (
	"os"

	"github.com/superlooper615/speedtest/cmd"

	_ "github.com/superlooper615/speedtest/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
