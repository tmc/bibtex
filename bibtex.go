package bibtex

import (
	"fmt"
)

type BibTeXEntry struct {
	Type           string
	Identifier     string
	Attributes     map[string]string
	attributeOrder []string
}

func NewBibTeXEntry(entryType, identifier string) BibTeXEntry {
	// @todo add type validation
	result := newBibTeXEntry()
	result.Type = entryType
	result.Identifier = identifier
	return result
}

func newBibTeXEntry() BibTeXEntry {
	// @todo add type validation
	return BibTeXEntry{
		Attributes:     make(map[string]string, 0),
		attributeOrder: make([]string, 0),
	}
}

func (bte *BibTeXEntry) AddNumericAttribute(key string, val int) error {
	return bte.addAttribute(key, fmt.Sprint(val))
}

func (bte *BibTeXEntry) AddStringAttribute(key, val string) error {
	return bte.addAttribute(key, fmt.Sprintf("\"%s\"", val))
}

func (bte *BibTeXEntry) addAttribute(key, value string) error {
	// @todo add key validation (based on type)
	if _, present := bte.Attributes[key]; !present {
		bte.attributeOrder = append(bte.attributeOrder, key)
	}
	bte.Attributes[key] = value
	return nil
}

func (bte BibTeXEntry) attributesString() (result string) {
	longestKey := 0
	for key := range bte.Attributes {
		if len(key) > longestKey {
			longestKey = len(key)
		}
	}
	for _, k := range bte.attributeOrder {
		result += fmt.Sprintf(" %-"+fmt.Sprint(longestKey)+"s = %s,\n", k, bte.Attributes[k])
	}
	return
}

func (bte BibTeXEntry) PrettyPrint() string {
	return fmt.Sprintf("@%s{%s,\n%s}", bte.Type, bte.Identifier, bte.attributesString())
}
