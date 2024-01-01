package loadbalance

import (
	"fmt"
	"testing"
)

// 轮询算法实现
func TestRoundRobin_Add(t *testing.T) {
	robin := RoundRobin{}

	// 添加服务器
	robin.Add("127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003", "127.0.0.1:8004")

	// 模拟发送10个请求
	for i := 0; i < 10; i++ {
		fmt.Println(robin.Next())
	}
}

// 随机负载均衡算法测试
func TestRandom_Add(t *testing.T) {

	random := Random{}

	// 添加服务器
	random.Add("127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003", "127.0.0.1:8004")

	// 模拟发送10个请求
	for i := 0; i < 10; i++ {
		fmt.Println(random.Next())
	}
}

func TestConsistentHashBalance_Add(t *testing.T) {
	consistentHashBalance := NewConsistentHashBalance(2, nil)
	consistentHashBalance.Add("127.0.0.1:8000", "127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003", "127.0.0.1:8004")
	fmt.Println(consistentHashBalance.hashKeys)
	for key, value := range consistentHashBalance.hashMap {
		fmt.Println("Hash:", key, " Addr:", value)

	}
	fmt.Println("----------------URL----------------")
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8080/demo"))
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8081/demo"))
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8082/demo"))
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8083/demo"))
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8084/demo"))
	fmt.Println(consistentHashBalance.Get("http://127.0.0.0.1:8085/demo"))

	fmt.Println("----------------IP----------------")
	fmt.Println(consistentHashBalance.Get("0127.0.0.0.1:8000"))
	fmt.Println(consistentHashBalance.Get("127.0.0.0.1:8003"))
	fmt.Println(consistentHashBalance.Get("127.0.0.0.1:8002"))
	fmt.Println(consistentHashBalance.Get("127.0.0.0.1:8001"))

}
