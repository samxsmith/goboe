package main

import (
	"bytes"
	"fmt"
	"github.com/samxsmith/goboe"
	flag "github.com/spf13/pflag"
	"io/ioutil"
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

	output := make([]string, len(vault.Notes()))
	for i, note := range vault.Notes() {
		vaultPathToNote := vault.LinkFromVaultRoot(note.Name())
		output[i] = fmt.Sprintf(`<a href="%s">%s</a>`, vaultPathToNote, note.Name())
	}

	fmt.Println("Writing contents page: ", o.outputPath)

	b := []byte(strings.Join(output, "<br>"))
	err = ioutil.WriteFile(o.outputPath, b, 0700)
	if err != nil {
		return fmt.Errorf("unable to write file: %w", err)
	}
	return nil
}
