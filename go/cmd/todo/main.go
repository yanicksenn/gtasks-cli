package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// By default, search the current directory. This can be overridden by a command-line argument.
	searchDir := "."
	flag.Parse()
	if flag.NArg() > 0 {
		searchDir = flag.Arg(0)
	}

	// Convert to an absolute path for cleaner output.
	absPath, err := filepath.Abs(searchDir)
	if err != nil {
		log.Fatalf("Error getting absolute path for %q: %v", searchDir, err)
	}

	fmt.Printf("Searching for TODOs in: %s\n", absPath)

	err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// If we can't access a path, report it and skip it.
			fmt.Fprintf(os.Stderr, "Error accessing path %q: %v\n", path, err)
			return nil
		}

		// Skip directories.
		if d.IsDir() {
			return nil
		}

		// For now, just print the file path.
		// Later, we will parse the file for TODOs here.
		fmt.Println(path)

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory %q: %v", absPath, err)
	}
}
