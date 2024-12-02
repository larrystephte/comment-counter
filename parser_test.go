package main

import (
	"comment-counter/parser"
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

// TestCountCommentLines tests the CountCommentLines function to ensure it correctly processes
// all C++ source files in the given directory and collects statistics about comments.
func TestCountCommentLines(t *testing.T) {
	commentParsers := []parser.ICommentParser{
		parser.NewCppCommentParser(),
	}

	// Count the comment lines in the specified directory using the provided parsers.
	stats, err := parser.CountCommentLines("testing/cpp", commentParsers)
	assert.NoError(t, err)

	for _, fs := range stats {
		fmt.Printf("filename: %s\t", fs.Filename)
		fmt.Printf("total: %d\t", fs.Total)
		fmt.Printf("inline: %d\t", fs.Inline)
		fmt.Printf("block: %d\n", fs.Block)
	}

}

// TestParse tests the Parse function to ensure it correctly parses individual C++ source files
// and collects statistics about comments.
//
// This test uses the CommonParseTest helper function to test multiple C++ source files
// located in different directories. The expected total lines, inline comments, and block comments
// are provided for each file to verify the correctness of the parsing logic.
func TestParse(t *testing.T) {
	CommonParseTest(t, "testing/cpp/lib_json", "json_reader.cpp", 1992, 134, 0)
	CommonParseTest(t, "testing/cpp/lib_json", "json_tool.h", 138, 13, 19)
	CommonParseTest(t, "testing/cpp/lib_json", "json_value.cpp", 1634, 111, 18)
	CommonParseTest(t, "testing/cpp/lib_json", "json_writer.cpp", 1259, 89, 0)
	CommonParseTest(t, "testing/cpp", "special_cases.cpp", 62, 6, 34)
	CommonParseTest(t, "testing/cpp/test_lib_json", "fuzz.cpp", 54, 5, 0)
	CommonParseTest(t, "testing/cpp/test_lib_json", "fuzz.h", 14, 5, 0)
	CommonParseTest(t, "testing/cpp/test_lib_json", "jsontest.cpp", 430, 54, 1)
	CommonParseTest(t, "testing/cpp/test_lib_json", "jsontest.h", 288, 52, 8)
	CommonParseTest(t, "testing/cpp/test_lib_json", "main.cpp", 3971, 182, 0)
}

// CommonParseTest is a helper function to test the Parse function for individual C++ source files.
// It verifies the correctness of the parsing logic by comparing the parsed statistics with the expected values.
func CommonParseTest(t *testing.T, path string, filename string, total int, inline int, block int) {
	cppParser := parser.NewCppCommentParser()

	testFile := filepath.Join(path, filename)

	// Parse the test file and collect statistics.
	stats, err := cppParser.Parse(testFile)

	fmt.Printf("File: %s\nTotal lines: %d\nInline comments: %d\nBlock comments: %d\n", stats.Filename, stats.Total, stats.Inline, stats.Block)

	assert.NoError(t, err)
	assert.Equal(t, testFile, stats.Filename)
	assert.Equal(t, total, stats.Total)
	assert.Equal(t, inline, stats.Inline)
	assert.Equal(t, block, stats.Block)
}
