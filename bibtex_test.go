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
	l := lexBibTeX(simpleBibTeXDocument)
	tokens := make([]lexeme, 0)
	for token := l.nextToken(); ; token = l.nextToken() {
            tokens = append(tokens, token)
			fmt.Println(token)
			if token.typ == tokenEOF || token.typ == tokenError {
			 break
			}
	}
	fmt.Println(tokens)
	if tokens[len(tokens)-1].typ != tokenEOF {
		t.Error("expected EOF, got:", tokens[len(tokens)-1])
	}

}
