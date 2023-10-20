//go:build v1

package queue

import (
	"context"
	"sync"
)

type ConcurrentBlockingQueue[T any] struct {
	mutex *sync.Mutex
	data  []T
	//notFull  chan struct{}
	//notEmpty chan struct{}
	maxSize int

	notEmptyCond *cond
	notFullCond  *cond
}

func NewConcurrentBlockingQueue[T any](maxSize int) *ConcurrentBlockingQueue[T] {
	m := &sync.Mutex{}
	return &ConcurrentBlockingQueue[T]{
		data:  make([]T, 0, maxSize),
		mutex: m,
		//notFull:  make(chan struct{}, 1),
		//notEmpty: make(chan struct{}, 1),
		maxSize: maxSize,
		notFullCond: &cond{
			sync.NewCond(m),
		},
		notEmptyCond: &cond{
			sync.NewCond(m),
		},
	}
}

func (c *ConcurrentBlockingQueue[T]) Enqueue(ctx context.Context, data T) error {

	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mutex.Lock()

	for c.isFull() {
		err := c.notFullCond.WaitTimeout(ctx)
		if err != nil {
			return err
		}
	}
	// 1. 缺点append会导致数组不断扩容，出队列没能解决缓存--可以用两个指针进行维护
	// 2.
	c.data = append(c.data, data)
	c.notFullCond.Signal()
	c.mutex.Unlock()

	return nil
}

func (c *ConcurrentBlockingQueue[T]) Dequeue(ctx context.Context) (T, error) {

	if ctx.Err() != nil {
		var t T
		return t, ctx.Err()
	}
	c.mutex.Lock()

	for c.isEmpty() {
		err := c.notEmptyCond.WaitTimeout(ctx)
		if err != nil {
			var t T
			return t, err
		}
	}
	// 队首
	// [a, b, c, d]
	t := c.data[0]
	c.data = c.data[1:]
	c.notEmptyCond.Signal()
	c.mutex.Unlock()

	return t, nil

}

func (c *ConcurrentBlockingQueue[T]) IsFull() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.isFull()
}

func (c *ConcurrentBlockingQueue[T]) isFull() bool {

	return len(c.data) == c.maxSize
}

func (c *ConcurrentBlockingQueue[T]) IsEmpty() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.isEmpty()
}

func (c *ConcurrentBlockingQueue[T]) isEmpty() bool {
	return len(c.data) == 0
}

func (c *ConcurrentBlockingQueue[T]) Len() uint64 {
	return uint64(len(c.data))
}

type cond struct {
	*sync.Cond
}

func (c *cond) WaitTimeout(ctx context.Context) error {
	ch := make(chan struct{})

	go func() {
		c.Cond.Wait()
		select {
		case ch <- struct{}{}:
		default:
			// 这里已经超时
			c.Cond.Signal()
			c.Cond.L.Unlock()
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		// 真的被唤醒
		return nil
	}

}
