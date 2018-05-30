package main

import (
	"net/http"

	"rocket-server/crypto"
	"rocket-server/jsonrpc"
	"rocket-server/service"
	"rocket-server/ws"
)

func main() {
	s := jsonrpc.NewServer()
	s.AddMethod("hello", jsonrpc.NewSayHelloMethod(&service.Service{}))

	c := crypto.Crypto{Key:[]byte("0123456789012345")}
	http.ListenAndServe(":4001", ws.NewHandlerFunc(s, c))
}
