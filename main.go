package main

import (
	"net/http"

	"rocket-server/jsonrpc"
	"rocket-server/service"
	"rocket-server/ws"
)

func main() {
	s := jsonrpc.NewServer()
	s.AddMethod("hello", jsonrpc.NewSayHelloMethod(&service.Service{}))
	http.ListenAndServe(":4001", ws.NewHandlerFunc(s))
}
