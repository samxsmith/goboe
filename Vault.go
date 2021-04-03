package goboe

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samxsmith/goboe/blocks"
)

type Vault struct {
	root         string
	notes        map[string]*NoteNode
	cursor       int
	noteIterator []string
	backlinks    map[string][]string
}

type NoteNode struct {
	name, path, fileExt string
	content             *Note
}

func LoadVault(vaultRoot string) Vault {
	noteInfo := map[string]*NoteNode{}
	filepath.Walk(vaultRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(info.Name()) == ".md" {
			filename := info.Name()
			ext := filepath.Ext(filename)
			name := filename[:len(filename)-3]

			noteInfo[name] = &NoteNode{
				name:    name,
				path:    path,
				fileExt: ext,
			}
		}

		return nil
	})

	return Vault{
		root:      vaultRoot,
		notes:     noteInfo,
		backlinks: map[string][]string{},
	}
}

func (v *Vault) SantiseMarkdown() {
	for name, note := range v.notes {
		c := NewNote(name, v)
		c.Build()
		note.content = c
		v.noteIterator = append(v.noteIterator, name)
	}
}

func (v *Vault) NextNote() (blocks.Content, string, bool) {
	if len(v.noteIterator) >= v.cursor+1 {
		nextNoteName := v.noteIterator[v.cursor]
		v.cursor++
		note := v.notes[nextNoteName]
		return note.content, note.name, true
	}
	return nil, "", false
}

func (v *Vault) RelativePathFromNoteToNote(fromNote, toNote string) string {
	// get path of each
	fromAbs, err := v.AbsolutePathToNote(fromNote)
	if err != nil {
		panic("cant find abs path for file: " + fromNote)
	}
	toAbs, err := v.AbsolutePathToNote(toNote)
	if err != nil {
		// the linked file may not exist
		// it's a common use case to link to a file then create later
		fmt.Printf("Can't find file <%s> linked from <%s> \n", toNote, fromNote)
		return ""
	}

	// remove file name
	fromAbsDir := filepath.Dir(fromAbs)
	toAbsDir := filepath.Dir(toAbs)

	// get relative path
	rel, err := filepath.Rel(fromAbsDir, toAbsDir)
	if err != nil {
		panic(fmt.Sprintf("cant find relative path from %s to %s", fromAbs, toAbs))
	}

	// add toNote file name
	return filepath.Join(rel, filepath.Base(toAbs))
}

func (v *Vault) AbsolutePathToNote(noteName string) (string, error) {
	n, ok := v.notes[noteName]
	if !ok {
		return "", ErrFileNotFound
	}
	return n.path, nil
}

func (v *Vault) GetNoteContent(noteName string) (string, error) {
	path, err := v.AbsolutePathToNote(noteName)
	if err != nil {
		return "", err
	}
	return ReadFile(path), nil
}

func (v *Vault) RegisterLink(fromNoteName, toNoteName string) {
	v.backlinks[toNoteName] = append(v.backlinks[toNoteName], fromNoteName)
}

func (v *Vault) GetBacklinksForNote(noteName string) []string {
	return v.backlinks[noteName]
}
