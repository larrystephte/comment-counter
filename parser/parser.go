package parser

import (
	"comment-counter/language"
	"fmt"
	"sync"
)

// CommentStats represents statistics about comments in a source code file.
type CommentStats struct {
	Filename string
	Total    int
	Inline   int
	Block    int
}

// CommentParser defines an interface for parsing comments in a source code file.
type ICommentParser interface {
	Parse(filename string) (CommentStats, error)
	GetLanguage() language.LanguageConfig
	NewCommentParser() ICommentParser
}

// CommentParser is a base struct implementing the ICommentParser interface.
type CommentParser struct {
	language language.LanguageConfig
}

// CountCommentLines walks through the given directory, processes all source files using the provided parsers,
// and collects statistics about the comments in each file.
func CountCommentLines(dir string, commentParsers []ICommentParser) ([]CommentStats, error) {
	files, err := GetSourceFiles(dir, commentParsers) // Get all source files in the directory
	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
		return nil, err
	}

	var wg sync.WaitGroup
	statsChan := make(chan CommentStats, len(files)) // Channel to collect comment statistics

	// Process each file concurrently
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			commentParser, err := GetParserForFile(file, commentParsers) // Get the appropriate parser for the file
			if err != nil {
				fmt.Printf("Error getting parser for file %s: %v\n", file, err)
				return
			}
			fs, err := commentParser.Parse(file) // Parse the file and collect statistics
			if err != nil {
				fmt.Printf("Error analyzing file %s: %v\n", file, err)
				return
			}
			statsChan <- fs // Send the statistics to the channel
		}(file)
	}

	// Close the stats channel once all goroutines are done
	go func() {
		wg.Wait()
		close(statsChan)
	}()

	var stats []CommentStats

	// Collect all statistics from the channel
	for fs := range statsChan {
		stats = append(stats, fs)
	}

	return stats, nil
}
