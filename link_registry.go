package goboe

type linkRegistry struct {
	backlinks map[string][]link
	links     map[string][]link
}

type link struct {
	name, path string
}

func (l link) GetLinkName() string {
	return l.name
}

func (l link) GetLinkPath() string {
	return l.path
}

func newLinkRegistry() linkRegistry {
	return linkRegistry{
		backlinks: map[string][]link{},
		links:     map[string][]link{},
	}
}

func (lR *linkRegistry) RegisterLinks(fromNote link, toNotes []link) {
	for _, toNote := range toNotes {
		lR.backlinks[toNote.name] = append(lR.backlinks[toNote.name], fromNote)
	}

	lR.links[fromNote.name] = toNotes
}

func (lR *linkRegistry) GetBacklinksForNote(name string) []link {
	return lR.backlinks[name]
}
func (lR *linkRegistry) GetLinksForNote(name string) []link {
	return lR.links[name]
}
