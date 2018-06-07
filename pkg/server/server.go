package server

import (
	"net/http"

	"rocket-server/internal/service"
	"rocket-server/internal/ws"
	"rocket-server/pkg/crypto"
	"rocket-server/pkg/jsonrpc"
)

func newSayHelloMethod(s *service.Service) jsonrpc.Method {
	return func(req jsonrpc.Request) jsonrpc.Response {
		return jsonrpc.Response{
			Result:  s.SayHello("you"),
			ID:      req.ID,
		}
	}
}

func Start() {
	s := jsonrpc.NewServer()
	s.AddMethod("hello", newSayHelloMethod(&service.Service{}))

	c := crypto.Crypto{Key: []byte("0123456789012345")}
	http.ListenAndServe(":4001", ws.NewHandlerFunc(s, c))
}
