package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"

	"github.com/req_counter_service/config"
)

func TestReturnReqCounterStatus(t *testing.T) {

	configurations, errConf := config.Initializer()
	if errConf != nil {
		t.Fatal(errConf)
	}

	info := serviceInfo {
		port : configurations.ServicePort,
		host : configurations.HostName,
	}
	// Create a request to pass to our handler in the 'http://localhost:PORT/'.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create ResponseRecorder.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(info.returnReqCounter)

	handler.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body.
	expected := fmt.Sprintf("You are talking to instance host1:%s.\nThis is request 1 to this instance and request 1 to the cluster.\n", configurations.ServicePort)
	fmt.Println(rr.Body.String())
	fmt.Println(expected)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
