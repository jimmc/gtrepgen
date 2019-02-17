package gen

import (
  "testing"

  gentest "github.com/jimmc/gtrepgen/test"
)

func TestFromSTring(t *testing.T) {
  basename := "fromstring"
  dot := "World"
  templ := "Hello {{.}}\n"

  d := gentest.Setup(t, basename)

  if err := FromString(d.OutW, "test", templ, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}

func TestFromPath(t *testing.T) {
  basename := "helloworld"
  dot := "World"

  d := gentest.Setup(t, basename)

  if err := FromPath(d.OutW, "test", d.TplFilePath, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}

func TestFromForm(t *testing.T) {
  formname := "org.jimmc.gtrepgen.test1"
  refdirpath := "testdata"
  dot := "World"

  d := gentest.Setup(t, formname)

  if err := FromForm(d.OutW, formname, refdirpath, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}
