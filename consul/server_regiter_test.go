package consul

import "testing"

func TestServerRegister(t *testing.T) {
	ServerRegister()
}

func TestConsulkv(t *testing.T) {
	Consulkv()
}

func TestKv(t *testing.T) {
	Kv()
}

func TestServerTest(t *testing.T) {
	ServerTest()
}

func TestTestHTTPServer(t *testing.T) {
	TestHTTPServer()
}

func TestSearchService(t *testing.T) {
	SearchService("Test")
}
