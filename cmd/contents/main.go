package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/samxsmith/goboe"
	flag "github.com/spf13/pflag"
)

var (
	templateFlag      = flag.StringP("template", "t", "", "(optional) You can specify a template to wrap your content")
	frontMatterFilter = flag.StringP("front-matter-filter", "f", "", "(optional) Specify a front matter filter. If present, only notes with this front matter key will be included.")
)

type opts struct {
	templatePath            string
	frontMatterFilter, root string
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
		return f, fmt.Errorf("You need to specify the path to your Goboe-built vault. \n\t e.g. goboe ~/Documents/my_vault")
	}

	root, err := goboe.PathToAbs(root)
	if err != nil {
		return f, fmt.Errorf("failed to find your vault: %w", err)
	}

	return opts{
		templatePath:      *templateFlag,
		frontMatterFilter: *frontMatterFilter,
		root:              root,
	}, nil
}

const contentPlaceholder = "{content}"

var templater = func(noteBodyHtml []byte) []byte {
	return noteBodyHtml
}

type fileTree struct {
	fullPath string
	subTrees map[string]*fileTree
	files    []string
}

func (ft *fileTree) Add(pathParts []string, name string) {

	if len(pathParts) == 0 || pathParts[0] == "." {
		ft.files = append(ft.files, name)
		return
	}

	if subtree, ok := ft.subTrees[pathParts[0]]; ok {
		subtree.Add(pathParts[1:], name)
		return
	}

	subtree := &fileTree{
		fullPath: filepath.Join(ft.fullPath, pathParts[0]),
		subTrees: map[string]*fileTree{},
	}

	subtree.Add(pathParts[1:], name)
	ft.subTrees[pathParts[0]] = subtree
	return
}

func run(o opts) error {
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

	err := filepath.WalkDir(o.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		relPathToNote, err := filepath.Rel(o.root, path)
		if err != nil {
			return err
		}
		dir := filepath.Dir(relPathToNote)
		p := strings.Split(dir, "/")
		baseTree.Add(p, filepath.Base(relPathToNote))

		return nil
	})

	if err != nil {
		panic(err)
	}

	buildIndexFile(o.root, baseTree)
	return nil
}

func buildIndexFile(outputBasePath string, t fileTree) {
	var subDirLinks []string
	for _, ds := range t.subTrees {
		relPath, _ := filepath.Rel(t.fullPath, ds.fullPath)
		subIndexPath := filepath.Join(relPath, "index.html")
		subDirLinks = append(subDirLinks, fmt.Sprintf(`<a href="%s">&#128193; %s</a>`, subIndexPath, ds.fullPath))
		buildIndexFile(outputBasePath, *ds)
	}

	subDirB := []byte(strings.Join(subDirLinks, "<br>"))

	noteLines := make([]string, len(t.files))
	for i, note := range t.files {
		noteLines[i] = fmt.Sprintf(`<a href="%s">%s</a>`, note, note)
	}

	noteB := []byte(strings.Join(noteLines, "<br>"))

	output := append(subDirB, []byte("<br><br>")...)
	output = append(output, noteB...)

	outputPath := filepath.Join(outputBasePath, t.fullPath, "index.html")

	fmt.Println("writing ", outputPath)

	os.MkdirAll(filepath.Dir(outputPath), 0700)

	if err := ioutil.WriteFile(outputPath, output, 0700); err != nil {
		panic(err)
	}
}
