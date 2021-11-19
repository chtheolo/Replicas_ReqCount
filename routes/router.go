package routes

import (
	"net/http"
	"fmt"
	"sync"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/req_counter_service/redis"
)

/*Set total transanction attempts for incrementing total count*/
var totalRedisTxAttempts = 2

/*Global static variable for counting the number of request of the instance.*/
var reqCount int = 0

/* HTTP plain/text return messages*/
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
	muTotaInc sync.Mutex
)

/*Function that increments the local counter protecting it with mutex*/
func incrCount() {
	muIncr.Lock()
	reqCount++
	defer muIncr.Unlock()
}


func incrTotalCount() (string, error) {
	muTotaInc.Lock()
	total, errIncrement := redis.IncrGetTotalCount()
	if errIncrement != nil {
		return "", errIncrement
	}
	defer muTotaInc.Unlock()
	return total, nil
}

func (info serviceInfo) returnReqCounter(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/text")

	var total string
	for attempts := 0; attempts < totalRedisTxAttempts; attempts++ {
		totalCount, errIncr := incrTotalCount()
		if errIncr != nil {
			if attempts == totalRedisTxAttempts-1 {
				fmt.Println(fmt.Errorf("Increment Transaction Error:%s", errIncr.Error()))
				w.Write([]byte("Increment TransactionError"))
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
			fmt.Println(fmt.Errorf("Increment Transaction Error:%s", errIncr.Error()))
		}
		
		if totalCount != "" {
			total = fmt.Sprintf("%s", totalCount)
			break
		}
	}

	/*Increment local counter with mutexes*/
	incrCount()

	/*Convert the reqCount to string and add in our response*/
	sReqCount := strconv.Itoa(reqCount)

	/*Build the http response string*/
	sResponse := fmt.Sprintf("%s%s:%s.\n%s%s%s%s%s", message1, info.host, info.port, message2, sReqCount, message3, total, message4)

	_, errResponse := w.Write([]byte(sResponse))
	if errResponse != nil {
		fmt.Println("Internal Service Error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

/*Router is the router handler of our service.*/
func Router(servicePort string, serviceHostname string) *mux.Router {

	info := serviceInfo {
		port : servicePort,
		host : serviceHostname,
	}
	
	Router := mux.NewRouter().StrictSlash(true)

	Router.HandleFunc("/", info.returnReqCounter).Methods("GET")
	return Router
}
