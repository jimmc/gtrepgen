package gen

import (
  "testing"

  "github.com/jimmc/gtrepgen/data"
  gentest "github.com/jimmc/gtrepgen/test"
)

func TestFromString(t *testing.T) {
  basename := "fromstring"
  dot := "World"
  templ := "Hello {{.}}\n"

  d := gentest.Setup(t, basename)

  g := New("test", false, d.OutW, &data.EmptySource{})
  if err := g.FromString(templ, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}

func TestFromPath(t *testing.T) {
  basename := "helloworld"
  dot := "World"

  d := gentest.Setup(t, basename)

  g := New("test", false, d.OutW, &data.EmptySource{})
  if err := g.FromPath(d.TplFilePath, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}

func TestFromForm(t *testing.T) {
  formname := "org.jimmc.gtrepgen.test1"
  refdirpath := "testdata"
  dot := "World"

  d := gentest.Setup(t, formname)

  g := New(formname, false, d.OutW, &data.EmptySource{})
  if err := g.FromForm(refdirpath, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}

func TestHTML(t *testing.T) {
  basename := "htmlfromstring"
  dot := "<World>"
  templ := "Hello {{.}}\n"

  d := gentest.Setup(t, basename)

  g := New("test", true, d.OutW, &data.EmptySource{})
  if err := g.FromString(templ, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}

func TestInclude(t *testing.T) {
  formname := "org.jimmc.gtrepgen.testinclude"
  refdirpath := "testdata"
  dot := "World"

  d := gentest.Setup(t, formname)

  g := New(formname, false, d.OutW, &data.EmptySource{})
  if err := g.FromForm(refdirpath, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}
