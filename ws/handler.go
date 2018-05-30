package ws

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"rocket-server/crypto"
	"rocket-server/middleware"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func NewHandlerFunc(m middleware.Middleware, c crypto.Crypto) http.HandlerFunc {
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

				var encReq []byte
				_, err = r.Read(encReq)
				if err != nil {
					break
				}
				message, err := c.Decrypt(encReq)
				if err != nil {
					break
				}

				buf := new(bytes.Buffer)
				buf.Write(message)

				m.Handle(buf)

				res := buf.Bytes()
				encRes, err := c.Encrypt(res)
				if err != nil {
					break
				}

				buf.Reset()
				buf.Write(encRes)

				io.Copy(w, buf)
				if err = w.Flush(); err != nil {
					break
				}
			}
		}()
	}
}
