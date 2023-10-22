package queue

import "context"

type Queue[T any] interface {
	// 定义方法
	// 入队和出队的方法
	// Enqueue()
	// Dequeue()
	// 其他语言会这样设计
	// Enqueue(duration time.Duration,  data any) error
	// Enqueue Go语言中用context.Context 来设计超时控制

	EnQueue(ctx context.Context, data T) error
	DeQueue(ctx context.Context) (T, error)

	IsFull() bool
	IsEmpty() bool
	Len() uint64
}
