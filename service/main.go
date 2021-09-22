package main

import (
	"axxon/service/handler"
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

	port = flag.Int("p", 9081, "service port")
	flag.Parse()

	fmt.Println("service run on port", *port)
	fmt.Println("to stop the service, press [Ctrl+C]")

	http.HandleFunc("/", server.HandlerWorkerRequest)

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Println("error:", err)
	}
}
