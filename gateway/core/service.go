package core

var Service *GateWayService
var Router *GateRouter

func init() {
	Service = &GateWayService{data: make(map[string][]ServiceItems, 10)}
	//Items = make([]ServiceItems, 0)
	Router = &GateRouter{data: make(map[string][]RouterItems, 10)}
}

type GateWayService struct {
	data map[string][]ServiceItems
}

type ServiceItems struct {
	Id        string
	Name      string
	CreatedAt int64
	UpdatedAt int64
	Host      string
	Protocol  string
}

type GateRouter struct {
	data map[string][]RouterItems
}

type RouterItems struct {
	Id        string
	CreatedAt int64
	UpdatedAt int64
	Path      []string
	Protocol  string
}
