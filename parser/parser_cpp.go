package parser

import (
	"bufio"
	"comment-counter/language"
	"os"
	"regexp"
	"strings"
)

// CppCommentParser represents a parser for C++ comments
// It keeps track of various states to correctly identify and count different types of comments
type CppCommentParser struct {
	CommentParser
	inBlockComment      bool           // Indicates if we're currently inside a block comment
	inStringLiteral     bool           // inside a string literal
	inRawStringLiteral  bool           // inside a raw string literal
	rawStringDelimiter  string         // Stores the delimiter for the current raw string
	previousLineInline  bool           // Indicates if the previous line ended with an inline comment continuation
	lastBlockCommentEnd int            // Tracks the end position of the last block comment on the current line
	reRawStringStart    *regexp.Regexp //Regex to identify the start of a raw string literal
}

func NewCppCommentParser() *CppCommentParser {
	return &CppCommentParser{
		CommentParser:    CommentParser{language: language.CppConfig{}},
		reRawStringStart: regexp.MustCompile(`R"([^()]*)\(`),
	}
}

func (parser *CppCommentParser) GetLanguage() language.LanguageConfig {
	return parser.language
}

func (parser *CppCommentParser) NewCommentParser() ICommentParser {
	return &CppCommentParser{
		CommentParser:    CommentParser{language: language.CppConfig{}},
		reRawStringStart: regexp.MustCompile(`R"([^()]*)\(`),
	}
}

// Parse processes the file and returns comment statistics
// It reads the file line by line, calling parseLine for each line
func (parser *CppCommentParser) Parse(filename string) (CommentStats, error) {
	file, err := os.Open(filename)
	if err != nil {
		return CommentStats{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	stats := CommentStats{Filename: filename}
	for scanner.Scan() {
		stats.Total++
		line := scanner.Text()
		inlineCount, blockCount := parser.parseLine(line)
		stats.Inline += inlineCount
		stats.Block += blockCount
	}
	if err := scanner.Err(); err != nil {
		return CommentStats{}, err
	}

	return CommentStats{
		Filename: filename,
		Total:    stats.Total,
		Inline:   stats.Inline,
		Block:    stats.Block,
	}, nil
}

// parseLine processes a single line of code
func (parser *CppCommentParser) parseLine(line string) (inlineCount, blockCount int) {
	parser.lastBlockCommentEnd = -1 //Reset for new line
	trimmedLine := strings.TrimSpace(line)

	if len(trimmedLine) == 0 {
		return parser.handleEmptyLine()
	}

	return parser.parseNonEmptyLine(trimmedLine)
}

// handleEmptyLine deals with empty lines
// It only counts as a block comment line if we're currently inside a block comment
func (parser *CppCommentParser) handleEmptyLine() (inlineCount, blockCount int) {
	if parser.inBlockComment {
		return 0, 1
	}
	return 0, 0
}

// It iterates through the line character by character, identifying and counting different types of comments
func (parser *CppCommentParser) parseNonEmptyLine(line string) (inlineCount, blockCount int) {
	for i := 0; i < len(line); {

		//Handle continued inline comment from previous line
		if parser.previousLineInline {
			inlineCount++
			parser.previousLineInline = false
			break
		}

		var nextIndex int
		var inlineIncrement, blockIncrement int

		//Process different states (block comment, string literal, etc.)
		switch {
		case parser.inBlockComment:
			nextIndex, blockIncrement = parser.handleBlockComment(line, i)
		case parser.inStringLiteral:
			nextIndex = parser.handleStringLiteral(line, i)
		case parser.inRawStringLiteral:
			nextIndex = parser.handleRawStringLiteral(line, i)
		default:
			nextIndex, inlineIncrement, blockIncrement = parser.handleDefaultCase(line, i)
		}

		inlineCount += inlineIncrement

		//Only increment block count if this is the first block comment on the line
		// This prevents counting adjacent block comments as separate
		if blockIncrement > 0 && parser.lastBlockCommentEnd == -1 {
			blockCount += blockIncrement
			parser.lastBlockCommentEnd = nextIndex
		}

		if nextIndex == -1 {
			break
		}

		i = nextIndex
	}

	return inlineCount, blockCount
}

// It searches for the end of the current block comment and updates the parser state accordingly
func (parser *CppCommentParser) handleBlockComment(line string, i int) (int, int) {
	if idx := strings.Index(line[i:], parser.language.GetCommentSymbols().Block[1]); idx != -1 {
		parser.inBlockComment = false
		return i + idx + len(parser.language.GetCommentSymbols().Block[1]), 1
	}
	return -1, 1
}

// It handles escaped quotes and the end of string literals
func (parser *CppCommentParser) handleStringLiteral(line string, i int) int {
	if line[i] == '\\' {
		return i + 2 // Skip escaped character
	}
	if line[i] == '"' {
		parser.inStringLiteral = false
		return i + 1
	}
	return i + 1
}

// It searches for the end of the raw string based on its unique delimiter
func (parser *CppCommentParser) handleRawStringLiteral(line string, i int) int {
	endRawString := ")" + parser.rawStringDelimiter + `"`
	if idx := strings.Index(line[i:], endRawString); idx != -1 {
		parser.inRawStringLiteral = false
		return i + idx + len(endRawString)
	}
	return -1 // Raw string continues to next line
}

// handleDefaultCase processes the default case when not in any special state
// It checks for the start of raw string literals, string literals, and comments
func (parser *CppCommentParser) handleDefaultCase(line string, i int) (int, int, int) {
	// Check for raw string literal start
	if match := parser.reRawStringStart.FindStringSubmatchIndex(line[i:]); match != nil {
		return parser.handleRawStringStart(line, i, match), 0, 0
	}

	// Check for string literal start
	if line[i] == '"' && !parser.isEscapedQuote(line, i) {
		parser.inStringLiteral = true
		return i + 1, 0, 0
	}

	// Check for inline comment
	if strings.HasPrefix(line[i:], parser.language.GetCommentSymbols().Inline) {
		return parser.handleInlineComment(line), 1, 0
	}

	// Check for block comment start
	if strings.HasPrefix(line[i:], parser.language.GetCommentSymbols().Block[0]) {
		return parser.handleBlockCommentStart(line, i), 0, 1
	}

	return i + 1, 0, 0 // Move to next character
}

// It sets up the parser state for a new raw string literal
func (parser *CppCommentParser) handleRawStringStart(line string, i int, match []int) int {
	if i+match[3]+1 < len(line) && line[i+match[3]+1] == ')' {
		return i + 1 // Not a valid raw string start, move to next character
	}
	parser.rawStringDelimiter = line[i+match[2] : i+match[3]]
	parser.inRawStringLiteral = true
	return i + match[1]
}

// isEscapedQuote checks if a quote is escaped (i.e., part of a character literal)
func (parser *CppCommentParser) isEscapedQuote(line string, i int) bool {
	return i > 0 && line[i-1] == '\'' && i+1 < len(line) && line[i+1] == '\''
}

// It checks if the comment continues to the next line
func (parser *CppCommentParser) handleInlineComment(line string) int {
	if strings.HasSuffix(line, "\\") {
		parser.previousLineInline = true
	}
	return -1
}

// It checks if the block comment ends on the same line or continues to the next
func (parser *CppCommentParser) handleBlockCommentStart(line string, i int) int {
	i += len(parser.language.GetCommentSymbols().Block[0])
	if idx := strings.Index(line[i:], parser.language.GetCommentSymbols().Block[1]); idx != -1 {
		return i + idx + len(parser.language.GetCommentSymbols().Block[1])
	}
	parser.inBlockComment = true
	return -1
}
