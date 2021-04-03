package goboe

import (
	"fmt"
	"goboe/blocks"
	"strings"
)

type Note struct {
	blocks       []blocks.Content
	name         string
	absolutePath string
	rawContent   string
	content      string
	title        string
	vault        VaultProvider
}

type VaultProvider interface {
	RelativePathFromNoteToNote(fromNote, toNote string) string
	AbsolutePathToNote(noteName string) (string, error)
	GetNoteContent(noteName string) (string, error)
	RegisterLink(fromNote, toNote string)
	GetBacklinksForNote(noteName string) []string
}

func NewNote(name string, vault VaultProvider) *Note {
	rawContent, err := vault.GetNoteContent(name)
	if err != nil {
		panic(err)
	}
	abs, err := vault.AbsolutePathToNote(name)
	if err != nil {
		panic(err)
	}

	n := Note{
		name:         name,
		absolutePath: abs,
		rawContent:   rawContent,
		vault:        vault,
	}

	n.title = strings.ToTitle(n.name)

	return &n
}

func (n *Note) Build() {
	// order is important
	// embeds must come before links
	// otherwise links will remove embeds

	n.content = n.rawContent

	n.content = blocks.BlockIdFinder.ReplaceAllStringFunc(n.content, func(match string) string {
		b := blocks.NewBlockID(match)
		return b.GetMarkdown()
	})

	n.content = blocks.EmbedFinder.ReplaceAllStringFunc(n.content, func(match string) string {
		e := blocks.NewEmbed(match, n)
		return e.GetMarkdown()
	})

	n.content = blocks.LinkNoteFinder.ReplaceAllStringFunc(n.content, func(match string) string {
		l := blocks.NewLink(match, n)
		return l.GetMarkdown()
	})
}

// TODO: create map of backlinks from all notes as we go, using Vault method
// TODO: append backlinks as blocks to notes at end

func (n Note) GetTitle() string {
	return fmt.Sprintf("# %s", n.title)
}

func (n Note) GetMarkdown() string {
	title := n.GetTitle() + "\n\n"
	c := title + n.content + "\n\n\n# Backlinks\n"

	backlinkedNotes := n.vault.GetBacklinksForNote(n.name)
	for _, bl := range backlinkedNotes {
		wikiLink := fmt.Sprintf("[[%s]]", bl)
		l := blocks.NewLink(wikiLink, n)
		c += l.GetMarkdown() + "\n\n"
	}

	return c
}

func (n Note) RelativePathToNote(toNote string) string {
	return n.vault.RelativePathFromNoteToNote(n.name, toNote)
}

func (n Note) GetNoteContent(noteName string) (string, error) {
	return n.vault.GetNoteContent(noteName)
}

func (n Note) RegisterLink(toNoteName string) {
	n.vault.RegisterLink(n.name, toNoteName)
}
