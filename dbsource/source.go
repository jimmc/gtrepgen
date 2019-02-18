package dbsource

import (
  "database/sql"
  "errors"
  "fmt"
  "log"
)

// dbQuery represents the functions we use from a database.
// It is designed to accept either a sql.DB or a sql.Tx.
type dbQuery interface{
  Query(query string, args ...interface{}) (*sql.Rows, error)
}

type SqlSource struct{
  db dbQuery;
}

// New creates a new SqlSource from a database or transaction.
func New(db dbQuery) *SqlSource {
  return &SqlSource{
    db: db,
  }
}

// Rows gets multiple rows from our database. The first arg is the query string, and all
// following args are passed to the Query function as arguments for the query string.
// This function may return zero or more rows.
func (s *SqlSource) Rows(args ...interface{}) (interface{}, error) {
  query, ok := args[0].(string)
  if !ok {
    return nil, errors.New("SqlSource.Data first arg must be string (query)")
  }
  queryArgs := args[1:]
  rr, err := s.db.Query(query, queryArgs...)
  if err != nil {
    return nil, err
  }
  log.Printf("Got query results")
  var colNames []string
  fieldCount := -1
  // We need to return either an array or a channel.
  // For now, read all the rows and return an array.
  data := make([]map[string]interface{}, 0)
  for rr.Next() {
    if fieldCount < 0 {
      colNames, err = rr.Columns()
      if err != nil {
        return nil, err
      }
      fieldCount = len(colNames)
    }
    values := make([]interface{}, fieldCount)
    targets := make([]interface{}, fieldCount)
    for i := 0; i < len(values); i++ {
      targets[i] = &values[i]
    }
    if err := rr.Scan(targets...); err != nil {
      return nil, err
    }
    // Build a map so we can access fields by name.
    m := make(map[string]interface{})
    for i := 0; i < len(values); i++ {
      // We assume values of type []byte are strings.
      switch v := values[i].(type) {
      case []byte:
        values[i] = string(v)
      }
      m[colNames[i]] = values[i]
    }
    log.Printf("row is %+v", values)
    data = append(data, m)
  }
  return data, nil
}

// Row gets exactly one row from our database. The first arg is the query string, and all
// following args are passed to the Query function as arguments for the query string.
// If the database returns either zero rows or two or more rows, this function returns an error.
func (s *SqlSource) Row(args ...interface{}) (interface{}, error) {
  data, err := s.Rows(args...)
  if err != nil {
    return nil, err
  }
  rows, ok := data.([]map[string]interface{})
  if !ok {
    return nil, fmt.Errorf("SqlSource.Row unexpected type returned from Rows")
  }
  if len(rows) != 1 {
    return nil, fmt.Errorf("SqlSource.Row expected one row, got %d", len(rows))
  }
  return rows[0], nil
}
