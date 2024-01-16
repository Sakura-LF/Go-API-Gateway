package loadbalance

import (
	"Go-API-Gateway/gateway/core"
	"strconv"
)

// RoundRobin 轮询算法实现
type WeightRoundRobinBalance struct {
	// 当前的索引值
	curIndex    int
	ServerAddrs []*WeightServer

	// 观察主体
	conf core.Subject
}

type WeightServer struct {
	// 服务器地址
	addr string
	// 服务器权重
	weight int
	// 服务器当前权重
	currentweight int
	// 服务器有效权重
	effective int
}

func (weightRoundRobin *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		panic("params need addr:weight")
	}
	weight, err := strconv.Atoi(params[1])
	if err != nil {
		panic("the second params need weight")
	}
	// 构建node
	server := &WeightServer{
		addr:          params[0],
		weight:        weight,
		currentweight: 0,
		effective:     weight,
	}
	weightRoundRobin.ServerAddrs = append(weightRoundRobin.ServerAddrs, server)
	return nil
}

func (weightRoundRobin *WeightRoundRobinBalance) Next() string {
	// 所有节点的有效权重之和
	effectiveTotal := 0

	// 1.maxWeightServer,记录本轮最大权重的服务器
	var maxWeightServer *WeightServer

	for _, weightServer := range weightRoundRobin.ServerAddrs {
		// 当前权重+有效权重
		weightServer.currentweight += weightServer.effective
		// 如果当前遍历的server大于定义的server
		if maxWeightServer == nil || weightServer.currentweight > maxWeightServer.currentweight {
			maxWeightServer = weightServer
		}
		// 3.记录所有有效权重之后
		effectiveTotal += weightServer.effective
	}
	// 4.降权,本次最大的节点降权
	maxWeightServer.currentweight -= effectiveTotal

	return maxWeightServer.addr
}

// Next 通过每轮降权，权值较大的服务器，被选中的次数较多，实现了按权重访问的逻辑
//
// 带权服务器轮询方式核心逻辑：
//
//	循环计算每个服务器的权值（currentWeight），选择最大的返回
//	对选中的服务器进行降权：
//		currentWeight - 本轮所有有效权重之和
//
// 实现步骤：
//
//	1.定义变量maxNode，记录本轮权值最大的服务器
//	2.循环计算每个服务器的权重：临时权重 + 有效权重，选择最大的临时权重节点
//	3.记录所有有效权重之和：effectiveTotal
//	4.对选中节点进行降权

//func (weightRoundRobin *WeightRoundRobin) Callback(addr string, flag bool) {
//	for i := 0; i < len(weightRoundRobin.ServerAddrs); i++ {
//		w := weightRoundRobin.ServerAddrs[i]
//		if w.weightRoundRobin == addr {
//			// 访问服务器成功
//			if flag {
//				// 有效权重默认与权重相同，通讯异常 -1，正常 +1，不能超过weight大小
//				if w.effectiveWeight < w.weight {
//					w.effectiveWeight++
//				}
//			} else {
//				// 访问服务器失败
//				w.effectiveWeight--
//
//				// 刷新错误记录，把过去超过 failTimeout 的错误记录删除
//				refreshErrRecords(w)
//
//				// 记录本次失败
//				w.failTimes = append(w.failTimes, time.Now())
//				// 当前节点剩余失败次数，可能小于0
//				w.maxFails = maxFails - len(w.failTimes)
//			}
//			break
//		}
//	}
//}
