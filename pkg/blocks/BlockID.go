package blocks

import (
	"fmt"
	"regexp"
	"strings"
)

type BlockID struct {
	rawInput, content, ID string
}

// find everything up to the block ID -- simpler than a complex regex for para
const blockIDFinderStr = "[\\S\\s]*?\\^(\\w+?)\\s"
const paraSplit = "\n\n"

var (
	BlockIdFinder = regexp.MustCompile(blockIDFinderStr)
	BlockIdRegex  = regexp.MustCompile("\\^(\\w+?)\\s$")
)

// NewBlockID accepts text that ends in a blockID
func NewBlockID(rawInput string) Content {
	b := BlockID{
		rawInput: rawInput,
	}

	matches := BlockIdRegex.FindStringSubmatch(rawInput)
	if len(matches) < 2 {
		b.content = b.rawInput
		return b
	}

	b.ID = matches[1]

	paras := strings.Split(rawInput, paraSplit)
	lastPara := paras[len(paras)-1]
	lastPara = fmt.Sprintf(`<a name="%s"></a>`, b.ID) + lastPara

	paras[len(paras)-1] = lastPara
	b.content = strings.Join(paras, paraSplit)

	return b
}

func (b BlockID) GetMarkdown() string {
	return b.content
}
