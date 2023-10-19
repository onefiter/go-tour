package queue

import (
	"context"
	"sync"
)

type ConcurrentBlockingQueue[T any] struct {
	mutex    *sync.Mutex
	data     []T
	notFull  chan struct{}
	notEmpty chan struct{}
	maxSize  int
}

func NewConcurrentBlockingQueue[T any](maxSize int) *ConcurrentBlockingQueue[T] {
	m := &sync.Mutex{}
	return &ConcurrentBlockingQueue[T]{
		data:     make([]T, 0, maxSize),
		mutex:    m,
		notFull:  make(chan struct{}, 1),
		notEmpty: make(chan struct{}, 1),
		maxSize:  maxSize,
	}
}

func (c *ConcurrentBlockingQueue[T]) Enqueue(ctx context.Context, data T) error {

	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mutex.Lock()

	for c.IsFull() {
		// 阻塞住:我阻塞我自己，直到有人唤醒我
		c.mutex.Unlock()

		select { // 再次加锁，控制住超时
		case <-c.notFull:
			c.mutex.Lock()
		case <-ctx.Done():
			return ctx.Err()

		}
	}
	// 1. 缺点append会导致数组不断扩容，出队列没能解决缓存--可以用两个指针进行维护
	// 2.
	c.data = append(c.data, data)
	if len(c.data) == 1 {
		// 只有从空变不空才会发信号
		c.notEmpty <- struct{}{}
	}
	c.mutex.Unlock()

	return nil
}

func (c *ConcurrentBlockingQueue[T]) Dequeue(ctx context.Context) (T, error) {

	if ctx.Err() != nil {
		var t T
		return t, ctx.Err()
	}
	c.mutex.Lock()

	for c.IsEmpty() {
		// 阻塞住我自己，等待有元素入队
		// 一睡不醒，超时也不知道
		c.mutex.Unlock()
		select {
		case <-c.notEmpty:
			c.mutex.Lock()
		case <-ctx.Done():
			var t T
			return t, ctx.Err()
		}

	}
	// 队首
	// [a, b, c, d]
	t := c.data[0]
	c.data = c.data[1:]
	if len(c.data) == c.maxSize-1 {
		// 只有从满变不满，才发信号
		c.notFull <- struct{}{}
	}
	c.mutex.Unlock()

	return t, nil

}

func (c *ConcurrentBlockingQueue[T]) IsFull() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.data) == c.maxSize
}

func (c *ConcurrentBlockingQueue[T]) IsEmpty() bool {
	return len(c.data) == 0
}

func (c *ConcurrentBlockingQueue[T]) Len() uint64 {
	return uint64(len(c.data))
}
