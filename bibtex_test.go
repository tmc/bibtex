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
	var lastToken lexeme
	for lexeme := range ls {
		fmt.Println("Got lexeme!", lexeme)
		lastToken = lexeme
	}
	if lastToken.typ != tokenEOF {
		t.Error("expected EOF:", lastToken)
	}

}
