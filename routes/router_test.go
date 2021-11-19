package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"

	"github.com/req_counter_service/config"

)

// type serviceInfo struct {
// 	port string
// 	host string
// }

func TestReturnReqCounterStatus(t *testing.T) {

	configurations, errConf := config.Initializer()
	if errConf != nil {
		panic(errConf)
	}

	info := serviceInfo {
		port : configurations.ServicePort,
		host : configurations.HostName,
	}
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(info.returnReqCounter)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// expected := "You are talking to instance host1:8083.\nThis is request 1 to this instance and request 1 to the cluster.\n"
	expected := fmt.Sprintf("You are talking to instance host1:%s.\nThis is request 1 to this instance and request 1 to the cluster.\n", configurations.ServicePort)
	fmt.Println(rr.Body.String())
	fmt.Println(expected)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
