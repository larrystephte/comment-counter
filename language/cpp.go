package language

// CppConfig implements the LanguageConfig interface for C and C++ languages.
//
//	It defines the file extensions and comment symbols used in C/C++ files.
type CppConfig struct {
}

func (CppConfig) GetExtensions() []string {
	return []string{".c", ".cpp", ".h", ".hpp"}
}

func (CppConfig) GetCommentSymbols() CommentSymbols {
	return CommentSymbols{
		Inline: "//",
		Block:  [2]string{"/*", "*/"},
	}
}
