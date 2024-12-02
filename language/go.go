package language

// GoConfig implements the LanguageConfig interface for C and C++ languages.
//
//	It defines the file extensions and comment symbols used in C/C++ files.
type GoConfig struct {
}

func (GoConfig) GetExtensions() []string {
	return []string{".go"}
}

func (GoConfig) GetCommentSymbols() CommentSymbols {
	return CommentSymbols{
		Inline: "//",
		Block:  [2]string{"/*", "*/"},
	}
}
