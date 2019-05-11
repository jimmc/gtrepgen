package dbsource

import (
  "database/sql"
  "io"
  "testing"

  "github.com/jimmc/gtrepgen/gen"

  goldenbase "github.com/jimmc/golden/base"
  goldendb "github.com/jimmc/golden/db"

  _ "github.com/mattn/go-sqlite3"
)

func TestDbSource(t *testing.T) {
  sqlFile := "testdata/dbsourcetest.sql"
  tplname := "org.jimmc.gtrepgen.sqltest"
  refdirpaths := []string{"testdata"}
  dot := "x"

  callback := func(db *sql.DB, w io.Writer) error {
    dataSource := New(db)
    g := gen.New(tplname, false, w, dataSource)
    return g.FromTemplate(refdirpaths, dot)
  }

  r := goldendb.NewTester(tplname, callback)
  r.SetupPath = sqlFile

  goldenbase.FatalIfError(t, goldenbase.Run(r), "Run")
}
