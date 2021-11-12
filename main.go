package main

import (
	"fmt"
	"net/http"

	"github.com/req_counter_service/routes"
)

func main() {

	fmt.Println("Running :: http://localhost:8083")

	r := routes.Router()
	fmt.Println(http.ListenAndServe(":8083", r))
}
