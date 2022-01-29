package linkmanagement

import (
	"path/filepath"
)

type LinkRegistry struct {
	backlinks map[string][]Link
	links     map[string][]Link
}

type Link struct {
	Name, PathFromRoot string

	// this is a link relative to the current note, to the desired note
	linkToNote string
}

func (l Link) GetLinkName() string {
	return l.Name
}

func (l Link) GetLinkPath() string {
	return l.linkToNote
}

func NewLinkRegistry() LinkRegistry {
	return LinkRegistry{
		backlinks: map[string][]Link{},
		links:     map[string][]Link{},
	}
}

func (lR *LinkRegistry) RegisterLinks(linking Link, linked []Link) {
	for _, linked := range linked {
		// backlink should go from linked back to linking
		linking.linkToNote = fromNoteToNote(linked, linking)
		lR.backlinks[linked.Name] = append(lR.backlinks[linked.Name], linking)

		// link goes from linking to linked
		linked.linkToNote = fromNoteToNote(linking, linked)
		lR.links[linking.Name] = append(lR.links[linking.Name], linked)
	}
}

func fromNoteToNote(fromNote, toNote Link) string {
	fromNoteDir := filepath.Dir(fromNote.PathFromRoot)
	linkToDir, _ := filepath.Rel(fromNoteDir, toNote.PathFromRoot)
	return linkToDir
}

func (lR *LinkRegistry) GetBacklinksForNote(name string) []Link {
	return lR.backlinks[name]
}
func (lR *LinkRegistry) GetLinksForNote(name string) []Link {
	return lR.links[name]
}
