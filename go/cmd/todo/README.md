# TODO CLI

This is a command-line tool for finding and validating TODO comments in source code.

## User Overview

### Installation

To build the tool, run the following command from the root of the repository:

```bash
bazel build //go/cmd/todo:todo
```

This will create a binary at `bazel-bin/go/cmd/todo/todo_/todo`.

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

- Find all TODOs in the current directory and its subdirectories:

```bash
./bazel-bin/go/cmd/todo/todo_/todo
```

- Find all TODOs in the `go/` directory and display them in an aggregated view:

```bash
./bazel-bin/go/cmd/todo/todo_/todo --aggregated go/
```

- Validate all TODOs in the current directory:

```bash
./bazel-bin/go/cmd/todo/todo_/todo --validate
```

- Validate all TODOs in the current directory and show the result in an aggregated view if the validation passes:

```bash
./bazel-bin/go/cmd/todo/todo_/todo --validate --aggregated
```

- Validate all TODOs in the current directory and suppress the output of invalid TODOs:

```bash
./bazel-bin/go/cmd/todo/todo_/todo --validate --quiet
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