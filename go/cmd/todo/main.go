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
	"sort"
	"strings"
	"time"
)

var (
	// Maps file extensions to their single-line comment markers.
	lineCommentMarkers = map[string]string{
		".go":   "//",
		".py":   "#",
		".js":   "//",
		".ts":   "//",
		".java": "//",
		".rs":   "//",
		".sh":   "#",
		".rb":   "#",
	}
)

// Todo represents a single TODO item found in the codebase.
type Todo struct {
	File    string
	Line    int
	Message string
	ModTime time.Time
}

// FileTodoCount holds the count of TODOs for a single file.
type FileTodoCount struct {
	File  string
	Count int
}

func main() {
	// CLI Flags
	aggregated := flag.Bool("aggregated", false, "Display TODOs in an aggregated view by file.")
	
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

	todosByFile := make(map[string][]Todo)

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

		// Determine the comment marker based on file extension.
		ext := filepath.Ext(path)
		marker, supported := lineCommentMarkers[ext]
		if !supported {
			return nil // Skip files with unsupported extensions.
		}

		fileInfo, err := d.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting file info for %q: %v\n", path, err)
			return nil
		}
		modTime := fileInfo.ModTime()

		// Now, parse the file for TODOs
		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %q: %v\n", path, err)
			return nil
		}
		defer file.Close()

		todoRegex := regexp.MustCompile(fmt.Sprintf(`^\s*%s\s*TODO:\s*(.*)`, regexp.QuoteMeta(marker)))
		
		scanner := bufio.NewScanner(file)
		lineNumber := 0
		var currentTodo *Todo
		for scanner.Scan() {
			lineNumber++
			line := scanner.Text()
			matches := todoRegex.FindStringSubmatch(line)

			if len(matches) > 1 { // Found a new TODO
				if currentTodo != nil && strings.HasSuffix(currentTodo.Message, ".") {
					todosByFile[currentTodo.File] = append(todosByFile[currentTodo.File], *currentTodo)
				}
				message := strings.TrimSpace(matches[1])
				currentTodo = &Todo{
					File:    path,
					Line:    lineNumber,
					Message: message,
					ModTime: modTime,
				}
			} else if currentTodo != nil { // Potentially part of a multi-line TODO
				trimmedLine := strings.TrimSpace(line)
				// Simple heuristic: if the line is a comment, append it.
				if strings.HasPrefix(trimmedLine, marker) {
					currentTodo.Message += " " + strings.TrimSpace(strings.TrimPrefix(trimmedLine, marker))
				} else {
					// Not a comment, so the multi-line TODO ends here.
					if strings.HasSuffix(currentTodo.Message, ".") {
						todosByFile[currentTodo.File] = append(todosByFile[currentTodo.File], *currentTodo)
					}
					currentTodo = nil
				}
			}
		}
		// Add the last todo if it exists
		if currentTodo != nil && strings.HasSuffix(currentTodo.Message, ".") {
			todosByFile[currentTodo.File] = append(todosByFile[currentTodo.File], *currentTodo)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory %q: %v", absPath, err)
	}

	if *aggregated {
		printAggregated(todosByFile)
	} else {
		printStandard(todosByFile)
	}
}

func printStandard(todosByFile map[string][]Todo) {
	var allTodos []Todo
	for _, todos := range todosByFile {
		allTodos = append(allTodos, todos...)
	}

	// Sort the todos by modification time (newest first).
	sort.Slice(allTodos, func(i, j int) bool {
		return allTodos[i].ModTime.After(allTodos[j].ModTime)
	})

	// Print all found TODOs
	for _, todo := range allTodos {
		fmt.Printf("%s:%d: %s\n", todo.File, todo.Line, todo.Message)
	}
}

func printAggregated(todosByFile map[string][]Todo) {
	counts := make([]FileTodoCount, 0, len(todosByFile))
	for file, todos := range todosByFile {
		counts = append(counts, FileTodoCount{File: file, Count: len(todos)})
	}

	// Sort by count (most first)
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	for _, item := range counts {
		fmt.Printf("%s: %d TODOs\n", item.File, item.Count)
	}
}
