package queue

import (
	"log"
	"sync/atomic"
	"testing"
)

func TestCAS(t *testing.T) {
	var value int64 = 10
	// 我准备把 value 更新为 12，当且仅当 value 原本的值是 10
	res := atomic.CompareAndSwapInt64(&value, 10, 12)

	log.Println(res)
	log.Println(value)
}
