// package main run simple http server
package main

import (
	"axxon/worker/handler"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	var (
		err error

		port   *int
		server handler.ServerType
	)

	// processing unexpected panic
	defer func() {
		if accident := recover(); accident != nil {
			fmt.Println("[recover] main accident:", accident)
		}
	}()

	port = flag.Int("p", 9080, "service port")
	flag.Parse()

	fmt.Println("service run on port", *port)
	fmt.Println("to stop the service, press [Ctrl+C]")

	server.Log = make(map[string]*handler.ServerLogRecordType)

	http.HandleFunc("/client/request", server.HandlerClientRequest)
	http.HandleFunc("/client/status", server.HandlerClientStatus)
	http.HandleFunc("/service/in", server.HandlerServiceIn)

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Println("error:", err)
	}
}
