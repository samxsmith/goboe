package blocks

import (
	"testing"
)

func TestFormatBullets(t *testing.T) {
	t.Run("Should add new line before bullets", func(t *testing.T) {
		input := `some text
then
- bullet
- bullet2
continue with text`

		expectedOutput := `some text
then

- bullet
- bullet2
continue with text`

		output := FormatBullets(input)
		if output != expectedOutput {
			t.Fatalf("Did not add new line correctly before bullets")
		}
	})

	t.Run("Should make no change when new line already present", func(t *testing.T) {
		input := `some text
then

- bullet
- bullet2
continue with text`

		output := FormatBullets(input)
		if output != input {
			t.Fatalf("Should not have modified input")
		}
	})
}
