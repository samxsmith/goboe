package blocks

import (
	"fmt"
	"regexp"
	"strings"
)

type EmbedNote struct {
	rawInput          string
	content           string
	embedNoteName     string
	embedNoteBlockRef string
	embedNotePath     string
	embedNoteContent  string
}

var EmbedFinder = regexp.MustCompile("\\!\\[\\[(.+?)\\]\\]")
var EmbedNameRegex = regexp.MustCompile("^\\!\\[\\[(.+?)\\]\\]$")

func NewEmbed(rawInput string, note NoteProvider) Content {
	e := EmbedNote{
		rawInput: rawInput,
	}

	matches := EmbedNameRegex.FindStringSubmatch(rawInput)
	if len(matches) < 2 {
		e.content = rawInput
		return e
	}
	e.embedNoteName, e.embedNoteBlockRef = getLinkedNoteParts(matches[1])

	embedNoteContent, err := note.GetNoteContent(e.embedNoteName)
	if err != nil {
		fmt.Printf("embedded note <%s> does not exist\n", e.embedNoteName)
		e.content = e.rawInput
		return e
	}

	if e.embedNoteBlockRef != "" {
		noteBlocks := strings.Split(embedNoteContent, paraSplit)
		ref := "^" + e.embedNoteBlockRef
		for _, b := range noteBlocks {
			if strings.Contains(b, ref) {
				embedNoteContent = b
				break
			}
		}
	}

	e.embedNoteContent = embedNoteContent
	e.content = fmt.Sprintf("---\n<span class=\"embed\">%s</span>\n\n---", e.embedNoteContent)

	return e
}

func (e EmbedNote) GetMarkdown() string {
	return e.content
}
