package data

type Source interface {
  Data(args ...string) interface{}
}

type EmptySource struct{}

func (s *EmptySource) Data(args ...string) interface{} {
  return nil
}
