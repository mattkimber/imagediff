package main

import (
	"flag"
	"fmt"
	"internal/differ"
)

var (
	GenerateMode bool
	ApplyMode bool
	InputDir string
	OutputDir string
	CompareDir string
)

func init() {
	flag.StringVar(&InputDir, "input", "", "Input directory")
	flag.StringVar(&OutputDir, "output", "", "Output directory")
	flag.StringVar(&CompareDir, "compare", "", "Comparison directory (when generating images) or diff directory (when applying)")
	flag.BoolVar(&GenerateMode, "generate", false, "Set to true to generate diffs")
	flag.BoolVar(&ApplyMode, "apply", false, "Set to true to apply diffs")

	flag.Parse()
}

func main() {
	if InputDir == "" || OutputDir == "" || CompareDir == "" {
		fmt.Printf("All three directories (-input, -output, -compare) must be specified")
		return
	}

	if ApplyMode && !GenerateMode {
		differ.Apply(InputDir, CompareDir, OutputDir)
	} else if GenerateMode && !ApplyMode {
		differ.Generate(InputDir, CompareDir, OutputDir)
	} else {
		fmt.Printf("Must run in either apply mode (with -apply) or generate mode (with -generate)")
		return
	}
}