package data

type Source interface {
  Row(args ...interface{}) (interface{}, error)
  Rows(args ...interface{}) (interface{}, error)
}

type EmptySource struct{}

func (s *EmptySource) Row(args ...interface{}) (interface{}, error) {
  return nil, nil
}

func (s *EmptySource) Rows(args ...interface{}) (interface{}, error) {
  return nil, nil
}
