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
			ServiceHost   	string
			ServicePort    	string
			ContainerDBhost string
			ContainerDBport string
			HostName 		string
		}
	*/
	configurations, err := config.Initializer()
	if err != nil {
		fmt.Println("Error in config file")
		panic(err)
	}

	serviceMessage := fmt.Sprintf("Running :: %s:%s ", configurations.ServiceHost, configurations.ServicePort)
	fmt.Println(serviceMessage)


	r := routes.Router(configurations.ServicePort, configurations.HostName)
	served := fmt.Sprintf(":%s", configurations.ServicePort)
	fmt.Println(http.ListenAndServe(served, r))
}
