//go:build ignore

// This program generates JSON Schema files from Go types.
// Run with: go run schema/cmd/generate/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grokify/prism-roadmap/schema"
)

func main() {
	gen := schema.NewGenerator()

	// Get schema directory relative to working directory
	schemaDir := "schema"
	if len(os.Args) > 1 {
		schemaDir = os.Args[1]
	}

	// Make absolute path
	absDir, err := filepath.Abs(schemaDir)
	if err != nil {
		fmt.Printf("Error resolving path: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generating schemas to %s...\n", absDir)

	// Generate all schemas
	if err := gen.GenerateAll(absDir); err != nil {
		fmt.Printf("Error generating schemas: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Schemas generated successfully!")
	fmt.Println("Generated files:")
	fmt.Println("  - prd.schema.json")
	fmt.Println("  - okr.schema.json")
	fmt.Println("  - v2mom.schema.json")
	fmt.Println("  - shapeup-pitch.schema.json")
	fmt.Println("  - shapeup-bet.schema.json")
	fmt.Println("  - shapeup-scope.schema.json")
	fmt.Println("  - discovery-snapshot.schema.json")
	fmt.Println("  - assumption-map.schema.json")
	fmt.Println("  - experience-map.schema.json")
	fmt.Println("  - leanstartup.schema.json")
	fmt.Println("  - designthinking.schema.json")
	fmt.Println("  - jtbd.schema.json")
	fmt.Println("  - journey-roadmap.schema.json")
}
