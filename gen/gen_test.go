package gen

import (
  "testing"
  "time"

  "github.com/jimmc/gtrepgen/data"

  goldenbase "github.com/jimmc/golden/base"
)

func TestFromString(t *testing.T) {
  basename := "fromstring"
  dot := "World"
  templ := "Hello {{.}}\n"

  r := goldenbase.NewTester(basename)
  goldenbase.FatalIfError(t, r.Arrange(), "Arrange")

  g := New("test", false, r.OutW, &data.EmptySource{})
  if err := g.FromString(templ, dot); err != nil {
    t.Fatal(err)
  }

  goldenbase.FatalIfError(t, r.Assert(), "Assert")
}

func TestFromPath(t *testing.T) {
  basename := "helloworld"
  dot := "World"

  r := goldenbase.NewTester(basename)
  goldenbase.FatalIfError(t, r.Arrange(), "Arrange")

  tplFilePath := "testdata/" + basename + ".tpl"
  g := New("test", false, r.OutW, &data.EmptySource{})
  if err := g.FromPath(tplFilePath, dot); err != nil {
    t.Fatal(err)
  }

  goldenbase.FatalIfError(t, r.Assert(), "Assert")
}

func TestFromTemplate(t *testing.T) {
  tplname := "org.jimmc.gtrepgen.test1"
  refdirpaths := []string{"testdata"}
  dot := "World"

  r := goldenbase.NewTester(tplname)
  goldenbase.FatalIfError(t, r.Arrange(), "Arrange")

  g := New(tplname, false, r.OutW, &data.EmptySource{})
  if err := g.FromTemplate(refdirpaths, dot); err != nil {
    t.Fatal(err)
  }

  goldenbase.FatalIfError(t, r.Assert(), "Assert")
}

func TestHTML(t *testing.T) {
  basename := "htmlfromstring"
  dot := "<World>"
  templ := "Hello {{.}}\n"

  r := goldenbase.NewTester(basename)
  goldenbase.FatalIfError(t, r.Arrange(), "Arrange")

  g := New("test", true, r.OutW, &data.EmptySource{})
  if err := g.FromString(templ, dot); err != nil {
    t.Fatal(err)
  }

  goldenbase.FatalIfError(t, r.Assert(), "Assert")
}

func TestInclude(t *testing.T) {
  tplname := "org.jimmc.gtrepgen.testinclude"
  refdirpaths := []string{"testdata"}
  dot := "World"

  r := goldenbase.NewTester(tplname)
  goldenbase.FatalIfError(t, r.Arrange(), "Arrange")

  g := New(tplname, false, r.OutW, &data.EmptySource{})
  if err := g.FromTemplate(refdirpaths, dot); err != nil {
    t.Fatal(err)
  }

  goldenbase.FatalIfError(t, r.Assert(), "Assert")
}

func TestIncludeResult(t *testing.T) {
  tplname := "org.jimmc.gtrepgen.testincluderesult"
  refdirpaths := []string{"testdata"}
  dot := "World"

  r := goldenbase.NewTester(tplname)
  goldenbase.FatalIfError(t, r.Arrange(), "Arrange")

  g := New(tplname, false, r.OutW, &data.EmptySource{})
  if err := g.FromTemplate(refdirpaths, dot); err != nil {
    t.Fatal(err)
  }

  goldenbase.FatalIfError(t, r.Assert(), "Assert")
}

func TestIncludeTwoDirs(t *testing.T) {
  tplname := "inc1"
  refdirpaths := []string{
    "testdata",
    "../dbsource/testdata",
  }
  dot := "World"

  r := goldenbase.NewTester(tplname)
  goldenbase.FatalIfError(t, r.Arrange(), "Arrange")

  g := New(tplname, false, r.OutW, &data.EmptySource{})
  if err := g.FromTemplate(refdirpaths, dot); err != nil {
    t.Fatal(err)
  }

  goldenbase.FatalIfError(t, r.Assert(), "Assert")
}

func TestEvenOdd(t *testing.T) {
  if got, want := evenodd(0, "a", "b"), "a"; got != want {
    t.Errorf("evenodd got %v, want %v", got, want)
  }
  if got, want := evenodd(3, "a", "b"), "b"; got != want {
    t.Errorf("evenodd got %v, want %v", got, want)
  }
  if got, want := evenodd(-4, "a", "b"), "a"; got != want {
    t.Errorf("evenodd got %v, want %v", got, want)
  }
  if got, want := evenodd(-3, "a", "b"), "b"; got != want {
    t.Errorf("evenodd got %v, want %v", got, want)
  }
}

func TestFormatTime(t *testing.T) {
  testTime, err := time.Parse("Jan 2, 2006 15:04:05", "Feb 1, 2019 14:34:20")
  if err != nil {
    t.Fatalf("Error parsing test time")
  }
  if got, want := formatTime("Jan 2, 2006 3:04:05 PM", testTime), "Feb 1, 2019 2:34:20 PM"; got != want {
    t.Errorf("formatTime got %v, want %v", got, want)
  }
}

func TestMkmap(t *testing.T) {
  m, err := mkmap("a", "abc", "b", "bbb", 1, 123)
  if err != nil {
    t.Fatalf("Unexpected error from mkmap: %v", err)
  }
  if got, want := len(m), 3; got != want {
    t.Fatalf("Wrong length, got %d, want %d", got, want)
  }
  if got, want := m["a"], "abc"; got != want {
    t.Errorf("m[a]: got %v, want %v", got, want)
  }
  if got, want := m[1], 123; got != want {
    t.Errorf("m[a]: got %v, want %v", got, want)
  }

  _, err = mkmap("x", "y", "z")
  if err == nil {
    t.Fatalf("Expected error from mkmap with odd arg count")
  }
}
