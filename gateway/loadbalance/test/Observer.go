package test

// Observer 观察者接口
// 多个观察者对应一个被观察者,想相当于多个监听者对应一个被监听者
// 观察者通过Update方法更新自己的可用服务列表
type Observer interface {
	Update()
}

// ConcreteSubject 观察者具体主体
type ConcreteSubject struct {
	observers []Observer // 观察者列表
	conf      []string   // 配置信息，即数据
	name      string     // 主体名称
}

// NewConcreteSubject 根据指定名称返回一个具体主体实例
func NewConcreteSubject(name string) (*ConcreteSubject, error) {
	mConf := &ConcreteSubject{name: name}
	return mConf, nil
}

// Attach 将新的观察者添加到观察者列表
func (s *ConcreteSubject) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

// NotifyAllObservers 通知所有观察者
func (s *ConcreteSubject) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

// GetConf 获取具体主题的配置信息
func (s *ConcreteSubject) GetConf() []string {
	return s.conf
}

// UpdateConf 更新配置消息时,通知所有观察者也更新
// 更新配置信息意味着发送了事件,所以要通知观察者
func (s *ConcreteSubject) UpdateConf(conf []string) {
	s.conf = conf
	//for _, obs := range s.observers {
	//	obs.Update()
	//}
	s.NotifyAllObservers()
}
