# dom

## Description

Minimal (DOM-like) interface to manipulate XML.

## Installation

```
go get github.com/xoinu/dom
```

## Examples

Use `xml.Unmarshal`, `xml.Marshal` and `xml.MarshalIndent` to read and to write XML.

```
func TestDom(t *testing.T) {
	input := `<PropertyGroup Condition="'$(CompileConfig)' == 'DEBUG'">
  <Optimization>false</Optimization>
  <Obfuscate>false</Obfuscate>
  <OutputPath>$(OutputPath)\debug</OutputPath>
</PropertyGroup>`

	elem := &Element{}
	xml.Unmarshal([]byte(input), &elem)
	b, err := xml.MarshalIndent(elem, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	ouput := strings.ReplaceAll(string(b), "&#39;", "'")

	if input != ouput {
		t.Fatal("input != output")
	}
}
```

`Element.Marshal` and `Element.MarshalIndent` provide a few more useful options to write XML.

The XML document is loaded into `Element` object which is as simple as follows:

```
Element struct {
    Name     xml.Name
    Attr     []xml.Attr
    Children []Node
}
```

where `Node` is either `Element`, `xml.CharData`, `xml.Comment`, or `xml.Directive`.

`Element.ForEachChild*` family lets you traverse child elements.
