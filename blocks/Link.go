package blocks

import (
	"fmt"
	"regexp"
	"strings"
)

type LinkNote struct {
	rawInput         string
	content          string
	linkAbsPath      string
	linkNoteName     string
	linkRelativePath string
	blockRef         string
}

const blockRefDivider = "#^"

var (
	LinkNoteFinder    = regexp.MustCompile("\\[\\[(.+?)\\]\\]")
	LinkNoteNameRegex = regexp.MustCompile("^\\[\\[(.+?)\\]\\]$")
)

func NewLink(rawInput string, note NoteProvider) Content {

	l := LinkNote{
		rawInput: rawInput,
	}

	// get name from link
	matches := LinkNoteNameRegex.FindStringSubmatch(rawInput)
	if len(matches) < 2 {
		l.content = rawInput
		return l
	}

	l.linkNoteName, l.blockRef = getLinkedNoteParts(matches[1])

	l.linkRelativePath = note.RelativePathToNote(l.linkNoteName)
	if l.linkRelativePath == "" {
		l.content = l.rawInput
		return l
	}

	// swap the .md extension for .html
	l.linkRelativePath = l.linkRelativePath[:len(l.linkRelativePath)-3] + ".html"

	if l.blockRef != "" {
		l.linkRelativePath += "#" + l.blockRef
		l.linkNoteName += "#" + l.blockRef
	}

	l.content = fmt.Sprintf("[%s](%s)", l.linkNoteName, l.linkRelativePath)

	note.RegisterLink(l.linkNoteName)
	return l
}

func (l LinkNote) GetMarkdown() string {
	return l.content
}

func getLinkedNoteParts(linkedNoteRef string) (string, string) {
	linkParts := strings.Split(linkedNoteRef, blockRefDivider)
	if len(linkParts) > 1 {
		return linkParts[0], linkParts[1]
	}

	return linkParts[0], ""
}
