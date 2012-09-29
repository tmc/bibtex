package bibtex

import (
	"fmt"
	"testing"
)

var simpleBibTeXDocument = `@Book{b_id,
 title = "Wonderful story",
}`

var exampleBibTeXDocument = `@Book{hicks2001,
 author    = "von Hicks, III, Michael",
 title     = "Design of a Carbon Fiber Composite Grid Structure for the GLAST
              Spacecraft Using a Novel Manufacturing Technique",
 publisher = "Stanford Press",
 year      =  2001,
 address   = "Palo Alto",
 edition   = "1st",
 isbn      = "0-69-697269-4"
}`

func TestBibTeXLexing(t *testing.T) {
	_, ls := lexBibTeX(simpleBibTeXDocument)
	tokens := make([]lexeme, 0)
	for lexeme := range ls {
		tokens = append(tokens, lexeme)
	}
	expected := `[entryStart "@" identifier "Book" { identifier "b_id" , identifier "title" = string "Wonderful story" , } EOF]`
	actual := fmt.Sprint(tokens)
	if actual != expected {
		t.Errorf("'%s' !+ '%s'", actual, expected)
	}
}
