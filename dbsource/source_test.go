package dbsource

import (
  "database/sql"
  "io/ioutil"
  "regexp"
  "strings"
  "testing"

  "github.com/jimmc/gtrepgen/gen"
  gentest "github.com/jimmc/gtrepgen/test"

  _ "github.com/mattn/go-sqlite3"
)

func EmptyDb() (*sql.DB, error) {
  return sql.Open("sqlite3", ":memory:")
}

func LoadSetupFile(db *sql.DB, filename string) error {
  setupSql, err := ioutil.ReadFile(filename)
  if err != nil {
    return err
  }
  return ExecMulti(db, string(setupSql))
}

// ExecMulti executes multiple sql statements from a string.
// It strips out all carriage returns, breaks the string into segments at double-newlines, removes lines
// starting with a "#" (as comments), and separately executes
// each segment. If any segment returns an error, it stop executing
// and returns that error.
func ExecMulti(db *sql.DB, sql string) error {
  re := regexp.MustCompile("\r")
  sql = re.ReplaceAllString(sql, "")
  segments := strings.Split(sql, "\n\n")
  for _, segment := range segments {
    if err := ExecSegment(db, segment); err != nil {
      return err
    }
  }
  return nil
}

// ExecSegment executes a single sql statement from a string.
// It removes lines starting with "#" (as comments).
func ExecSegment(db *sql.DB, segment string) error {
  lines := strings.Split(segment, "\n")
  sqlLines := make([]string, 0)
  for _, line := range lines {
    if !strings.HasPrefix(line, "#") {
      sqlLines = append(sqlLines, line)
    }
  }
  segment = strings.Join(sqlLines, "\n")
  _, err := db.Exec(segment)
  return err
}

func TestDbSource(t *testing.T) {
  sqlInitFile := "testdata/dbsourcetest.sql"
  db, err := EmptyDb()
  if err != nil {
    t.Fatal(err)
  }
  if err := LoadSetupFile(db, sqlInitFile); err != nil {
    t.Fatal(err)
  }
  dataSource := New(db)
  formname := "org.jimmc.gtrepgen.sqltest"
  refdirpath := "testdata"
  dot := "x"

  d := gentest.Setup(t, formname)

  g := gen.New(formname, false, d.OutW, dataSource)
  if err := g.FromForm(refdirpath, dot); err != nil {
    t.Fatal(err)
  }

  gentest.Finish(t, d)
}
