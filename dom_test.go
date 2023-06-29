package dom

import (
	"encoding/xml"
	"log"
	"strings"
	"testing"
)

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

func TestForEachChild(t *testing.T) {
	elem := &Element{}
	xml.Unmarshal([]byte(`<a><b/><c/><d/>text<e/>text</a>`), elem)
	childCount := 0
	elem.ForEachChild(func(child *Element) error {
		childCount++
		return nil
	})

	if childCount != 4 {
		t.Fatal("childCount != 4")
	}

	b, err := elem.ForEachChild(func(child *Element) error {
		if child.Name.Local == "b" {
			return ErrBreak
		}
		return nil
	})

	if err != nil || b == nil || b.Name.Local != "b" {
		t.Fatal("ForEachChild with ErrBreak failed.")
	}
}

func TestForEachChildNamed(t *testing.T) {
	elem := &Element{}
	xml.Unmarshal([]byte(`<a><b/><c/><d/>text<e/>text<!--comment--><c/><c></c></a>`), elem)
	childCount := 0
	elem.ForEachChildNamed("c", func(child *Element) error {
		childCount++
		return nil
	})

	if childCount != 3 {
		t.Fatal("childCount != 3")
	}

	childCount = 0
	c, err := elem.ForEachChildNamed("c", func(child *Element) error {
		childCount++
		return ErrBreak
	})

	if err != nil || c == nil || c.Name.Local != "c" {
		t.Fatal("ForEachChildNamed with ErrBreak failed.")
	}

	if childCount != 1 {
		t.Fatal("childCount != 1")
	}
}

func TestError(t *testing.T) {
	elem := &Element{}
	err := xml.Unmarshal([]byte(`<a><b/><c/><d/>text<e/>text<!--comment--><x><c/</x><c></c></a>`), elem)
	if err == nil {
		t.Fatal("Unmarshal error is expected.")
	}
}

func TestMarshal(t *testing.T) {
	input := `<PropertyGroup Condition="'$(CompileConfig)' == 'DEBUG'">
  <Optimization>false</Optimization>
  <Obfuscate>false</Obfuscate>
  <OutputPath>$(OutputPath)\debug</OutputPath>
</PropertyGroup>`

	elem := Must(input)
	m0, err := elem.MarshalIndent("", "  ", true, false, false)
	if err != nil {
		t.Fatal(err)
	}

	elem = &Element{}
	err = xml.Unmarshal([]byte(m0), elem)
	if err != nil {
		t.Fatal(err)
	}

	m1, err := elem.Marshal(false, false)
	if err != nil {
		t.Fatal(err)
	}

	elem = &Element{}
	err = xml.Unmarshal([]byte(m1), elem)
	if err != nil {
		t.Fatal(err)
	}

	m2, err := elem.MarshalIndent("", "  ", false, false, false)
	if err != nil {
		t.Fatal(err)
	}

	if m2 != input {
		t.Fatal("m1 != input")
	}

	elem = nil
	if res, _ := elem.MarshalIndent("", "  ", false, false, false); len(res) != 0 {
		t.Fatal("elem.MarshalIndent() must return empty string if elem == nil")
	}
}

func TestIsEmpty(t *testing.T) {
	elem := &Element{}
	if elem.IsEmpty() == false {
		t.Fatal("elem.IsEmpty() == false")
	}

	elem.Children = append(elem.Children, &Element{})
	if elem.IsEmpty() == true {
		t.Fatal("elem.IsEmpty() == true")
	}

	elem = nil
	if elem.IsEmpty() == false {
		t.Fatal("elem.IsEmpty() == false")
	}

	elem = &Element{}
	elem.Attr = append(elem.Attr, xml.Attr{})
	if elem.IsEmpty() == true {
		t.Fatal("elem.IsEmpty() == true")
	}
}

func TestText(t *testing.T) {
	elem := Must(`<a><s1/><s2></s2><s3>text</s3></a>`)
	text, res := elem.Text()
	if len(text) > 0 || res == true {
		t.Fatal(`len(text) > 0 || res == true`)
	}

	elem.ForEachChild(func(child *Element) error {
		switch child.Name.Local {
		case "s1", "s2":
			if text, res = child.Text(); len(text) > 0 || res == true {
				t.Fatal(`len(text) > 0 || res == true`)
			}
		case "s3":
			if text, res = child.Text(); res == false || text != "text" {
				t.Fatal(`res == false || text != "text"`)
			}
		}
		return nil
	})

	// It replaces children with a text
	elem.SetText("text")
	text, res = elem.Text()
	if res == false || text != "text" {
		t.Fatal(`res == false || text != "text"`)
	}

	// It clears children if text is empty
	elem.SetText("")
	text, res = elem.Text()
	if len(text) > 0 || res == true {
		t.Fatal(`len(text) > 0 || res == true`)
	}

	// Nothing happens if elem is nil
	elem = nil
	elem.SetText("text")
	text, res = elem.Text()
	if len(text) > 0 || res == true {
		t.Fatal(`len(text) > 0 || res == true`)
	}
}

func TestFindAttr(t *testing.T) {
	elem := Must(`<a attr1="test1" attr2="test2" attr3="test3"/>`)
	attr := elem.FindAttr("attr1")
	if attr == nil {
		t.Fatal(`attr == nil`)
	}
	if attr.Name.Local != "attr1" || attr.Value != "test1" {
		t.Fatal(`attr.Name.Local != "attr1" || attr.Value != "test1"`)
	}
	if elem.HasAttr("attr2") == false {
		t.Fatal(`elem.HasAttr("attr2") == false`)
	}
	if elem.HasAttr("attr4") == true {
		t.Fatal(`elem.HasAttr("attr4") == true`)
	}
	elem = nil
	if elem.FindAttr("test1") != nil {
		log.Fatal(`elem.FindAttr("test1") != nil`)
	}
}

func TestTextRecurse(t *testing.T) {
	input := `<PropertyGroup Condition="'$(CompileConfig)' == 'DEBUG'">
	ThisIs
  <Optimization>Test</Optimization>
  <Obfuscate>De</Obfuscate>
  <OutputPath>Gonsu</OutputPath>Gonsu
</PropertyGroup>`
	elem := Must(input)
	res := elem.TextRecurse()
	if res != "ThisIsTestDeGonsuGonsu" {
		t.Fatal(res)
	}
}
