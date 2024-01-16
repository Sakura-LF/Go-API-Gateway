package core

import (
	"time"
)

var ConcreteSubjectConsul *LoadBalanceConsul

// Subject 抽象主体(发布者)
type Subject interface {
	// Attach 增加观察者
	Attach(Observer)

	// Detach 减少观察者
	//Detach(Observer)

	// Notify 通知观察者
	Notify()
}

// LoadBalanceConsul (ConcreteSubject) 具体主体(发布者)
type LoadBalanceConsul struct {
	Observers      []Observer                // 观察者列表
	ConfigIpWeight map[string][]ServiceItems // IP 和权重的映射表
}

// NewConcreteSubject 根据指定名称返回一个具体主体(观察者)实例
func NewConcreteSubject() {
	ConcreteSubjectConsul = &LoadBalanceConsul{
		ConfigIpWeight: make(map[string][]ServiceItems),
	}
	// 从Service中得到地址
	ConcreteSubjectConsul.WatchConf()
	//return ConcreteSubjectConsul, nil
}

// Attach 将新的观察者添加到观察者列表
func (s *LoadBalanceConsul) Attach(observer Observer) {
	s.Observers = append(s.Observers, observer)
}

// Notify  通知所有观察者
func (s *LoadBalanceConsul) Notify() {
	for _, concreteObserver := range s.Observers {
		concreteObserver.Update()
	}
}

// Attach 将新的观察者添加到观察者列表
func (s *LoadBalanceConsul) WatchConf() {
	GetService()
	go func() {
		for {
			for key, value := range Service.Data {
				s.ConfigIpWeight[key] = value
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

// Observer 抽象观察者
type Observer interface {
	// Update 更新
	Update()
}

//// ConcreteObserver 具体观察者
//type ConcreteObserver struct {
//	subject *LoadBalanceConsul
//}
//
//func NewConcreteObserver(subject *LoadBalanceConsul) (*ConcreteObserver, error) {
//	concreteObserver := &ConcreteObserver{subject: subject}
//	return concreteObserver, nil
//}
