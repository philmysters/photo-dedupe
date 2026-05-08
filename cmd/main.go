package main

import (
	"fmt"
	"log"
	"os"

	"github.com/philmysters/photo-dedupe/internal"
)

func main() {
	fmt.Println("photo-dedupe - A utility for deduplicating and organizing photos")

	if len(os.Args) < 4 {
		fmt.Println("Usage: photo-dedupe -in1 <input_folder_1> -in2 <input_folder_2> -out <output_folder> [-config <config_file>] [--dryrun]")
		os.Exit(1)
	}

	input1, input2, output := "", "", ""
	var configPath = "photo_dedupe.yaml"
	var dryrun bool

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-in1" && i+1 < len(os.Args) {
			input1 = os.Args[i+1]
			i++
		} else if arg == "-in2" && i+1 < len(os.Args) {
			input2 = os.Args[i+1]
			i++
		} else if arg == "-out" && i+1 < len(os.Args) {
			output = os.Args[i+1]
			i++
		} else if arg == "-config" && i+1 < len(os.Args) {
			configPath = os.Args[i+1]
			i++
		} else if arg == "--dryrun" {
			dryrun = true
		}
	}

	if input1 == "" || input2 == "" || output == "" {
		fmt.Println("Missing one or more required arguments: -in1, -in2, -out")
		os.Exit(2)
	}

	cfg, err := internal.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Could not load config file (%s): %v", configPath, err)
	}
	fmt.Printf("Loaded config: %+v\n", cfg)

	fmt.Printf("Input 1: %s\nInput 2: %s\nOutput: %s\nConfig: %s\nDry run: %t\n", input1, input2, output, configPath, dryrun)

	err = os.MkdirAll(output, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	photoFiles1, err := internal.FindPhotoFiles(input1, cfg.SupportedExtensions)
	if err != nil {
		log.Fatalf("Error scanning %s: %v", input1, err)
	}
	photoFiles2, err := internal.FindPhotoFiles(input2, cfg.SupportedExtensions)
	if err != nil {
		log.Fatalf("Error scanning %s: %v", input2, err)
	}

	fmt.Printf("Discovered %d photo(s) in %s and %d photo(s) in %s\n", len(photoFiles1), input1, len(photoFiles2), input2)
	// Placeholder: Next step would be deduplication, hashing, output...

	for _, pf := range photoFiles1 {
		fmt.Printf("Input1 file: %s\n", pf.Path)
	}
	for _, pf := range photoFiles2 {
		fmt.Printf("Input2 file: %s\n", pf.Path)
	}
}
