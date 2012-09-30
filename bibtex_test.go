package bibtex

import (
	"fmt"
	"testing"
)

var simpleBibTeXDocument = `@book{b_id,
 title = "Wonderful story",
 year  = 1999,
}`

var exampleBibTeXDocument = `@book{hicks2001,
 author    = "von Hicks, III, Michael",
 title     = "Design of a Carbon Fiber Composite Grid Structure for the GLAST Spacecraft Using a Novel Manufacturing Technique",
 publisher = "Stanford Press",
 year      = 2001,
 address   = "Palo Alto",
 edition   = "1st",
 isbn      = "0-69-697269-4",
}`

func TestBibTeXPrinting(t *testing.T) {
	b := NewBibTeXEntry("book", "b_id")
	b.AddStringAttribute("title", "Wonderful story")
	b.AddNumericAttribute("year", 1999)
	actual := b.PrettyPrint()
	expected := simpleBibTeXDocument
	if actual != expected {
		t.Errorf("'%s' != '%s'", actual, expected)
	}
}

func TestBibTeXLexing(t *testing.T) {
	for doc, expected := range map[string]string{
		simpleBibTeXDocument:  `[entryStart "@" identifier "book" { identifier "b_id" , identifier "title" = string "Wonderful story" , identifier "year" = number "1999" , } EOF]`,
		exampleBibTeXDocument: `[entryStart "@" identifier "book" { identifier "hicks2001" , identifier "author" = string "von Hicks, III, Michael" , identifier "title" = string "Design of a Carbon Fiber Compo"... , identifier "publisher" = string "Stanford Press" , identifier "year" = number "2001" , identifier "address" = string "Palo Alto" , identifier "edition" = string "1st" , identifier "isbn" = string "0-69-697269-4" , } EOF]`,
	} {
		_, ls := lexBibTeX(doc)
		tokens := make([]lexeme, 0)
		for lexeme := range ls {
			tokens = append(tokens, lexeme)
		}
		//expected := `[entryStart "@" identifier "Book" { identifier "b_id" , identifier "title" = string "Wonderful story" , } EOF]`
		actual := fmt.Sprint(tokens)
		if actual != expected {
			t.Errorf("'%s' != '%s'", actual, expected)
		}
	}
}

func TestSimpleBibTeXParsing(t *testing.T) {
	b := NewBibTeXEntry("book", "b_id")
	b.AddStringAttribute("title", "Wonderful story")

	parsed, err := ParseBibTeXEntry(b.PrettyPrint())
	if err != nil {
		t.Error(err)
	}
	actual := parsed.PrettyPrint()
	expected := b.PrettyPrint()
	if actual != expected {
		t.Errorf("'%s' != '%s'", actual, expected)
	}
}

func TestMoreComplexBibTeXParsing(t *testing.T) {
	for _, doc := range []string{simpleBibTeXDocument, exampleBibTeXDocument} {
		parsed, err := ParseBibTeXEntry(doc)
		if err != nil {
			t.Error(err)
		}
		actual := parsed.PrettyPrint()
		expected := doc
		if actual != expected {
			t.Errorf("'%s' != '%s'", actual, expected)
		}
	}

}

func TestParsingMultiple(t *testing.T) {
	gibberish := " a bunch of gibberish "
	almostMatch := "@fooo}=\n\t"

	cases := map[string]int{
		simpleBibTeXDocument:                                     1,
		exampleBibTeXDocument:                                    1,
		exampleBibTeXDocument + gibberish + simpleBibTeXDocument: 2,
		simpleBibTeXDocument + gibberish + simpleBibTeXDocument:  2,
		exampleBibTeXDocument + almostMatch:                      1,
		almostMatch + simpleBibTeXDocument + gibberish:           1,
		exampleBibTeXDocument + gibberish + almostMatch +
			exampleBibTeXDocument + almostMatch: 2,
	}

	for doc, expected := range cases {
		entries := ParseBibTeXEntries(doc)
		if len(entries) != expected {
			for _, result := range entries {
				fmt.Println(result.PrettyPrint())
			}
			t.Error("Expected", expected, "entries, got", len(entries))
		}
	}
}

func TestParsingBrokenInputs(t *testing.T) {
	brokenDocuments := []string{
		"@article{",
		"@article{id",
		"@article{id,",
		"@article{id,foo ",
		"@article{id,foo = ",
		"@article{id,foo = 123",
	}
	for i := 0; i < len(exampleBibTeXDocument); i++ {
		brokenDocuments = append(brokenDocuments, exampleBibTeXDocument[0:i])
	}

	for _, doc := range brokenDocuments {
		_, err := ParseBibTeXEntry(doc)
		if err == nil {
			t.Error("expected error but did not get one on ", doc)
		}
	}
}
