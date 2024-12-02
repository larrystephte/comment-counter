package language

// PythonConfig implements the LanguageConfig interface for python languages.
//
//	It defines the file extensions and comment symbols used in python files.
type PythonConfig struct {
}

func (PythonConfig) GetExtensions() []string {
	return []string{".py"}
}

func (PythonConfig) GetCommentSymbols() CommentSymbols {
	return CommentSymbols{
		Inline: "#",
		Block:  [2]string{`'''`, `'''`},
	}
}
