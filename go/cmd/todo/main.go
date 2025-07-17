package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Todo represents a single TODO item found in the codebase.
type Todo struct {
	File    string
	Line    int
	Message string
}

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

	var todos []Todo
	todoRegex := regexp.MustCompile(`^\s*(//|#)\s*TODO:\s*(.*)`)

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

		// Now, parse the file for TODOs
		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %q: %v\n", path, err)
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNumber := 0
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			matches := todoRegex.FindStringSubmatch(line)

			if len(matches) > 2 {
				message := strings.TrimSpace(matches[2])
				if strings.HasSuffix(message, ".") {
					todos = append(todos, Todo{
						File:    path,
						Line:    lineNumber,
						Message: message,
					})
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory %q: %v", absPath, err)
	}

	// Print all found TODOs
	for _, todo := range todos {
		fmt.Printf("%s:%d: %s\n", todo.File, todo.Line, todo.Message)
	}
}
