package gen

import (
  "strings"
  "testing"

  "github.com/google/go-cmp/cmp"
)

func TestExtractAttributesEmpty(t *testing.T) {
  got, err := extractAttributeString(strings.NewReader(""))
  if err != nil {
    t.Fatalf("Extracting string: %v", err)
  }
  want := ""
  if got != want {
    t.Fatalf("Empty string: got %q, want %q", got, want)
  }
}

func TestExtractAttributesNoFields(t *testing.T) {
  got, err := extractAttributeString(strings.NewReader("{{/*GT: */ -}}"))
  if err != nil {
    t.Fatalf("Extracting string: %v", err)
  }
  want := ""
  if got != want {
    t.Fatalf("No fields: got %q, want %q", got, want)
  }
}

func TestExtractAttributesOneFields(t *testing.T) {
  got, err := extractAttributeString(strings.NewReader("{{/*GT: 123 */ -}}"))
  if err != nil {
    t.Fatalf("Extracting string: %v", err)
  }
  want := "123"
  if got != want {
    t.Fatalf("No fields: got %q, want %q", got, want)
  }
}

func TestExtractAttributesThreeLines(t *testing.T) {
  got, err := extractAttributeString(strings.NewReader("{{/*GT: abc\ndef\nghi */ -}}"))
  if err != nil {
    t.Fatalf("Extracting string: %v", err)
  }
  want := "abc\ndef\nghi"
  if got != want {
    t.Fatalf("No fields: got %q, want %q", got, want)
  }
}

func TestReadAttributesFromStringEmpty(t *testing.T) {
  a, err := ReadTemplateAttributesFromString("")
  if err != nil {
    t.Fatalf("Reading attributes: %v", err)
  }
  if a != nil {
    t.Fatalf("Expected no attributes, got %v", a)
  }
}

func TestReadAttributesFromStringOneLine(t *testing.T) {
  a, err := ReadTemplateAttributesFromString("{{/*GT: 123 */ -}}")
  if err != nil {
    t.Fatalf("Reading attributes: %v", err)
  }
  if got, want := a, float64(123); got != want {
    t.Fatalf("One value: got %v(%T), want %v(%T)", got, got, want, want)
  }
}

func TestReadAttributesFromStringOneValueTwoLines(t *testing.T) {
  a, err := ReadTemplateAttributesFromString("{{/*GT: 123\n*/ -}}")
  if err != nil {
    t.Fatalf("Reading attributes: %v", err)
  }
  if got, want := a, 123.0; got != want {
    t.Fatalf("One value: got %v(%T), want %v(%T)", got, got, want, want)
  }
}

func TestReadAttributesFromPath(t *testing.T) {
  expected := map[string]interface{}{
    "display": "Hello World",
    "x": float64(1),
  };
  a, err := ReadTemplateAttributesFromPath("testdata/helloworld.tpl")
  if err != nil {
    t.Fatalf("Reading attributes from file: %v", err)
  }
  if got, want := a, expected; !cmp.Equal(got, want) {
    t.Fatalf("Attributes: got %+v, want %+v", got, want)
  }
}

func TestReadDirFilesAttributes(t *testing.T) {
  expected := []*TemplateAttributes{
    {
      Name: "helloworld",
      Attributes: map[string]interface{}{"display": "Hello World", "x": float64(1)},
    },
    {
      Name: "org.jimmc.gtrepgen.test1",
      Attributes: map[string]interface{}{"display": "again"},
    },
  }
  attrsList, err := ReadDirFilesAttributes("testdata")
  if err != nil {
    t.Fatalf("Error reading dir files attributes: %v", err)
  }
  if attrsList == nil {
    t.Fatal("Nil attrs list")
  }
  got, want := attrsList, expected
  if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("ReadDirFilesAttributes() mismatch (-want +got):\n%s", diff)
  }
}
