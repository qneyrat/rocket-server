package jsonrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"rocket-server/middleware"
)

type Request struct {
	Method  string `json:"method"`
	Params  string `json:"params"`
	ID      string `json:"id"`
}

type Response struct {
	Result  string `json:"result"`
	Error   string `json:"error,omitempty"`
	ID      string `json:"id,omitempty"`
}

type Server struct {
	next middleware.Middleware

	methods map[string]Method
}

func NewServer() *Server {
	return &Server{
		methods: make(map[string]Method),
	}
}

func (s *Server) AddMiddleware(next middleware.Middleware) {
	s.next = next
}

type Method func(req Request) Response

func (s *Server) AddMethod(name string, method Method) {
	s.methods[name] = method
}

func (s *Server) getMethod(name string) (Method, error) {
	e, ok := s.methods[name]
	if !ok {
		return nil, errors.New("method not allowed")
	}
	return e, nil
}

func (s *Server) Handle(w *bytes.Buffer, r *bytes.Buffer) {
	var j Request
	if r == nil {
		return
	}
	err := json.NewDecoder(r).Decode(&j)
	if err != nil {
		WriteError(w, err)
		return
	}

	e, err := s.getMethod(j.Method)
	if err != nil {
		WriteError(w, err)
		return
	}
	res := e(j)
	b, err := json.Marshal(res)
	if err != nil {
		WriteError(w, err)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		WriteError(w, err)
		return
	}

	if s.next != nil {
		s.next.Handle(w, r)
	}
}

func WriteError(buf *bytes.Buffer, err error)  {
	buf.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err)))
}
