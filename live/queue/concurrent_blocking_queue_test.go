package queue

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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