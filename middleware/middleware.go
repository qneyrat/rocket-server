package middleware

import (
	"bytes"
	"fmt"
)

type Middleware interface {
	Handle(*bytes.Buffer, *bytes.Buffer)
}

type ResetMiddleware struct {}

func (ResetMiddleware) Handle(buf *bytes.Buffer, r *bytes.Buffer)  {
	fmt.Println(buf.String())
	buf.Reset()
	buf.Write([]byte(`{}`))
}
