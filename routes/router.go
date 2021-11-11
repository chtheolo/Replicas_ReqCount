package routes

import (
	"net/http"
	// "os"
	"strconv"

	"github.com/gorilla/mux"
)

// type response struct {
// 	message   string
// 	count_req int
// }

var count_req = 0

func returnReqCounter(w http.ResponseWriter, r *http.Request) {
	count_req++
	s2 := strconv.Itoa(count_req)
	// name, err := os.Hostname()
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = w.Write([]byte(name, " ", counter_req))
	// _, err = w.Write([]byte(count_req))
	w.Header().Set("Content-Type", "application/text")
	_, err := w.Write([]byte(s2))
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}

/*export Router*/
func Router() *mux.Router {
	Router := mux.NewRouter().StrictSlash(true)

	Router.HandleFunc("/", returnReqCounter).Methods("GET")
	return Router
}
