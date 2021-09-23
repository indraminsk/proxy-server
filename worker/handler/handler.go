// package handler release function to processing http requests
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// describe structure of log record
type ServerLogRecordType struct {
	Received    time.Time
	ClientData  ClientDataType
	ServiceData ServiceResponseType
}

// struct to manage client's requests
// map key is request key (uuid)
type ServerType struct {
	IP  string
	Log map[string]*ServerLogRecordType
}

// data type for client request
type ClientDataType struct {
	Method  string
	Url     string
	Headers http.Header
}

// data type for client check request
type ClientStatusType struct {
	Request string
}

// data type for service request
type ServiceResponseType struct {
	Id      string
	Status  string
	Headers http.Header
	Length  string
}

// all action to processing client request:
// log request, check request, parse client json, send request to service
func (server *ServerType) HandlerClientRequest(w http.ResponseWriter, r *http.Request) {
	var (
		err error

		decoder       *json.Decoder
		clientRequest ClientDataType
		id            string
	)

	defer func() {
		if accident := recover(); accident != nil {
			fmt.Println()
			fmt.Println("[recover] handler client accident:", accident)

			http.Error(w, "accident", http.StatusInternalServerError)
		}
	}()

	fmt.Println()

	// we wait only POST request
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte("I'm ready to POST only"))
		if err != nil {
			fmt.Println("[error] processing wrong request type", err)
			http.Error(w, "error", http.StatusInternalServerError)
		}

		return
	}

	// convert request body to request structure
	decoder = json.NewDecoder(r.Body)
	err = decoder.Decode(&clientRequest)
	if err != nil {
		fmt.Println("[error] decode request params:", err)
		http.Error(w, "error", http.StatusInternalServerError)

		return
	}

	// log events
	id = uuid.New().String()
	server.Log[id] = &ServerLogRecordType{
		Received:   time.Now(),
		ClientData: clientRequest,
	}

	fmt.Println("===", id, "===")
	fmt.Println()

	err = validateClientRequest(server.Log[id].ClientData)
	if err != nil {
		fmt.Println("[error] validate client's data:", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	fmt.Println("[CLIENT REQUEST]")
	fmt.Println(server.Log[id].Received.String())
	fmt.Println(server.Log[id].ClientData)
	fmt.Println()

	// send to 3d-party service request
	go server.sendServiceRequest(id)

	// prepare response
	w.WriteHeader(http.StatusAccepted)

	_, err = w.Write([]byte(id))
	if err != nil {
		fmt.Println("[error] build success response:", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}

// the simplest validation client's json
func validateClientRequest(clientData ClientDataType) (err error) {
	if clientData.Url == "" {
		return errors.New("bad url")
	}

	if (clientData.Method != http.MethodPost) && (clientData.Method != http.MethodGet) {
		return errors.New("not allowed http method")
	}

	if len(clientData.Headers) == 0 {
		return errors.New("not set headers")
	}

	return err
}

func (server *ServerType) HandlerClientStatus(w http.ResponseWriter, r *http.Request) {
	var (
		err error

		decoder     *json.Decoder
		clientCheck ClientStatusType
		jsonData    []byte
	)

	defer func() {
		if accident := recover(); accident != nil {
			fmt.Println()
			fmt.Println("[recover] handler client accident:", accident)

			http.Error(w, "accident", http.StatusInternalServerError)
		}
	}()

	fmt.Println()

	// we wait only POST request
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte("I'm ready to POST only"))
		if err != nil {
			fmt.Println("[error] processing wrong request type", err)
			http.Error(w, "error", http.StatusInternalServerError)
		}

		return
	}

	// convert request body to request structure
	decoder = json.NewDecoder(r.Body)
	err = decoder.Decode(&clientCheck)
	if err != nil {
		fmt.Println("[error] decode request params:", err)
		http.Error(w, "error", http.StatusInternalServerError)

		return
	}

	fmt.Println("===", clientCheck.Request, "===")
	fmt.Println()
	fmt.Println("[STATUS REQUEST]")

	// check request status and prepare response
	if _, ok := server.Log[clientCheck.Request]; !ok {
		fmt.Println("request isn't registered")

		w.WriteHeader(http.StatusNotFound)
		_, err = w.Write([]byte("request isn't registered"))
	} else {
		if server.Log[clientCheck.Request].ServiceData.Id == "" {
			fmt.Println("continue waiting response")

			w.WriteHeader(http.StatusAccepted)
			_, err = w.Write([]byte("continue waiting response"))

			return
		} else {
			server.Log[clientCheck.Request].Received = time.Now()

			fmt.Println(server.Log[clientCheck.Request].Received.String())
			fmt.Println(server.Log[clientCheck.Request].ServiceData)

			w.WriteHeader(http.StatusOK)

			jsonData, err = json.Marshal(server.Log[clientCheck.Request].ServiceData)
			if err != nil {
				fmt.Println("[error] convert to json service response:", err)
				http.Error(w, "error", http.StatusInternalServerError)
				return
			}

			_, err = w.Write(jsonData)
		}
	}

	if err != nil {
		fmt.Println("[error] build success response:", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}

func (server *ServerType) HandlerServiceResponse(w http.ResponseWriter, r *http.Request) {
	var (
		err error

		id string
	)

	defer func() {
		if accident := recover(); accident != nil {
			fmt.Println()
			fmt.Println("[recover] handler client accident:", accident)

			http.Error(w, "accident", http.StatusInternalServerError)
		}
	}()

	fmt.Println()

	// we wait only POST request
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte("I'm ready to POST only"))
		if err != nil {
			fmt.Println("[error] processing wrong request type", err)
			http.Error(w, "error", http.StatusInternalServerError)
		}

		return
	}

	id = r.Header.Get("ID")

	server.Log[id].Received = time.Now()
	server.Log[id].ServiceData = ServiceResponseType{
		Id:      id,
		Status:  r.Header.Get("Status"),
		Headers: r.Header,
		Length:  r.Header.Get("Content-Length"),
	}

	fmt.Println("===", id, "===")
	fmt.Println()
	fmt.Println("[SERVICE RESPONSE]")
	fmt.Println(server.Log[id].Received.String())
	fmt.Println(server.Log[id].ServiceData)
	fmt.Println()
}

func (server *ServerType) sendServiceRequest(requestKey string) {
	var (
		err error

		clientData ClientDataType
		client     *http.Client
		request    *http.Request
		response   *http.Response
	)

	clientData = server.Log[requestKey].ClientData

	client = &http.Client{}

	switch clientData.Method {
	case http.MethodGet:
		request, err = http.NewRequest(clientData.Method, clientData.Url, nil)
		if err != nil {
			return
		}

		request.Header.Add("ID", requestKey)
		request.Header.Add("Worker-url", server.IP)

		for header, values := range clientData.Headers {
			for _, value := range values {
				request.Header.Add(header, value)
			}
		}

		response, err = client.Do(request)
		if err != nil {
			return
		}

		defer func() {
			if err = response.Body.Close(); err != nil {
				fmt.Println("[error] clear response memory:", err)
			}
		}()

		if response.StatusCode != http.StatusAccepted {
			return
		}

	case http.MethodPost:
	}
}
