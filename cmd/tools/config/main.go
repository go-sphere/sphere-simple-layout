//go:build spheretools

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/go-sphere/sphere-simple-layout/internal/config"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	output := fs.String("output", "config_gen.json", "output file path")
	if len(args) > 0 && args[0] == "gen" {
		args = args[1:]
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(config.NewEmptyConfig(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(*output, raw, 0o644)
}
