package routes

import (
	"net/http"
	"fmt"
	"sync"
	"context"

	"github.com/gorilla/mux"
	"github.com/req_counter_service/redis"
	"github.com/req_counter_service/logs"
)


// Global static variable for counting the number of request of the instance.
var reqCount int = 0

// HTTP plain/text return messages.
var (
	message1 = "You are talking to instance "
	message2 = "This is request "
	message3 = " to this instance and request "
	message4 = " to the cluster.\n"
)

type serviceInfo struct {
	port string
	host string
}

var (
	muIncr sync.Mutex
	muRead sync.Mutex
)

// Function that increments the local counter protecting it with mutex and return him back
// as a string.
func incrCount() {
	muIncr.Lock()
	defer muIncr.Unlock()
	reqCount++
}

// Function that reads the local count
func readCount() int{
	muRead.Lock()
	defer muRead.Unlock()
	return reqCount
}

func (info serviceInfo) returnReqCounter(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	w.Header().Set("Content-Type", "application/text")

	// Increment cluster counter
	total, errIncrement := redis.IncrGetTotalCount(ctx)
	if errIncrement != nil {
		fmt.Println(fmt.Errorf("Increment Transaction Error:%s", errIncrement.Error()))
		logs.WriteLog(errIncrement.Error())

		w.Write([]byte("Increment Transaction Error"))
		w.WriteHeader(http.StatusNotImplemented)
		return 
	}

	// Increment local counter
	incrCount()

	// Read local counter
	localCount := readCount()

	// Build the http response string
	strSuccessResponse := fmt.Sprintf("%s%s:%s.\n%s%d%s%s%s", message1, info.host, info.port, message2, localCount, message3, total, message4)

	_, errResponse := w.Write([]byte(strSuccessResponse))
	if errResponse != nil {
		fmt.Println("Internal Service Error")
		logs.WriteLog(errResponse.Error())

		w.Write([]byte("Internal Service Error"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Router is the router handler of our service.
func Router(servicePort string, serviceHostname string) *mux.Router {

	info := serviceInfo {
		port : servicePort,
		host : serviceHostname,
	}
	
	Router := mux.NewRouter().StrictSlash(true)

	Router.HandleFunc("/", info.returnReqCounter).Methods("GET")
	return Router
}
