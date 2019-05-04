package dbsource

import (
  "testing"

  "github.com/jimmc/gtrepgen/gen"
  gentest "github.com/jimmc/gtrepgen/test"

  goldendb "github.com/jimmc/golden/db"

  _ "github.com/mattn/go-sqlite3"
)

func TestDbSource(t *testing.T) {
  sqlInitFile := "testdata/dbsourcetest.sql"
  db, err := goldendb.EmptyDb()
  if err != nil {
    t.Fatal(err)
  }
  if err := goldendb.LoadSetupFile(db, sqlInitFile); err != nil {
    t.Fatal(err)
  }
  dataSource := New(db)
  tplname := "org.jimmc.gtrepgen.sqltest"
  refdirpaths := []string{"testdata"}
  dot := "x"

  r := gentest.NewRunner(tplname)
  r.SetupT(t)

  g := gen.New(tplname, false, r.OutW, dataSource)
  if err := g.FromTemplate(refdirpaths, dot); err != nil {
    t.Fatal(err)
  }

  r.FinishT(t)
}
