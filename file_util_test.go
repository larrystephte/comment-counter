package main

import (
	"comment-counter/parser"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetSourceFiles tests the GetSourceFiles function to ensure it correctly identifies
// and returns the expected number of source files in the given directory.
func TestGetSourceFiles(t *testing.T) {
	// Initialize the list of comment parsers. Currently, only C++ parser is included.
	commentParsers := []parser.ICommentParser{
		parser.NewCppCommentParser(),
	}

	// Call GetSourceFiles to retrieve all source files in the "testing/cpp" directory.
	files, err := parser.GetSourceFiles("testing/cpp", commentParsers)
	if err == nil {
		for _, file := range files {
			fmt.Printf("File: %s\n", file)
		}
	}
	assert.NoError(t, err)

	assert.Equal(t, len(files), 10)
}
