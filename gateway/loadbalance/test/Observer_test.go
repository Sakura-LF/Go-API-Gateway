package test

import (
	"fmt"
	"testing"
)

func TestObserver(t *testing.T) {
	// 创建一个具体主体:被观察者
	obConf, err := NewConcreteSubject("rs_server")
	if err != nil {
		return
	}
	// 创建一个观察者,维护他要观察的主体
	loadBalanceObserver := &LoadBalanceObserver{
		Conf: obConf,
	}
	fmt.Println("loadBalanceObserver conf:", loadBalanceObserver.Conf.conf)

	// 注册监听,将观察者绑定到具体主体的观察者列表中
	obConf.Attach(loadBalanceObserver)

	// 具体主体数据更新,
	obConf.UpdateConf([]string{"127.0.0.1"})
}

// 具体观察者,观察者接口的实现
// 维护一个本地的具体主体(可用服务器列表):Conf
// 实现抽象观察者所要求的更新接口,通过接口方法更新主体状态
type LoadBalanceObserver struct {
	Conf *ConcreteSubject
}

func (L *LoadBalanceObserver) Update() {
	fmt.Println("update get conf", L.Conf.conf)
}
