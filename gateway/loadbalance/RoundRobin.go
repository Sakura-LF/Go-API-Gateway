package loadbalance

import (
	"errors"
)

// RoundRobin 轮询算法实现
type RoundRobin struct {
	// 当前的索引值
	curIndex    int
	ServerAddrs []string
}

func (r *RoundRobin) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params was Empty")
	}
	r.ServerAddrs = append(r.ServerAddrs, params...)
	return nil
}

func (r *RoundRobin) Next() string {
	if len(r.ServerAddrs) == 0 {
		return ""
	}
	addr := r.ServerAddrs[r.curIndex]
	// 使用队列中的算法
	r.curIndex = (r.curIndex + 1) % len(r.ServerAddrs)
	return addr
}
