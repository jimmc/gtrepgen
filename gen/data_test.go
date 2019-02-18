package gen

import (
  "fmt"
  "testing"

  gentest "github.com/jimmc/gtrepgen/test"
)

type TestSource struct{}

func (s *TestSource) Data(args ...interface{}) (interface{}, error) {
  a0, ok := args[0].(string)
  if !ok {
    return nil, fmt.Errorf("TestSource.Data first arg must be string")
  }
  a1, ok := args[1].(string)
  if !ok {
    return nil, fmt.Errorf("TestSource.Data second arg must be string")
  }
  switch a0 {
  case "field":
    return "<" + a1 + ">", nil
  case "row":
    return map[string]interface{}{
      "a": 1,
      "b": 2.2,
      "c": "Three:"+a1,
    }, nil
  case "rows":
    r0 := map[string]interface{}{
      "a": 1,
      "b": 2.2,
      "c": "Three:"+a1,
    }
    r1 := map[string]interface{}{
      "a": 11,
      "b": 12.2,
      "c": "Thirteen:"+a1,
    }
    r2 := map[string]interface{}{
      "a": 21,
      "b": 22.2,
      "c": "Twentythree:"+a1,
    }
    return []interface{}{r0, r1, r2}, nil
  }
  return nil, fmt.Errorf("unknown return type specifier %q", a0)
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
