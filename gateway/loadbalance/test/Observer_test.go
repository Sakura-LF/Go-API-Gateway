package test

//func TestObserver(t *testing.T) {
//	// 创建发布者
//	subject, err := loadbalance.NewConcreteSubject("Consul")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	// 创建观察者,并将观察者绑定到发布者
//	observer1, err := loadbalance.NewConcreteObserver(subject)
//	observer2, err := loadbalance.NewConcreteObserver(subject)
//
//	//subject := &ConcreteSubject{}
//	//observer1 := &ConcreteObserver{subject}
//	//observer2 := &ConcreteObserver{subject}
//
//	subject.Attach(observer1)
//	subject.Attach(observer2)
//
//	subject.conf = []string{"192.168.2.7:8080", "192.168.2.7:8081"}
//
//	subject.Notify()
//}
