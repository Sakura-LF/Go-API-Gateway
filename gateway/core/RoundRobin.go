package core

import (
	"errors"
	"fmt"
)

// RoundRobin 轮询算法实现
type RoundRobinBalance struct {
	// 当前的索引值
	CurIndex    int
	ServerAddrs []string

	// 观察主体
	Conf Subject
}

func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params was Empty")
	}
	r.ServerAddrs = append(r.ServerAddrs, params...)
	return nil
}

func (r *RoundRobinBalance) Next() string {
	if len(r.ServerAddrs) == 0 {
		return ""
	}
	addr := r.ServerAddrs[r.CurIndex]
	// 使用队列中的算法
	r.CurIndex = (r.CurIndex + 1) % len(r.ServerAddrs)
	return addr
}

func (r *RoundRobinBalance) Update() {
	// 从具体发布中将conf传入ServerAddrs
	if conf, ok := r.Conf.(*LoadBalanceConsul); ok {
		for _, items := range conf.ConfigIpWeight {
			for _, serviceItems := range items {
				r.ServerAddrs = append(r.ServerAddrs, serviceItems.Host)
			}
		}
	}
	fmt.Println(r)
}

func (r *RoundRobinBalance) SetConf(conf Subject) {
	r.Conf = conf
}

func (r *RoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}
