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

type nameAndNum struct {
  Name string
  Num int
}

func TestReadAttributesFromStringInto(t *testing.T) {
  got := nameAndNum{}
  err := ReadTemplateAttributesFromStringInto(`{{/*GT: {"Name":"foo", "Num":123} */ -}}`, &got)
  if err != nil {
    t.Fatalf("Reading attributes: %v", err)
  }
  want := nameAndNum{"foo", 123}
  if got != want {
    t.Fatalf("FromStringInto: got %v(%T), want %v(%T)", got, got, want, want)
  }
}

func TestReadAttributesFromPath(t *testing.T) {
  expected := map[string]interface{}{
    "display": "Hello World",
    "x": float64(1),
  }
  a, err := ReadTemplateAttributesFromPath("testdata/helloworld.tpl")
  if err != nil {
    t.Fatalf("Reading attributes from file: %v", err)
  }
  if got, want := a, expected; !cmp.Equal(got, want) {
    t.Fatalf("Attributes: got %+v, want %+v", got, want)
  }
}

type displayAndX struct {
  Display string
  X float64
}

func TestReadAttributesFromPathInto(t *testing.T) {
  got := displayAndX{}
  want := displayAndX{
    Display: "Hello World",
    X: float64(1),
  }
  if err := ReadTemplateAttributesFromPathInto("testdata/helloworld.tpl", &got); err != nil {
    t.Fatalf("Reading attributes from file: %v", err)
  }
  if !cmp.Equal(got, want) {
    t.Fatalf("Attributes: got %+v, want %+v", got, want)
  }
}

type displayOnly struct {
  Display string
}

func TestReadAttributesFromPathIntoPartialMatch(t *testing.T) {
  got := displayOnly{}
  want := displayOnly{
    Display: "Hello World",
  }
  if err := ReadTemplateAttributesFromPathInto("testdata/helloworld.tpl", &got); err != nil {
    t.Fatalf("Reading attributes from file: %v", err)
  }
  if !cmp.Equal(got, want) {
    t.Fatalf("Attributes: got %+v, want %+v", got, want)
  }
}

type displayAndString struct {
  Display string
  X string
}

func TestReadAttributesFromPathIntoMismatch(t *testing.T) {
  got := displayAndString{}
  err := ReadTemplateAttributesFromPathInto("testdata/helloworld.tpl", &got)
  if err == nil {
    t.Fatalf("Expected error about type mismatch")
  }
}

func TestReadDirFilesAttributes(t *testing.T) {
  var m0, m1 interface{}
  m0 = map[string]interface{}{"display": "Hello World", "x": float64(1)}
  m1 = map[string]interface{}{"display": "again"}
  expected := []*TemplateAttributes{
    {
      Name: "helloworld",
      Attributes: &m0,
    },
    {
      Name: "org.jimmc.gtrepgen.test1",
      Attributes: &m1,
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

func TestReadDirFilesAttributesAs(t *testing.T) {
  expected := []*TemplateAttributes{
    {
      Name: "helloworld",
      Attributes: &displayAndX{Display: "Hello World", X: 1},
    },
    {
      Name: "org.jimmc.gtrepgen.test1",
      Attributes: &displayAndX{Display: "again"},
    },
  }
  newDest := func() interface{} {
    return &displayAndX{}
  }
  attrsList, err := ReadDirFilesAttributesAs("testdata", newDest)
  if err != nil {
    t.Fatalf("Error reading dir files attributes: %v", err)
  }
  if attrsList == nil {
    t.Fatal("Nil attrs list")
  }
  got, want := attrsList, expected
  if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("ReadDirFilesAttributesAs() mismatch (-want +got):\n%s", diff)
  }
}
