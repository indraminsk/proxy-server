package handler

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// struct to manage worker's requests
type ServerType struct {
}

// all action to processing worker request:
// check request, send response
func (server *ServerType) HandlerWorkerRequest(w http.ResponseWriter, r *http.Request) {
	var (
		err error

		sleepTime int
	)

	defer func() {
		if accident := recover(); accident != nil {
			fmt.Println()
			fmt.Println("[recover] handler worker accident:", accident)

			http.Error(w, "accident", http.StatusInternalServerError)
		}
	}()

	fmt.Println()

	// prepare response
	switch r.Method {
	case http.MethodGet:
		sleepTime = rand.Intn(14) + 1

		fmt.Println("request:", r.Header.Get("ID"))
		fmt.Println("request:", r.Header.Get("Worker-Url"))
		fmt.Println("sleep:", sleepTime)

		time.Sleep(time.Duration(sleepTime) * time.Second)

		err = sendWorkerRequest(r, sleepTime)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err = w.Write([]byte(fmt.Sprintf("wrong method: %s", r.Method)))
	}

	if err != nil {
		fmt.Println("[error] build response:", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}

func sendWorkerRequest(r *http.Request, sleepTime int) (err error) {
	var (
		client   *http.Client
		request  *http.Request
		response *http.Response

		url  string
		body []byte
	)

	client = &http.Client{}

	url = fmt.Sprintf("http://%s/service/in", r.Header.Get("Worker-Url"))
	body = []byte(fmt.Sprintf("sleep: %ds", sleepTime))

	request, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header.Add("ID", r.Header.Get("ID"))
	request.Header.Add("Status", fmt.Sprintf("%d", http.StatusOK))

	response, err = client.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		if err = response.Body.Close(); err != nil {
			fmt.Println("[error] clear response memory:", err)
		}
	}()

	return nil
}
