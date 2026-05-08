package main

import (
	"fmt"
	"io"
	"os"

	"github.com/philmysters/photo-dedupe/internal"
)

var osExit = os.Exit

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		osExit(1)
	}
}

func run(args []string, out io.Writer) error {
	_, _ = fmt.Fprintln(out, "photo-dedupe - A utility for deduplicating and organizing photos")

	if len(args) < 3 {
		return fmt.Errorf("usage: photo-dedupe -in1 <input_folder_1> -in2 <input_folder_2> -out <output_folder> [-config <config_file>] [--dryrun]")
	}

	input1, input2, output := "", "", ""
	configPath := "photo_dedupe.yaml"
	var dryrun bool

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-in1" && i+1 < len(args):
			input1 = args[i+1]
			i++
		case arg == "-in2" && i+1 < len(args):
			input2 = args[i+1]
			i++
		case arg == "-out" && i+1 < len(args):
			output = args[i+1]
			i++
		case arg == "-config" && i+1 < len(args):
			configPath = args[i+1]
			i++
		case arg == "--dryrun":
			dryrun = true
		}
	}

	if input1 == "" || input2 == "" || output == "" {
		return fmt.Errorf("missing one or more required arguments: -in1, -in2, -out")
	}

	cfg, err := internal.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("could not load config file (%s): %w", configPath, err)
	}
	_, _ = fmt.Fprintf(out, "Loaded config: %+v\n", cfg)
	_, _ = fmt.Fprintf(out, "Input 1: %s\nInput 2: %s\nOutput: %s\nConfig: %s\nDry run: %t\n", input1, input2, output, configPath, dryrun)

	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	photoFiles1, err := internal.FindPhotoFiles(input1, cfg.SupportedExtensions)
	if err != nil {
		return fmt.Errorf("error scanning %s: %w", input1, err)
	}
	photoFiles2, err := internal.FindPhotoFiles(input2, cfg.SupportedExtensions)
	if err != nil {
		return fmt.Errorf("error scanning %s: %w", input2, err)
	}

	_, _ = fmt.Fprintf(out, "Discovered %d photo(s) in %s and %d photo(s) in %s\n", len(photoFiles1), input1, len(photoFiles2), input2)

	for _, pf := range photoFiles1 {
		_, _ = fmt.Fprintf(out, "Input1 file: %s\n", pf.Path)
	}
	for _, pf := range photoFiles2 {
		_, _ = fmt.Fprintf(out, "Input2 file: %s\n", pf.Path)
	}
	return nil
}
