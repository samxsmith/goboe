package blocks

import "testing"

func TestCombForLinks(t *testing.T) {
	t.Run("should find links", func(t *testing.T) {
		content := `
		here is [[the first]] link and
		then [[there]] [[is a]] second
		`

		links := CombForLinks(content)
		if len(links) != 3 {
			t.Fatalf("len(links) != 3")
		}

		if links[0] != "the first" {
			t.Fatalf("links[0] != \"the first\"")
		}
		if links[1] != "there" {
			t.Fatalf("links[1] != \"there\"")
		}
		if links[2] != "is a" {
			t.Fatalf("links[2] != \"is a\"")
		}
	})
}
