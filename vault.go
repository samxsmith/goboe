package goboe

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samxsmith/goboe/pkg/linkmanagement"
)

type Vault struct {
	notes        map[string]Note
	linkRegistry linkmanagement.LinkRegistry
	root         string
}

func OpenVault(vaultRoot, frontMatterFilter string) (Vault, error) {
	notes, err := traverseVault(vaultRoot)
	if err != nil {
		return Vault{}, fmt.Errorf("traverseVault: %w", err)
	}

	v := newVault(vaultRoot)

	for _, notePath := range notes {
		note, err := OpenNote(notePath)
		if err != nil {
			return v, fmt.Errorf("OpenNote: %w", err)
		}
		v.notes[note.name] = note
	}

	// do this before generating link registry
	if frontMatterFilter != "" {
		v.ApplyFrontMatterFilter(frontMatterFilter)
	}

	// now we have all notes in a map, and we've applied filters
	// we can do another pass to assemble backlinks
	v.generateLinks()

	return v, nil
}

func newVault(root string) Vault {
	v := Vault{}
	v.root = root
	v.notes = map[string]Note{}
	v.linkRegistry = linkmanagement.NewLinkRegistry()
	return v
}

func traverseVault(vaultRoot string) ([]string, error) {
	var noteList []string

	err := filepath.Walk(vaultRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		filename := info.Name()
		ext := filepath.Ext(filename)
		if ext != ".md" {
			return nil
		}

		noteList = append(noteList, path)
		return nil
	})

	return noteList, err
}

func (v *Vault) ApplyFrontMatterFilter(desiredFrontMatterKey string) {
	for noteName, note := range v.notes {
		if _, ok := note.frontMatter[desiredFrontMatterKey]; !ok {
			delete(v.notes, noteName)
		}
	}
}

func (v *Vault) generateLinks() {
	for _, note := range v.notes {
		linksInThisNote := note.linkedNoteNames
		var links []linkmanagement.Link
		for _, linkedNoteFileName := range linksInThisNote {
			linkedNoteName := strings.TrimSuffix(linkedNoteFileName, filepath.Ext(linkedNoteFileName))

			link := linkmanagement.Link{
				Name:         linkedNoteFileName,
				PathFromRoot: v.LinkFromVaultRoot(linkedNoteName),
			}
			links = append(links, link)
		}

		// not necessary at this point, but may wish to record this link direction in future
		noteLink := linkmanagement.Link{
			Name:         note.name,
			PathFromRoot: v.LinkFromVaultRoot(note.name),
		}
		v.linkRegistry.RegisterLinks(noteLink, links)
	}
}

func (v *Vault) Notes() []Note {
	notes := make([]Note, len(v.notes))
	i := 0
	for _, n := range v.notes {
		notes[i] = n
		i++
	}
	return notes
}

func (v *Vault) GetLinkRegistry() linkmanagement.LinkRegistry {
	return v.linkRegistry
}

func (v *Vault) LinkFromVaultRoot(noteName string) string {
	note, ok := v.notes[noteName]
	if !ok {
		return ""
	}

	noteDir := filepath.Dir(note.Path())
	noteOutputFilename := fmt.Sprintf("%s.html", note.Name())
	noteDirRelativeToVaultRoot, _ := filepath.Rel(v.root, noteDir)
	return filepath.Join(noteDirRelativeToVaultRoot, noteOutputFilename)
}
