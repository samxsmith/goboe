package blocks

type Content interface {
	GetMarkdown() string
}

type NoteProvider interface {
	RelativePathToNote(noteName string) string
	GetNoteContent(noteName string) (string, error)
	RegisterLink(toNoteName string)
}

type Block struct {
	rawInput string
}

// NewBlock creates a generic block for markdown types that don't need to be converted
func NewBlock(rawInput string) Content {
	return Block{rawInput: rawInput}
}

func (b Block) GetMarkdown() string {
	return b.rawInput
}
