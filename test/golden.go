package test

import (
  "testing"

  goldenbase "github.com/jimmc/golden/base"
)

type Runner struct {
  goldenbase.Runner
  TplFilePath string
}

func NewRunner(basename string) *Runner {
  r := &Runner{}
  r.BaseName = basename
  r.TplFilePath = "testdata/" + basename + ".tpl"
  return r
}

func (r *Runner) Setup(t *testing.T) {
  t.Helper()
  r.SetupT(t)
}

func (r *Runner) Finish(t *testing.T) {
  t.Helper()
  r.FinishT(t)
}
