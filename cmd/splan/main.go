// Command splan works with structured planning documents (PRD, MRD, TRD, OKR, V2MOM, Roadmap).
package main

import (
	"fmt"
	"os"

	"github.com/grokify/prism-roadmap/cli"
)

func main() {
	// Override the command name for standalone use
	cli.RootCmd.Use = "splan"
	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
