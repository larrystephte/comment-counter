# Comment Line Counter for C/C++ Source Code

## Usage

Run the tool with:
```shell
go run . testing/cpp
```

## Testing
Write Go test cases under the ./ directory. Use example_test.go as a reference, and test your implementation with:

```shell
go test ./...
```

##  Description

A lightweight command-line tool written in Go to count comment lines in C/C++ source code. It recursively processes files (*.c, *.cpp, *.h, *.hpp) in a given directory, categorizing comments into inline (//) and block (/* … */) types, and outputs a detailed summary. The tool handles special cases, excludes EOF trailing lines, and provides alphabetically sorted results for easy readability. It includes a testing framework for robust validation and can be extended to support other programming languages.

