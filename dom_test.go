package dom

import (
	"encoding/xml"
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
	elem := &Element{}
	err := xml.Unmarshal([]byte(input), elem)
	if err != nil {
		t.Fatal(err)
	}

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
}
