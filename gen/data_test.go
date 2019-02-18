package gen

import (
  "fmt"
  "testing"

  gentest "github.com/jimmc/gtrepgen/test"
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
  formname := "org.jimmc.gtrepgen.datatest"
  refdirpath := "testdata"
  dot := "top"

  d := gentest.Setup(t, formname)

  if err := FromForm(d.OutW, &TestSource{}, formname, refdirpath, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}
