package main

import (
	"fmt"
	"os"

	"github.com/sudabon/dotenv_cryption/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
