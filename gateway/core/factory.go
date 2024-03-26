package core

type LoadBalanceType int

const (
	RoundRobin LoadBalanceType = iota
	Random
	WeightRoundRobin
	ConsistenHash
)

type LoadBalance interface {
	Add(...string) error
	Next() string

	Get(string) (string, error)

	//后期服务发现补充
	Update()
}

// LoadBalanceFactory m Random or RoundRobin or WeightRoundRobin or ConsistenHash
func LoadBalanceFactory(balanceType LoadBalanceType, consul Subject) LoadBalance {
	switch balanceType {
	case RoundRobin:
		lb := &RoundRobinBalance{}
		lb.SetConf(consul)
		consul.Attach(lb)
		lb.Update()
		return lb
	case Random:
		lb := &RandomBalance{}
		lb.SetConf(consul)
		consul.Attach(lb)
		lb.Update()
		return lb
	case WeightRoundRobin:
		lb := &RoundRobinBalance{}
		lb.SetConf(consul)
		consul.Attach(lb)
		lb.Update()
		return lb
	case ConsistenHash:
		lb := &RoundRobinBalance{}
		lb.SetConf(consul)
		consul.Attach(lb)
		lb.Update()
		return lb
	}
	return nil
}
