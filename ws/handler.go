package ws

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"rocket-server/middleware"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func NewHandlerFunc(m middleware.Middleware) http.HandlerFunc {
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

				bufw := new(bytes.Buffer)
				bufr := new(bytes.Buffer)

				bufr.ReadFrom(r)
				m.Handle(bufw, bufr)
				io.Copy(w, bufw)

				if err = w.Flush(); err != nil {
					break
				}
			}
		}()
	}
}
