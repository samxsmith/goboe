package goboe

import "testing"

func TestFrontMatter(t *testing.T) {
	contentWithFrontMatter :=
		`---
key1: val1
key2: val2
---
body of
notes
`

	contentWithoutFrontMatter := `text in a note`

	contentWithDashes := `notes
with
---
dashes`

	contentWithFrontMatterAndDashes := `---
keyx: valy
---
plus
---
dashes`

	t.Run("Should extract front matter", func(t *testing.T) {
		content, fm := extractFrontMatter(contentWithFrontMatter)
		if content != "body of\nnotes\n" {
			t.Fatalf("content != \"body of \\nnotes\"")
		}

		if fm["key1"] != "val1" {
			t.Fatalf("fm[\"key1\"] != \"val1\"")
		}
		if fm["key2"] != "val2" {
			t.Fatalf("fm[\"key2\"] != \"val2\"")
		}
	})

	t.Run("Should work without front matter", func(t *testing.T) {
		content, fm := extractFrontMatter(contentWithoutFrontMatter)
		if content != "text in a note" {
			t.Fatalf("content != \"text in a note\"")
		}

		if len(fm) != 0 {
			t.Fatalf("len(fm) != 0")
		}
	})

	t.Run("Should work with other dashes", func(t *testing.T) {
		content, fm := extractFrontMatter(contentWithDashes)
		if content != "notes\nwith\n---\ndashes" {
			t.Fatalf("content != \"notes\\nwith\\n---\\ndashes\"")
		}

		if len(fm) != 0 {
			t.Fatalf("len(fm) != 0")
		}
	})

	t.Run("Should work with other dashes", func(t *testing.T) {
		content, fm := extractFrontMatter(contentWithFrontMatterAndDashes)
		if content != "plus\n---\ndashes" {
			t.Fatalf("plus\n---\ndashes")
		}

		if fm["keyx"] != "valy" {
			t.Fatalf("fm[\"keyx\"] != \"valy\"")
		}
	})
}
