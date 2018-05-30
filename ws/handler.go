package ws

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"rocket-server/crypto"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Middleware interface {
	Handle(r io.Reader, w io.Writer)
}

type Message struct {
	Body string `json:"body"`
}

func NewHandlerFunc(h Middleware, c crypto.Crypto) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w, nil)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer conn.Close()

			var (
				state  = ws.StateServerSide
				r = wsutil.NewReader(conn, state)
				w = wsutil.NewWriter(conn, state, ws.OpText)
			)

			for {
				header, err := r.NextFrame()
				if err != nil {
					log.Println(err)
					break
				}

				w.Reset(conn, state, header.OpCode)

				var reqMessage Message
				err = json.NewDecoder(r).Decode(&reqMessage)
				if err != nil {
					break
				}

				clearMessage, err := c.Decrypt([]byte(reqMessage.Body))
				if err != nil {
					break
				}

				req := bytes.NewBuffer(clearMessage)
				rw := new(bytes.Buffer)
				h.Handle(req, rw)

				encMessage, err := c.Encrypt(rw.Bytes())
				if err != nil {
					break
				}

				resMessage := Message{Body: string(encMessage)}
				err = json.NewEncoder(w).Encode(&resMessage)
				if err != nil {
					break
				}

				if err = w.Flush(); err != nil {
					break
				}
			}
		}()
	}
}
