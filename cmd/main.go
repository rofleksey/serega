package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rofleksey/serega/internal/cli"
)

func main() {
	if err := cli.Execute(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "serega:", err)
		os.Exit(1)
	}
}
