package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

/*Configuration ... */
type Configuration struct {
	ServiceHost   	string
	ServicePort    	string
	ContainerDBhost string
	ContainerDBport string
	HostName 		string
}


/*Initializer ... 
@Returns: a pointer to a configuration struct variable OR an error opening file
@Functionality: reads the .env file
*/
func Initializer() (*Configuration, error) {
	// open .env file in the local directory
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file!")
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return &Configuration{
		ServiceHost:     os.Getenv("SERVICE_HOST"),
		ServicePort:     os.Getenv("SERVICE_PORT"),
		ContainerDBhost: os.Getenv("CONTAINER_DB_HOST"),
		ContainerDBport: os.Getenv("CONTAINER_DB_PORT"),
		HostName: 		 hostname,
	}, nil
}
