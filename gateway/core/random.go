package core

import (
	"errors"
	"fmt"
	"math/rand"
)

// Random 随机算法实现
type RandomBalance struct {
	// 当前的索引值
	CurIndex int
	Addrs    []string
	// 观察主体
	Conf Subject
}

func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params was Empty")
	}
	r.Addrs = append(r.Addrs, params...)
	return nil
}

func (r *RandomBalance) Next() string {
	if len(r.Addrs) == 0 {
		return ""
	}
	// 随机返回一个服务器
	addr := r.Addrs[rand.Intn(len(r.Addrs))]
	return addr
}

func (r *RandomBalance) SetConf(conf Subject) {
	r.Conf = conf
}

func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RandomBalance) Update() {
	// 从具体发布中将conf传入ServerAddrs
	if conf, ok := r.Conf.(*LoadBalanceConsul); ok {
		for _, items := range conf.ConfigIpWeight {
			for _, serviceItems := range items {
				r.Addrs = append(r.Addrs, serviceItems.Host)
			}
		}
	}
	fmt.Println(r)
}
