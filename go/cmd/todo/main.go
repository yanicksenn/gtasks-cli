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
	lineCommentMarkers = map[string]string{
		".go": "//", ".py": "#", ".js": "//", ".ts": "//", ".java": "//", ".rs": "//", ".sh": "#", ".rb": "#",
	}
)

type Todo struct {
	File    string
	Line    int
	Message string
	ModTime time.Time
}

type InvalidTodo struct {
	File    string
	Line    int
	Content string
	Reason  string
}

type FileTodoCount struct {
	File  string
	Count int
}

func main() {
		aggregated := flag.Bool("aggregated", false, "Group TODOs by file and display the count for each file.")
	validate := flag.Bool("validate", false, "Validate TODO format and exit with an error if invalid TODOs are found.")
	searchDir := "."
	flag.Parse()
	if flag.NArg() > 0 {
		searchDir = flag.Arg(0)
	}

	absPath, err := filepath.Abs(searchDir)
	if err != nil {
		log.Fatalf("Error getting absolute path for %q: %v", searchDir, err)
	}

	todosByFile, invalidTodos, err := findTodos(absPath)
	if err != nil {
		log.Fatalf("Error walking directory %q: %v", absPath, err)
	}

	if *validate {
		if len(invalidTodos) > 0 {
			printInvalid(invalidTodos, absPath)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if *aggregated {
		printAggregated(todosByFile, absPath)
	} else {
		printStandard(todosByFile, absPath)
	}
}

func findTodos(searchDir string) (map[string][]Todo, []InvalidTodo, error) {
	todosByFile := make(map[string][]Todo)
	var invalidTodos []InvalidTodo

	err := filepath.WalkDir(searchDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			dirName := d.Name()
			if dirName == ".git" || dirName == ".hg" || dirName == "node_modules" || strings.HasPrefix(dirName, "bazel-") || (strings.HasPrefix(dirName, ".") && dirName != ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		marker, supported := lineCommentMarkers[filepath.Ext(path)]
		if !supported {
			return nil
		}

		fileInfo, _ := d.Info()
		modTime := fileInfo.ModTime()
		file, _ := os.Open(path)
		defer file.Close()

		lenientTodoRegex := regexp.MustCompile(fmt.Sprintf(`(?i)^\s*%s.*todo`, regexp.QuoteMeta(marker)))
		validTodoRegex := regexp.MustCompile(fmt.Sprintf(`^\s*%s\s*TODO:\s*(.*)\s*\.`, regexp.QuoteMeta(marker)))

		scanner := bufio.NewScanner(file)
		for i := 1; scanner.Scan(); i++ {
			line := scanner.Text()
			if !lenientTodoRegex.MatchString(line) {
				continue
			}
			if matches := validTodoRegex.FindStringSubmatch(line); len(matches) > 1 {
				todosByFile[path] = append(todosByFile[path], Todo{
					File: path, Line: i, Message: strings.TrimSpace(matches[1]), ModTime: modTime,
				})
			} else {
				reason := "Invalid format."
				if !strings.Contains(line, "TODO:") {
					reason = "Use uppercase 'TODO:'."
				} else if !strings.HasSuffix(strings.TrimSpace(line), ".") {
					reason = "Missing trailing period."
				}
				invalidTodos = append(invalidTodos, InvalidTodo{
					File: path, Line: i, Content: strings.TrimSpace(line), Reason: reason,
				})
			}
		}
		return nil
	})

	return todosByFile, invalidTodos, err
}

func printStandard(todosByFile map[string][]Todo, searchDir string) {
	var allTodos []Todo
	for _, todos := range todosByFile {
		allTodos = append(allTodos, todos...)
	}
	sort.Slice(allTodos, func(i, j int) bool { return allTodos[i].ModTime.After(allTodos[j].ModTime) })
	for _, todo := range allTodos {
		relPath, _ := filepath.Rel(searchDir, todo.File)
		fmt.Printf("%s:%d: %s\n", relPath, todo.Line, todo.Message)
	}
}

func printAggregated(todosByFile map[string][]Todo, searchDir string) {
	counts := make([]FileTodoCount, 0, len(todosByFile))
	for file, todos := range todosByFile {
		counts = append(counts, FileTodoCount{File: file, Count: len(todos)})
	}
	sort.Slice(counts, func(i, j int) bool { return counts[i].Count > counts[j].Count })
	for _, item := range counts {
		relPath, _ := filepath.Rel(searchDir, item.File)
		fmt.Printf("%s: %d TODOs\n", relPath, item.Count)
	}
}

func printInvalid(invalidTodos []InvalidTodo, searchDir string) {
	fmt.Print("Invalid TODOs found:\n")
	for _, todo := range invalidTodos {
		relPath, _ := filepath.Rel(searchDir, todo.File)
		fmt.Printf("%s:%d: [%s] %s\n", relPath, todo.Line, todo.Reason, todo.Content)
	}
}
