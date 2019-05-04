package gen

import (
  "fmt"
  "testing"

  goldenbase "github.com/jimmc/golden/base"
)

type TestSource struct{}

func (s *TestSource) Row(args ...interface{}) (interface{}, error) {
  a0, ok := args[0].(string)
  if !ok {
    return nil, fmt.Errorf("TestSource.Row first arg must be string")
  }
  return map[string]interface{}{
    "a": 1,
    "b": 2.2,
    "c": "Three:"+a0,
  }, nil
}

func (s *TestSource) Rows(args ...interface{}) (interface{}, error) {
  a0, ok := args[0].(string)
  if !ok {
    return nil, fmt.Errorf("TestSource.Row first arg must be string")
  }
  r0 := map[string]interface{}{
    "a": 1,
    "b": 2.2,
    "c": "Three:"+a0,
  }
  r1 := map[string]interface{}{
    "a": 11,
    "b": 12.2,
    "c": "Thirteen:"+a0,
  }
  r2 := map[string]interface{}{
    "a": 21,
    "b": 22.2,
    "c": "Twentythree:"+a0,
  }
  return []interface{}{r0, r1, r2}, nil
}

func TestDataSource(t *testing.T) {
  tplname := "org.jimmc.gtrepgen.datatest"
  refdirpaths := []string{"testdata"}
  dot := "top"

  r := goldenbase.NewTester(tplname)
  r.SetupT(t)

  g := New(tplname, false, r.OutW, &TestSource{})
  if err := g.FromTemplate(refdirpaths, dot); err != nil {
    t.Fatal(err)
  }

  r.FinishT(t)
}
