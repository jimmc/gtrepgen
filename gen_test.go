package gen

import (
  "bufio"
  "os"
  "testing"

  gentest "github.com/jimmc/gtrepgen/test"
)

func TestFromSTring(t *testing.T) {
  basename := "fromstring"
  dot := "World"
  templ := "Hello {{.}}\n"

  outfilepath := "testdata/" + basename + ".out"
  goldenfilepath := "testdata/" + basename + ".golden"

  os.Remove(outfilepath)
  f, err := os.Create(outfilepath)
  if err != nil {
    t.Fatal(err)
  }
  w := bufio.NewWriter(f)

  if err := FromString(w, "test", templ, dot); err != nil {
    t.Fatal(err)
  }

  w.Flush()
  f.Close()

  if err := gentest.CompareOutToGolden(outfilepath, goldenfilepath); err != nil {
    t.Fatal(err)
  }
}

func TestFromPath(t *testing.T) {
  basename := "helloworld"
  dot := "World"

  infilepath := "testdata/" + basename + ".tpl"
  outfilepath := "testdata/" + basename + ".out"
  goldenfilepath := "testdata/" + basename + ".golden"

  os.Remove(outfilepath)
  f, err := os.Create(outfilepath)
  if err != nil {
    t.Fatal(err)
  }
  w := bufio.NewWriter(f)

  if err := FromPath(w, "test", infilepath, dot); err != nil {
    t.Fatal(err)
  }

  w.Flush()
  f.Close()

  if err := gentest.CompareOutToGolden(outfilepath, goldenfilepath); err != nil {
    t.Fatal(err)
  }
}

func TestFromForm(t *testing.T) {
  formname := "org.jimmc.gtrepgen.test1"
  dot := "World"

  refdirpath := "testdata"
  outfilepath := "testdata/" + formname + ".out"
  goldenfilepath := "testdata/" + formname + ".golden"

  os.Remove(outfilepath)
  f, err := os.Create(outfilepath)
  if err != nil {
    t.Fatal(err)
  }
  w := bufio.NewWriter(f)

  if err := FromForm(w, formname, refdirpath, dot); err != nil {
    t.Fatal(err)
  }

  w.Flush()
  f.Close()

  if err := gentest.CompareOutToGolden(outfilepath, goldenfilepath); err != nil {
    t.Fatal(err)
  }
}
