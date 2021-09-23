// test for handler package. this is not full lists of tests of course. at least we need check others params in client's
// json. also useful make tests to check service response
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// technical function, not test, just read test data file
func fetchJsonData(t *testing.T, pathToJson string) (clientData ClientDataType, err error) {
	var (
		jsonFile *os.File
		jsonData []byte
	)

	jsonFile, err = os.Open(pathToJson)
	if err != nil {
		return clientData, errors.New(fmt.Sprint("[error] can't open file:", pathToJson))
	}

	defer func() {
		if err = jsonFile.Close(); err != nil {
			fmt.Println("[error] clear memory file")
		}
	}()

	jsonData, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		return clientData, errors.New("[error] read json data")
	}

	err = json.Unmarshal(jsonData, &clientData)
	if err != nil {
		return clientData, errors.New("[error] unmarshal json data")
	}

	return clientData, err
}

// check what we will have if client will send correct request
func TestValidateClientRequestOk(t *testing.T)  {
	var(
		err error

		clientData ClientDataType
	)

	clientData, err = fetchJsonData(t, "../testdata/client_request_all_ok.json")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	err = validateClientRequest(clientData)
	if err != nil {
		t.Error(err.Error())
	}
}

// check what we will have if client will send request with bad http method
func TestValidateClientRequestFail(t *testing.T)  {
	var(
		err error

		clientData ClientDataType
	)

	clientData, err = fetchJsonData(t, "../testdata/client_request_fail_method.json")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	err = validateClientRequest(clientData)
	if err == nil {
		t.Error("unexpected success")
	}
}