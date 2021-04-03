package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/samxsmith/goboe"
	flag "github.com/spf13/pflag"
)

const (
	htmlHeader = "<html><head><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"></head><body>"
	htmlFooter = "</body></html>"
)

func main() {
	flag.Parse()

	root := flag.Arg(0)
	if root == "" {
		println("You need to specify the path to your Goboe HTML. \n\t e.g. goboe ~/Documents/my_vault/public")
		return
	}

	root, err := goboe.PathToAbs(root)
	if err != nil {
		fmt.Println("failed to find your vault:", err)
		return
	}

	outputFile := filepath.Join(root, "index.html")

	fmt.Println("Looking for Goboe HTML in: ", root)
	fmt.Println("Will write index.html to: ", outputFile)

	htmlList := ""

	count := 0
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(info.Name()) == ".html" {
			filename := info.Name()
			name := filename[:len(filename)-5]

			pathToDir := filepath.Dir(path)

			relToDir, err := filepath.Rel(root, pathToDir)
			if err != nil {
				panic(err)
			}

			escapedName := url.PathEscape(filename)
			escapedRel := filepath.Join(relToDir, escapedName)
			htmlList += newNoteItem(name, escapedRel)
			count++
		}

		return nil
	})

	output := htmlHeader + htmlList + htmlFooter

	err = ioutil.WriteFile(outputFile, []byte(output), 0700)
	if err != nil {
		fmt.Println("failed to write file:", err)
		return
	}

	fmt.Printf("Index created for %d files \n", count)
}

func newNoteItem(name, link string) string {
	return fmt.Sprintf(`<a href="%s" class="NoteLink">%s</a><br>`, link, name)
}
