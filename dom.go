// Package dom privides the minimal interfaces to manipulate XML, which is developed on top of the standard xml package.
package dom

import (
	"encoding/xml"
	"errors"
	"regexp"
	"strings"
)

type (
	// Node is an interface that holds Element, xml.Comment or xml.CharData
	Node interface{}

	// Element represents an XML element
	Element struct {
		Name     xml.Name
		Attr     []xml.Attr
		Children []Node
	}
)

var (
	// ErrBreak ...
	ErrBreak = errors.New("Break")

	regSelfClosing = regexp.MustCompile(`></[^>]+>`)
)

// MarshalXML implements xml.Marshaler interface
func (elem *Element) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	s := xml.StartElement{Name: elem.Name, Attr: elem.Attr}
	if err = e.EncodeToken(s); err != nil {
		return
	}

	for _, child := range elem.Children {
		switch node := child.(type) {
		case *Element:
			if err = e.Encode(node); err != nil {
				return
			}
		case xml.CharData, xml.Comment, xml.Directive:
			if err = e.EncodeToken(node); err != nil {
				return
			}
		}
	}

	if err = e.EncodeToken(xml.EndElement{Name: elem.Name}); err != nil {
		return
	}

	return
}

// UnmarshalXML implements xml.Unmarshaler interface
func (elem *Element) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	copy := start.Copy()
	elem.Name.Local = copy.Name.Local
	elem.Attr = copy.Attr
	var next xml.Token

loop:
	for {
		switch next, err = d.Token(); token := next.(type) {
		case xml.CharData:
			// Ignore whitespaces
			if text := strings.Trim(string(token), " \r\n\t"); len(text) > 0 {
				elem.Children = append(elem.Children, token.Copy())
			}
		case xml.Comment, xml.Directive:
			elem.Children = append(elem.Children, xml.CopyToken(token))
		case xml.StartElement:
			child := &Element{}
			if err = d.DecodeElement(child, &token); err != nil {
				break loop
			}
			elem.Children = append(elem.Children, child)
		case xml.EndElement:
			break loop
		default:
			if err != nil {
				break loop
			}
		}
	}
	return
}

// ForEachChild invokes fn on each child element.
//
// The iteration can be broken when fn returns ErrBreak.
// This function returns a child element where fn returned ErrBreak.
// Any other errors from fn causes the iteration to be broken immediately and the error is
// directly returned from this function with nil Element.
func (elem *Element) ForEachChild(fn func(child *Element) error) (res *Element, err error) {
	for _, child := range elem.Children {
		if childElem, ok := child.(*Element); ok == true {
			if err = fn(childElem); err != nil {
				if err == ErrBreak {
					err = nil
					res = childElem
				}
				return
			}
		}
	}
	return
}

// ForEachChildPred invokes fn on each child element where pred returns true.
// See also ForEachChild for the specifications of the return values.
func (elem *Element) ForEachChildPred(pred func(child *Element) bool, fn func(child *Element) error) (res *Element, err error) {
	return elem.ForEachChild(func(child *Element) error {
		if pred(child) == false {
			return nil
		}
		return fn(child)
	})
}

// ForEachChildNamed invokes fn on each child element whose Name is equal to name.
// See also ForEachChild for the specifications of the return values.
func (elem *Element) ForEachChildNamed(name string, fn func(child *Element) error) (res *Element, err error) {
	return elem.ForEachChildPred(
		func(child *Element) bool {
			return child.Name.Local == name
		},
		func(child *Element) error {
			return fn(child)
		})
}

// Marshal returns the XML encoding of elem.
func (elem *Element) Marshal(escQuot, escApos bool) (res string, err error) {
	dat, err := xml.Marshal(elem)
	if err != nil {
		return "", err
	}

	res = string(dat)

	if escQuot == false {
		res = strings.ReplaceAll(res, "&#34;", `"`)
	}

	if escApos == false {
		res = strings.ReplaceAll(res, "&#39;", "'")
	}

	return
}

// MarshalIndent works like Marshal, but XML element begins on a new indented line that starts
// with prefix and is followed by one or more copies of indent according to the nesting depth.
func (elem *Element) MarshalIndent(prefix, indent string, withDecl, escQuot, escApos bool) (res string, err error) {
	dat, err := xml.MarshalIndent(elem, prefix, indent)
	if err != nil {
		return "", err
	}

	res = string(dat)

	if escQuot == false {
		res = strings.ReplaceAll(res, "&#34;", `"`)
	}

	if escApos == false {
		res = strings.ReplaceAll(res, "&#39;", "'")
	}

	res = regSelfClosing.ReplaceAllString(res, " />")

	if withDecl == true {
		res = `<?xml version="1.0" encoding="utf-8"?>` + "\n" + res
	}

	return
}
