package main

import (
	"fmt"
	"net/http"

	"github.com/req_counter_service/logs"
	"github.com/req_counter_service/config"
	"github.com/req_counter_service/routes"
)

func main() {

	// Configurations struct
	configurations, err := config.Initializer()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %s", err.Error()))
		logs.WriteLog(err.Error())
	}

	serviceMessage := fmt.Sprintf("Running :: %s:%s ", configurations.ServiceHost, configurations.ServicePort)
	fmt.Println(serviceMessage)


	r := routes.Router(configurations.ServicePort, configurations.HostName)
	served := fmt.Sprintf(":%s", configurations.ServicePort)
	fmt.Println(http.ListenAndServe(served, r))
}
