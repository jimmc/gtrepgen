package dbsource

import (
  "database/sql"
  "errors"
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

// Data gets data from our database. The first arg is the query string, and all
// following args are passed to the Query function as arguments for the query string.
func (s *SqlSource) Data(args ...interface{}) (interface{}, error) {
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
