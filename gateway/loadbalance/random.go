package loadbalance

import (
	"errors"
	"math/rand"
)

// Random 随机算法实现
type Random struct {
	// 当前的索引值
	curIndex int
	addrs    []string
}

func (r *Random) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params was Empty")
	}
	r.addrs = append(r.addrs, params...)
	return nil
}

func (r *Random) Next() string {
	if len(r.addrs) == 0 {
		return ""
	}
	// 随机返回一个服务器
	addr := r.addrs[rand.Intn(len(r.addrs))]
	return addr
}
