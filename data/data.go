package data

type Source interface {
  Data(args ...interface{}) (interface{}, error)
}

type EmptySource struct{}

func (s *EmptySource) Data(args ...interface{}) (interface{}, error) {
  return nil, nil
}
