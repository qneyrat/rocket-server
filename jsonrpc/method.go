package jsonrpc

import (
	"rocket-server/service"
)

func NewSayHelloMethod(s *service.Service) Method {
	return func(req Request) Response {
		return Response{
			Result:  s.SayHello("you"),
			ID:      req.ID,
		}
	}
}
