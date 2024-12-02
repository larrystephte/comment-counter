package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// GetSourceFiles walks through all directories starting from the root directory,
// and collects all source files that match the extensions provided by the comment parsers.
// It uses a wait group and channels to handle concurrent directory traversal and file collection.
func GetSourceFiles(root string, commentParser []ICommentParser) ([]string, error) {
	var files []string
	var wg sync.WaitGroup
	filesChan := make(chan string, 1000) // Buffered channel to store file paths

	//Goroutine to close the channel once all directory walking is done
	go func() {
		wg.Wait()
		close(filesChan)
	}()

	// Start the directory walking process
	wg.Add(1)
	go walkDirectory(root, commentParser, &wg, filesChan)

	// Collect file paths from the channel
	for file := range filesChan {
		files = append(files, file)
	}

	return files, nil
}

// walkDirectory recursively walks through the given directory and its subdirectories,
// sending file paths that match the source file extensions to the filesChan channel.
func walkDirectory(path string, commentParser []ICommentParser, wg *sync.WaitGroup, filesChan chan string) {
	defer wg.Done()

	entries, err := os.ReadDir(path) // Read the directory entries
	if err != nil {
		return // If there's an error, return immediately
	}
	for _, entry := range entries {
		if entry.IsDir() {
			// If the entry is a directory, walk into it recursively
			wg.Add(1)
			go walkDirectory(filepath.Join(path, entry.Name()), commentParser, wg, filesChan)
		} else if isSourceFile(entry.Name(), commentParser) {
			// If the entry is a source file, send its path to the filesChan channel
			filesChan <- filepath.Join(path, entry.Name())
		}
	}
}

// isSourceFile checks if a given file name matches any of the source file extensions
// provided by the comment parsers.
func isSourceFile(filename string, commentParser []ICommentParser) bool {
	for _, parser := range commentParser {
		for _, ext := range parser.GetLanguage().GetExtensions() {
			if strings.HasSuffix(filename, ext) {
				return true
			}
		}
	}
	return false
}

// GetParserForFile returns the appropriate comment parser for the given file based on its extension.
func GetParserForFile(filename string, commentParser []ICommentParser) (ICommentParser, error) {
	ext := filepath.Ext(filename)
	for _, parser := range commentParser {
		for _, langExt := range parser.GetLanguage().GetExtensions() {
			if ext == langExt {
				return parser.NewCommentParser(), nil
			}
		}
	}
	return nil, fmt.Errorf("unsupported file type: %s", filename)
}
