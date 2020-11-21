package discovery

import (
	// "sync"
	"sync/atomic"
)

const (
	name = "ROUND_BALANCER"
)

var (
	count     uint64
	maxUint64 = ^uint64(0)
	// pool      sync.Pool
)

func GetModVal(modVal uint64) uint64 {
	if count > maxUint64 {
		atomic.StoreUint64(&count, 0)
	}
	return atomic.AddUint64(&count, 1) % modVal

	// if val, ok := pool.Get().(uint64); ok {
	// 	return val % modVal
	// }
	// return 0
}
