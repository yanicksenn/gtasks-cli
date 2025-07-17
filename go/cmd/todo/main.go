package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
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
		// Handle potential errors walking the path
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %q: %v\n", path, err)
			return nil // Continue walking
		}

		// Heuristic to skip common ignored directories and hidden files/dirs
		if d.IsDir() {
			dirName := d.Name()
			if dirName == ".git" || dirName == ".hg" || dirName == "node_modules" || strings.HasPrefix(dirName, "bazel-") {
				return filepath.SkipDir // Skip this directory and all its contents
			}
		}
		
		// Skip all hidden files and directories (e.g. .DS_Store, .idea)
		if strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil // Skip hidden file
		}

		// We only want to parse files, not directories
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
