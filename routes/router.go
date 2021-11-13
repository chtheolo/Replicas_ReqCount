package routes

import (
	"fmt"
	"net/http"

	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/req_counter_service/redis"
)

/*Global static variable for counting the number of request of the instance.*/
var req_count = 0

/* HTTP plain/text return messages*/
var (
	message_1 = "You are talking to instance "
	message_2 = "This is request "
	message_3 = " to this instance and request "
	message_4 = " to the cluster.\n"
)

var PORT string

func returnReqCounter(w http.ResponseWriter, r *http.Request) {

	total := redis.GetData()
	total_count, err_convert := strconv.Atoi(total)

	if err_convert != nil {
		fmt.Println("convert to int total count")
		panic(err_convert)
	}

	req_count++
	total_count++

	redis.SetData(total_count)

	s_total_count := strconv.Itoa(total_count)
	s_req_count := strconv.Itoa(req_count)

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	/*Build the http response string*/
	s_response := message_1 + hostname + ":" + PORT + ".\n" + message_2 + s_req_count + message_3 + s_total_count + message_4

	w.Header().Set("Content-Type", "application/text")
	_, err1 := w.Write([]byte(s_response))
	if err1 != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

/*export Router*/
func Router(port string) *mux.Router {
	PORT = port
	Router := mux.NewRouter().StrictSlash(true)

	Router.HandleFunc("/", returnReqCounter).Methods("GET")
	return Router
}
