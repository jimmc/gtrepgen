package gen

import (
  "testing"
  "time"

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

func TestFromTemplate(t *testing.T) {
  tplname := "org.jimmc.gtrepgen.test1"
  refdirpaths := []string{"testdata"}
  dot := "World"

  d := gentest.Setup(t, tplname)

  g := New(tplname, false, d.OutW, &data.EmptySource{})
  if err := g.FromTemplate(refdirpaths, dot); err != nil {
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
  tplname := "org.jimmc.gtrepgen.testinclude"
  refdirpaths := []string{"testdata"}
  dot := "World"

  d := gentest.Setup(t, tplname)

  g := New(tplname, false, d.OutW, &data.EmptySource{})
  if err := g.FromTemplate(refdirpaths, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}

func TestIncludeResult(t *testing.T) {
  tplname := "org.jimmc.gtrepgen.testincluderesult"
  refdirpaths := []string{"testdata"}
  dot := "World"

  d := gentest.Setup(t, tplname)

  g := New(tplname, false, d.OutW, &data.EmptySource{})
  if err := g.FromTemplate(refdirpaths, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}

func TestIncludeTwoDirs(t *testing.T) {
  tplname := "inc1"
  refdirpaths := []string{
    "testdata",
    "../dbsource/testdata",
  }
  dot := "World"

  d := gentest.Setup(t, tplname)

  g := New(tplname, false, d.OutW, &data.EmptySource{})
  if err := g.FromTemplate(refdirpaths, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
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
