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
