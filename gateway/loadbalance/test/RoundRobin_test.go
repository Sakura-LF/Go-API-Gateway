package test

import (
	"Go-API-Gateway/gateway/loadbalance"
	"fmt"
	"testing"
)

// 轮询算法实现
//func TestRoundRobin_Add(t *testing.T) {
//	robin := core.RoundRobin{}
//
//	// 添加服务器
//	robin.Add("127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003", "127.0.0.1:8004")
//
//	// 模拟发送10个请求
//	for i := 0; i < 10; i++ {
//		fmt.Println(robin.Next())
//	}
//}

// 随机负载均衡算法测试
//func TestRandom_Add(t *testing.T) {
//
//	random := core.Random{}
//
//	// 添加服务器
//	random.Add("127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003", "127.0.0.1:8004")
//
//	// 模拟发送10个请求
//	for i := 0; i < 10; i++ {
//		fmt.Println(random.Next())
//	}
//}

func TestConsistentHashBalance_Add(t *testing.T) {
	consistentHashBalance := loadbalance.NewConsistentHashBalance(2, nil)
	consistentHashBalance.Add("127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003", "127.0.0.1:8004")
	fmt.Println(consistentHashBalance.HashKeys)
	for key, value := range consistentHashBalance.HashKeys {
		fmt.Println("Hash:", key, " Addr:", value)

	}
	fmt.Println("----------------URL----------------")
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8080/demo"))
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8080/index"))
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8080/demo"))
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8080/demo"))
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8080/demo"))
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8080/demo"))

	fmt.Println("----------------IP----------------")
	fmt.Println(consistentHashBalance.Get("127.0.0.0.1:8000"))
	fmt.Println(consistentHashBalance.Get("127.0.0.0.1:8003"))
	fmt.Println(consistentHashBalance.Get("127.0.0.0.1:8002"))
	fmt.Println(consistentHashBalance.Get("127.0.0.0.1:8001"))

}

//func TestRoundRobin_Next(t *testing.T) {
//	weightRoundRobin := core.WeightRoundRobin{}
//	weightRoundRobin.Add("127.0.0.1:8080", "6")
//	weightRoundRobin.Add("127.0.0.1:8081", "3")
//	weightRoundRobin.Add("127.0.0.1:8082", "1")
//
//	one := 0
//	two := 0
//	three := 0
//	for i := 0; i < 20; i++ {
//		addr := weightRoundRobin.Next()
//		fmt.Println(addr)
//
//		if addr[len(addr)-1:len(addr)] == "0" {
//			one += 1
//		} else if addr[len(addr)-1:len(addr)] == "1" {
//			two += 1
//		} else if addr[len(addr)-1:len(addr)] == "2" {
//			three += 1
//		}
//	}
//	fmt.Println("8080: ", one)
//	fmt.Println("8081: ", two)
//	fmt.Println("8082: ", three)
//}
