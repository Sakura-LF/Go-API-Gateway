package rpc

import "testing"

func TestHelloWorld_HelloWorld(t *testing.T) {
	RpcServer()
}

func TestRpcClient(t *testing.T) {
	RpcClient()
}
