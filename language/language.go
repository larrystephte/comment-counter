package language

/*
LanguageConfig defines an interface for language-specific configurations.
This interface allows for easy extension to support different programming languages
by implementing language-specific details while sharing common functionality.
*/
type LanguageConfig interface {
	//GetExtensions returns a slice of file extensions associated with the language.
	// For example: [".cpp", ".hpp"] for C++, [".go"] for Go, [".py"] for Python.
	GetExtensions() []string

	//GetCommentSymbols returns a CommentSymbols struct containing the comment
	//symbols used by the language for both inline and block comments.
	GetCommentSymbols() CommentSymbols
}

/*
CommentSymbols represents the comment symbols used in a programming language.
It includes symbols for both inline (single-line) and block (multi-line) comments.
*/
type CommentSymbols struct {
	//Inline represents the symbol(s) used for single-line comments.
	//For example: "//" in C-like languages, "#" in Python.
	Inline string

	//Block is an array of two strings representing the opening and closing
	//symbols for multi-line comments. For example: ["/*", "*/"] in C-like languages.
	Block [2]string
}
