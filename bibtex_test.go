package bibtex

import (
	"fmt"
	"testing"
)

func TestBibTeXStringing(t *testing.T) {
	b := NewBibTeXEntry("book", "b_id")
	b.AddAttribute("title", "Wonderful story")
	actual := fmt.Sprint(b)
	expected := `@book{b_id,
title = "Wonderful story"
}`
	if actual != expected {
		t.Errorf("'%s' !+ '%s'", actual, expected)
	}
}

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
	for doc,expected := range map[string]string{
		simpleBibTeXDocument: 
		`[entryStart "@" identifier "Book" { identifier "b_id" , identifier "title" = string "Wonderful story" , } EOF]`,
		exampleBibTeXDocument: 
		`[entryStart "@" identifier "Book" { identifier "hicks2001" , identifier "author" = string "von Hicks, III, Michael" , identifier "title" = string "Design of a Carbon Fiber Compo"... , identifier "publisher" = string "Stanford Press" , identifier "year" = identifier "2001" , identifier "address" = string "Palo Alto" , identifier "edition" = string "1st" , identifier "isbn" = string "0-69-697269-4" } EOF]`,
	} {
		_, ls := lexBibTeX(doc)
		tokens := make([]lexeme, 0)
		for lexeme := range ls {
			tokens = append(tokens, lexeme)
		}
		//expected := `[entryStart "@" identifier "Book" { identifier "b_id" , identifier "title" = string "Wonderful story" , } EOF]`
		actual := fmt.Sprint(tokens)
		if actual != expected {
			t.Errorf("'%s' !+ '%s'", actual, expected)
		}
	}
}
