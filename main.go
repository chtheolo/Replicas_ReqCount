package main

import (
	"fmt"
	"net/http"

	"github.com/req_counter_service/config"

	"github.com/req_counter_service/routes"
)

func main() {
	configurations, err := config.ConfigService()
	if err != nil {
		fmt.Println("Error in config file")
		panic(err)
	}

	// fmt.Println("Running :: http://localhost:8083")
	fmt.Sprintf("Running :: %s:%s ", configurations.Service_host, configurations.Service_port)

	r := routes.Router()
	served := fmt.Sprintf(":%s", configurations.Service_port)
	fmt.Println(http.ListenAndServe(served, r))
}
