package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/samxsmith/goboe"

	"github.com/gomarkdown/markdown"
	flag "github.com/spf13/pflag"
)

var (
	outputDirFlag = flag.StringP("output-dir", "o", "", "(required) Where would you like the us to store the output files?")
	templateFlag  = flag.StringP("template", "t", "", "(optional) You can specify a template to wrap your content")
)

func main() {
	flag.Parse()

	root := flag.Arg(0)
	if root == "" {
		println("You need to specify the path to your Obsidian vault. \n\t e.g. goboe ~/Documents/my_vault")
		return
	}

	root, err := filepath.Abs(root)
	if err != nil {
		fmt.Println("failed to find your vault:", err)
		return
	}

	if outputDirFlag == nil || *outputDirFlag == "" {
		println("Missing required flag: ", "-o")
		flag.Usage()
		return
	}

	fmt.Println("Using Obsidian Vault: ", root)
	fmt.Println("Will output to : ", *outputDirFlag)

	var templateParts [2][]byte
	var useTemplate bool
	if templateFlag != nil && *templateFlag != "" {
		templateB, err := ioutil.ReadFile(*templateFlag)
		if err != nil {
			println("Unable to read template: " + err.Error())
			return
		}
		templateS := string(templateB)
		parts := strings.Split(templateS, "{content}")
		if len(parts) != 2 {
			println("Template must include `{content}`, and only once")
			return
		}
		templateParts[0], templateParts[1] = []byte(parts[0]), []byte(parts[1])
		useTemplate = true
	}

	if strings.HasPrefix(root, "~") {
		home, ok := os.LookupEnv("HOME")
		if ok {
			root = strings.Replace(root, "~", home, 1)
		}
	}

	v := goboe.LoadVault(root)
	v.SantiseMarkdown()

	for {
		content, name, ok := v.NextNote()
		if !ok {
			break
		}

		markdownPath, err := v.AbsolutePathToNote(name)
		if err != nil {
			panic(err)
		}
		markdownDirPath := filepath.Dir(markdownPath)
		pathFromRoot, err := filepath.Rel(root, markdownDirPath)
		if err != nil {
			panic(err)
		}

		outputDir := filepath.Join(*outputDirFlag, pathFromRoot)
		outputPath := filepath.Join(outputDir, name+".html")

		md := content.GetMarkdown()

		html := markdown.ToHTML([]byte(md), nil, nil)
		if useTemplate {
			html = append(templateParts[0], html...)
			html = append(html, templateParts[1]...)
		}

		os.MkdirAll(outputDir, 0700)
		if err = ioutil.WriteFile(outputPath, html, 0700); err != nil {
			panic(err)
		}
	}

	fmt.Println("Goboe complete.")
}
