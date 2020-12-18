package config

import (
	"errors"
	"strconv"

	"github.com/qiwenilli/auv.kit/discovery"
)

// 如果存在则返回error
func Put(k, v string) error {
	if Get(k) != "" {
		return errors.New("key exist")
	}
	return Set(k, v)
}

func Set(k, v string) error {
	return discovery.KvPut(k, v)
}

func Get(k string) string {
	return discovery.KvGet(k)
}

func String(k string) string {
	return Get(k)
}

func Float64(k string) float64 {
	v, _ := strconv.ParseFloat(Get(k), 64)
	return v
}

func UInt64(k string) uint64 {
	v, _ := strconv.ParseUint(Get(k), 10, 64)
	return v
}

func Int64(k string) int64 {
	v, _ := strconv.ParseInt(Get(k), 10, 64)
	return v
}
