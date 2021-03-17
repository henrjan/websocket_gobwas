package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func main() {
	go websocket()

	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {

		u := fmt.Sprintf("http://127.0.0.1:8080?%s", r.URL.RawQuery)
		newUrl, err := url.Parse(u)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		reqProxy := httputil.NewSingleHostReverseProxy(newUrl)
		reqProxy.ServeHTTP(rw, r)

	})

	http.ListenAndServe(":3300", nil)
}

func websocket() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			fmt.Println(err)
			return
		}

		go func() {
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)

				if err != nil {
					fmt.Println(err)
					return
				}
				err = wsutil.WriteServerMessage(conn, op, msg)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}()
	}))
}
