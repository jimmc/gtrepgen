package test

import (
  "testing"

  goldenbase "github.com/jimmc/golden/base"
)

type Data struct {
  goldenbase.RunData
  TplFilePath string
  r goldenbase.Runner
}

func Setup(t *testing.T, basename string) *Data {
  t.Helper()
  tplfilepath := "testdata/" + basename + ".tpl"
  r := goldenbase.Runner{
    BaseName: basename,
  }
  runData := r.SetupT(t)

  return &Data{
    RunData: *runData,
    TplFilePath: tplfilepath,
    r: r,
  }
}

func Finish(t *testing.T, data *Data) {
  data.r.FinishT(t, &data.RunData)
}
