package goboe

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/gomarkdown/markdown"
	"github.com/samxsmith/goboe/pkg/blocks"
	"github.com/samxsmith/goboe/pkg/linkmanagement"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Note struct {
	name            string
	filepath        string
	body            string
	frontMatter     map[string]string
	linkedNoteNames []string
}

func OpenNote(path string) (Note, error) {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	noteName := strings.TrimSuffix(filename, ext)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return Note{}, fmt.Errorf("ioutil.ReadFile: %w", err)
	}
	body, frontMatter := extractFrontMatter(string(b))

	allLinks := blocks.CombForLinks(body)
	linksSeen := map[string]bool{}
	var uniqLinks []string
	for _, link := range allLinks {
		if _, ok := linksSeen[link]; !ok {
			linksSeen[link] = true
			uniqLinks = append(uniqLinks, link)
		}
	}

	return Note{
		name:            noteName,
		filepath:        path,
		body:            body,
		frontMatter:     frontMatter,
		linkedNoteNames: uniqLinks,
	}, nil
}

func extractFrontMatter(rawContent string) (bodyContent string, frontMatter map[string]string) {
	if !strings.HasPrefix(rawContent, "---\n") {
		return rawContent, frontMatter
	}
	parts := strings.SplitN(rawContent, "---\n", 3)
	var nonEmptyParts []string
	for _, p := range parts {
		if p == "" {
			continue
		}
		nonEmptyParts = append(nonEmptyParts, p)
	}

	if len(nonEmptyParts) != 2 {
		return rawContent, frontMatter
	}

	if err := yaml.Unmarshal([]byte(nonEmptyParts[0]), &frontMatter); err != nil {
		fmt.Printf("could not pass front matter for file %s: %s\n", rawContent, err)
		return rawContent, frontMatter
	}
	return nonEmptyParts[1], frontMatter
}

func (n Note) Name() string {
	return n.name
}
func (n Note) Path() string {
	return n.filepath
}

func (n Note) Html(lR linkmanagement.LinkRegistry) []byte {
	backlinks := lR.GetBacklinksForNote(n.name)

	links := map[string]linkmanagement.Link{}
	for _, l := range lR.GetLinksForNote(n.name) {
		links[l.Name] = l
	}

	md := markdownifyWikiLinks(n.body, links)

	md = blocks.FormatBullets(md)
	md = addNoteTitle(md, n.name)
	if len(backlinks) == 0 {
		return toHtml(md)
	}

	md += "\n\n## Backlinks\n"
	for _, backlink := range backlinks {
		// TODO: contextual backlinks
		backlinkMd := fmt.Sprintf("- [%s](%s)\n", backlink.GetLinkName(), backlink.GetLinkPath())
		md += backlinkMd
	}

	return toHtml(md)
}

func markdownifyWikiLinks(contentWithWikiLinks string, links map[string]linkmanagement.Link) string {
	return blocks.LinkNoteFinder.ReplaceAllStringFunc(contentWithWikiLinks, func(wikiLink string) string {
		noteName := blocks.GetWikiLinkContent(wikiLink)

		link, ok := links[noteName]
		if !ok || link.GetLinkPath() == "" {
			return blocks.DeadLink(noteName)
		}
		return blocks.MdLink(noteName, link.GetLinkPath())
	})
}

func addNoteTitle(body, title string) string {
	return fmt.Sprintf("# %s \n\n%s", title, body)
}

func toHtml(md string) []byte {
	return markdown.ToHTML([]byte(md), nil, nil)
}
