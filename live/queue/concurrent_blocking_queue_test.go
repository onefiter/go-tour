package queue

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestConcurrentBlockingQueue_EnQueue(t *testing.T) {
	testCases := []struct {
		name    string
		q       *ConcurrentBlockingQueue[int]
		timeout time.Duration
		value   int
		data    []int
		wantErr error
	}{
		{
			name:    "入队成功",
			q:       NewConcurrentBlockingQueue[int](10),
			value:   1,
			timeout: time.Minute,
			data:    []int{1},
		},
		{
			name: "blocking and timeout",
			q: func() *ConcurrentBlockingQueue[int] {
				res := NewConcurrentBlockingQueue[int](2)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				err := res.Enqueue(ctx, 1)
				require.NoError(t, err)
				err = res.Enqueue(ctx, 2)
				require.NoError(t, err)
				return res
			}(),
			value:   3,
			timeout: time.Second,
			data:    []int{1, 2},
			wantErr: context.DeadlineExceeded,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			err := tc.q.Enqueue(ctx, tc.value)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.data, tc.q.data)
		})
	}
}

func TestConcurrentBlockingQueue(t *testing.T) {
	// 要测试什么场景
	// 1. 从goroutine数量考虑
	//  a. enqueue和dequeue数量一样
	//  b. enqueue > dequeue
	//  c. enqueue < dequeue
	// 加锁的话，只能确保没有死锁
	q := NewConcurrentBlockingQueue[int](10000)
	for i := 0; i < 20; i++ {
		go func() {
			for {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				// 没有办法校验这里面的中间结果
				_ = q.Enqueue(ctx, rand.Int())
				// 怎么断言error
				cancel()

			}
		}()
	}

	for i := 0; i < 20; i++ {
		go func() {
			for {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)

				_, _ = q.Dequeue(ctx)
				// 怎么断言 val ,err
				cancel()

			}
		}()
	}
	// 怎么校验q对不对
}

func BenchmarkConcurrentQueue(b *testing.B) {
	var wg sync.WaitGroup

	q := NewConcurrentBlockingQueue[int](100)

	wg.Add(2)

	go func() {
		for i := 0; i < b.N; i++ {
			_ = q.Enqueue(context.Background(), i)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < b.N; i++ {
			_, _ = q.Dequeue(context.Background())

		}
		wg.Done()
	}()

	wg.Wait()
}
