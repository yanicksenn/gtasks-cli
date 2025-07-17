# TODO CLI

This is a command-line tool for finding and validating TODO comments in source code.

## User Overview

### Installation

To build the tool, run the following command from the root of the repository:

```bash
bazel build //go/cmd/todo:todo
```

This will create a binary at `bazel-bin/go/cmd/todo/todo_/todo`.

### Valid TODO Format

The tool validates TODOs based on a specific format. A TODO is considered valid if it follows these rules:

1.  It must start with the language's line comment marker (e.g., `//` for Go, `#` for Python).
2.  It must contain the uppercase `TODO:` keyword, followed by a space.
3.  The TODO message must end with a period (`.`).

**Example of a valid TODO in Go:**
```go
// TODO: This is a valid todo.
```

**Examples of invalid TODOs:**
```go
// todo: This is invalid because of the lowercase 'todo'.
// TODO: This is invalid because it is missing a period
```

### Usage

To run the tool, you can use the following command:

```bash
./bazel-bin/go/cmd/todo/todo_/todo [command] [flags] [directory]
```

If no directory is specified, it will search the current directory.

#### Commands

- `help`: Display the help message.

#### Flags

- `--aggregated`: Group TODOs by file and display the count for each file. Can be combined with `--validate`.
- `--validate`: Validate TODO format and exit with an error if invalid TODOs are found.
- `--quiet`: Suppress output of invalid TODOs when validating. Only works with `--validate`.

#### Examples

- Find all TODOs in the `go/cmd/todo/testdata/simple` directory and its subdirectories:

```bash
./bazel-bin/go/cmd/todo/todo_/todo go/cmd/todo/testdata/simple
```
Output:
```
go/cmd/todo/testdata/simple/test.go:3: This is a valid todo.
```

- Find all TODOs in the `go/cmd/todo/testdata/aggregated` directory and display them in an aggregated view:

```bash
./bazel-bin/go/cmd/todo/todo_/todo --aggregated go/cmd/todo/testdata/aggregated
```
Output:
```
go/cmd/todo/testdata/aggregated/test.go: 2 TODOs
go/cmd/todo/testdata/aggregated/test.py: 1 TODOs
```

- Validate all TODOs in the `go/cmd/todo/testdata/validation-invalid` directory:

```bash
./bazel-bin/go/cmd/todo/todo_/todo --validate go/cmd/todo/testdata/validation-invalid
```
Output:
```
Invalid TODOs found:
go/cmd/todo/testdata/validation-invalid/test.go:2: [Use uppercase 'TODO:'.] // todo: this is invalid because it is lowercase.
go/cmd/todo/testdata/validation-invalid/test.go:3: [Missing trailing period.] // TODO: This is invalid because it is missing a period
```

- Validate all TODOs in the `go/cmd/todo/testdata/validation-aggregated-valid-only` directory and show the result in an aggregated view if the validation passes:

```bash
./bazel-bin/go/cmd/todo/todo_/todo --validate --aggregated go/cmd/todo/testdata/validation-aggregated-valid-only
```
Output:
```
go/cmd/todo/testdata/validation-aggregated-valid-only/test.go: 2 TODOs
```

- Validate all TODOs in the `go/cmd/todo/testdata/validation-quiet-invalid` directory and suppress the output of invalid TODOs:

```bash
./bazel-bin/go/cmd/todo/todo_/todo --validate --quiet go/cmd/todo/testdata/validation-quiet-invalid
```
Output:
```
(no output, exit code 1)
```

## Engineering Overview

### Code Structure

- `main.go`: The main entry point for the application. It contains the logic for parsing command-line arguments, walking the file system, finding and validating TODOs, and printing the output.
- `main_test.go`: The tests for the application. It uses a `testdata` directory to store test cases.
- `testdata/`: Contains test cases for the different features of the application. Each test case is a directory with test files and an `expected.out` file.

### Contributing

To contribute to the project, you can follow these steps:

1.  Make your changes to the code.
2.  Add a new test case to the `testdata` directory if you are adding a new feature.
3.  Run the tests to make sure everything is working correctly:

```bash
bazel test //go/cmd/todo:todo_test
```

4.  Commit your changes and create a pull request.
