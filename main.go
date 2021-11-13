package main

import (
	"fmt"
	"net/http"

	"github.com/req_counter_service/config"
	"github.com/req_counter_service/routes"
)

func main() {
	/* configurations =>
	type Configuration struct {
		Service_host      string
		Service_port      string
		Container_db_host string
		Container_db_port string
	}
	*/
	configurations, err := config.ConfigService()
	if err != nil {
		fmt.Println("Error in config file")
		panic(err)
	}

	service_message := fmt.Sprintf("Running :: %s:%s ", configurations.Service_host, configurations.Service_port)
	fmt.Println(service_message)

	r := routes.Router(configurations.Service_port)
	served := fmt.Sprintf(":%s", configurations.Service_port)
	fmt.Println(http.ListenAndServe(served, r))
}
