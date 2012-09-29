package bibtex

import (
    "fmt"
)

type BibTeXEntry struct {
	Type string
        Identifier string
        Attributes map[string]string
        attributeOrder []string
}

func NewBibTeXEntry(entryType, identifier string) BibTeXEntry {
    // @todo add type validation
    return BibTeXEntry{
        Type: entryType,
        Identifier: identifier,
        Attributes: make(map[string]string, 0),
        attributeOrder: make([]string, 0),
    }
}

func (bte *BibTeXEntry) AddAttribute(key, value string) error {
    // @todo add key validation (based on type)
    if _, present := bte.Attributes[key]; !present {
        bte.attributeOrder = append(bte.attributeOrder, key)
    }
    bte.Attributes[key] = value
    return nil
}

func (bte BibTeXEntry) attributesString() (result string) {
    for _,k := range bte.attributeOrder {
        result += fmt.Sprintf("%s = \"%s\"\n", k, bte.Attributes[k])
    }
    return
}

func (bte BibTeXEntry) String() string {
    return fmt.Sprintf("@%s{%s,\n%s}", bte.Type, bte.Identifier, bte.attributesString())
}

