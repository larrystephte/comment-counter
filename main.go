package main

import (
	"comment-counter/parser"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		printHelp()
	} else {
		dir := args[0]
		if err := countCommentLines(dir); err != nil {
			fmt.Println(err)
		}
	}
}

func printHelp() {
	fmt.Println("usage: \n\tgo run . <directory>")
}

// countCommentLines walks through the given directory, processes all C++ source files using the provided parsers,
// and prints statistics about the comments in each file.
func countCommentLines(dir string) error {
	// Initialize the list of comment parsers. Currently, only C++ parser is included.
	commentParsers := []parser.ICommentParser{
		parser.NewCppCommentParser(),
	}

	// Count the comment lines in the specified directory using the provided parsers.
	stats, err := parser.CountCommentLines(dir, commentParsers)

	if err != nil {
		return err
	}

	// Print the statistics for each file.
	for _, fs := range stats {
		fmt.Printf("filename: %s\t", fs.Filename)
		fmt.Printf("total: %d\t", fs.Total)
		fmt.Printf("inline: %d\t", fs.Inline)
		fmt.Printf("block: %d\n", fs.Block)
	}

	return nil
	//	return errors.New(fmt.Sprintf(`
	//error:		not implemented.
	//directory:	%s`, dir))
}
