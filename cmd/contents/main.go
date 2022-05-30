package main

import (
	"bytes"
	"fmt"
	"github.com/samxsmith/goboe"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var (
	outputFileFlag    = flag.StringP("output-dir", "o", "", "(required) Where would you like the us to store the output files?")
	templateFlag      = flag.StringP("template", "t", "", "(optional) You can specify a template to wrap your content")
	frontMatterFilter = flag.StringP("front-matter-filter", "f", "", "(optional) Specify a front matter filter. If present, only notes with this front matter key will be included.")
)

type opts struct {
	outputPath, templatePath     string
	frontMatterFilter, vaultRoot string
}

func main() {
	flags, err := getFlags()
	if err != nil {
		fmt.Println("ERROR) \n", err)
		return
	}

	if err = run(flags); err != nil {
		fmt.Println("ERROR) \n", err)
		return
	}
}

func getFlags() (f opts, e error) {
	flag.Parse()

	root := flag.Arg(0)
	if root == "" {
		return f, fmt.Errorf("You need to specify the path to your Obsidian vault. \n\t e.g. goboe ~/Documents/my_vault")
	}

	root, err := goboe.PathToAbs(root)
	if err != nil {
		return f, fmt.Errorf("failed to find your vault: %w", err)
	}

	if outputFileFlag == nil || *outputFileFlag == "" {
		println("Missing required flag: ", "-o")
		flag.Usage()
		return
	}

	outputDir, err := goboe.PathToAbs(*outputFileFlag)
	if err != nil {
		return f, fmt.Errorf("failed to parse your output path: %w", err)
	}

	return opts{
		outputPath:        outputDir,
		templatePath:      *templateFlag,
		frontMatterFilter: *frontMatterFilter,
		vaultRoot:         root,
	}, nil
}

const contentPlaceholder = "{content}"

var templater = func(noteBodyHtml []byte) []byte {
	return noteBodyHtml
}

type fileTree struct {
	fullPath string
	subTrees map[string]*fileTree
	notes    []goboe.Note
}

func (ft *fileTree) Add(pathParts []string, n goboe.Note) {

	if len(pathParts) == 0 || pathParts[0] == "." {
		ft.notes = append(ft.notes, n)
		return
	}

	if subtree, ok := ft.subTrees[pathParts[0]]; ok {
		subtree.Add(pathParts[1:], n)
		return
	}

	subtree := &fileTree{
		fullPath: filepath.Join(ft.fullPath, pathParts[0]),
		subTrees: map[string]*fileTree{},
	}

	subtree.Add(pathParts[1:], n)
	ft.subTrees[pathParts[0]] = subtree
	return
}

func run(o opts) error {
	vault, err := goboe.OpenVault(o.vaultRoot, o.frontMatterFilter)
	if err != nil {
		return fmt.Errorf("OpenVault: %w", err)
	}

	if o.templatePath != "" {
		b, err := ioutil.ReadFile(o.templatePath)
		if err != nil {
			return fmt.Errorf("unable to read template file: %w", err)
		}

		template := string(b)
		if !strings.Contains(template, contentPlaceholder) {
			return fmt.Errorf("template file does not have '%s' placeholder", contentPlaceholder)
		}

		contentPlaceholderB := []byte(contentPlaceholder)

		templater = func(noteBodyHtml []byte) []byte {
			return bytes.Replace(b, contentPlaceholderB, noteBodyHtml, 1)
		}
	}

	baseTree := fileTree{
		// empty path at base
		fullPath: "",
		subTrees: map[string]*fileTree{},
	}

	for _, note := range vault.Notes() {
		vaultPathToNote := vault.LinkFromVaultRoot(note.Name())
		dir := filepath.Dir(vaultPathToNote)
		p := strings.Split(dir, "/")
		baseTree.Add(p, note)
	}

	buildIndexFile(o.outputPath, baseTree)

	return nil
}

func buildIndexFile(outputBasePath string, t fileTree) {
	var subDirLinks []string
	for _, ds := range t.subTrees {
		relPath, _ := filepath.Rel(t.fullPath, ds.fullPath)
		subIndexPath := filepath.Join(relPath, "index.html")
		subDirLinks = append(subDirLinks, fmt.Sprintf(`<a href="%s">%s</a>`, subIndexPath, ds.fullPath))
		buildIndexFile(outputBasePath, *ds)
	}

	subDirB := []byte(strings.Join(subDirLinks, "<br>"))

	noteLines := make([]string, len(t.notes))
	for i, note := range t.notes {
		noteLines[i] = fmt.Sprintf(`<a href="%s">%s</a>`, note.Path(), note.Name())
	}

	noteB := []byte(strings.Join(noteLines, "<br>"))

	output := append(subDirB, []byte("<br><br>")...)
	output = append(output, noteB...)

	basePath := filepath.Dir(outputBasePath)
	outputPath := filepath.Join(basePath, t.fullPath, "index.html")

	fmt.Println("writing ", outputPath)

	if err := ioutil.WriteFile(outputPath, output, 0700); err != nil {
		panic(err)
	}
}
