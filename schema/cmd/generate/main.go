//go:build ignore

// This program generates JSON Schema files for prism-roadmap types.
// Run from the schema directory:
//
//	go run cmd/generate/main.go
package main

import (
	"fmt"
	"os"

	"github.com/grokify/prism-roadmap/schema"
)

func main() {
	g := schema.NewGenerator()

	// Generate schemas to the schema directory
	outputDir := "."
	if err := g.GenerateAll(outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated all prism-roadmap schemas")
}
