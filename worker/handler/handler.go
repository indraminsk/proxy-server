// package handler release function to processing http requests
package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const DefaultRequestKey = "f229fcff-012d-4ad7-8b1a-88bf0ce2aba9"

// describe structure of log record
type ServerLogRecordType struct {
	Received time.Time
	ClientData ClientRequestType
	ServiceData ServiceResponseType
}

// struct to manage client's requests
// map key is request key (uuid)
type ServerType struct {
	Log map[string]*ServerLogRecordType
}

// data type for request/response headers
type HeadersType map[string]string

// data type for client request
type ClientRequestType struct {
	Method string
	Url string
	Headers HeadersType
}

// data type for client check request
type ClientCheckType struct {
	Request string
}

// data type for service request
type ServiceResponseType struct {
	Id string
	Status int
	Headers HeadersType
	Length int
}

// all action to processing client request:
// log request, check request, parse client json, send request to service
func (server *ServerType) HandlerClientRequest(w http.ResponseWriter, r *http.Request) {
	var (
		err error

		decoder *json.Decoder
		clientRequest ClientRequestType
		requestKey string
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
	requestKey = buildRequestKey(DefaultRequestKey)
	server.Log[requestKey] = &ServerLogRecordType{
		Received: time.Now(),
		ClientData:  clientRequest,
	}

	fmt.Println("===", requestKey, "===")
	fmt.Println()
	fmt.Println("[CLIENT REQUEST]")
	fmt.Println(server.Log[requestKey].Received.String())
	fmt.Println(server.Log[requestKey].ClientData)
	fmt.Println()

	// send to 3d-party service request



	// prepare response
	w.WriteHeader(http.StatusAccepted)

	_, err = w.Write([]byte(requestKey))
	if err != nil {
		fmt.Println("[error] build success response:", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}

func (server *ServerType) HandlerClientStatus(w http.ResponseWriter, r *http.Request) {
	var (
		err error

		decoder *json.Decoder
		clientCheck ClientCheckType
		jsonData []byte
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
	if server.Log[clientCheck.Request].ServiceData.Id == "" {
		fmt.Println("continue waiting response")

		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write([]byte(clientCheck.Request))
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

	if err != nil {
		fmt.Println("[error] build success response:", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
}

func (server *ServerType) HandlerServiceIn(w http.ResponseWriter, r *http.Request) {
	//server.Log[recordKey].ServiceData = serviceResponse
	//
	//fmt.Println("[SERVICE RESPONSE]")
	//fmt.Println("resource:", server.Log[recordKey].ServiceData)
	//fmt.Println()
}

// key is calculated uuid value
// in test aim we can use const DefaultRequestKey to get always the same key
func buildRequestKey(key string) string {
	if key != "" {
		return key
	}

	return uuid.New().String()
}