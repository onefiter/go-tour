package demo

import "context"

type Selector1[T any] struct {
}

func (s *Selector1[T]) Build() (*Query, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector1[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector1[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
