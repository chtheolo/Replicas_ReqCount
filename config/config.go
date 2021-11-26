package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Configuration struct has all the necessary data that needs the service to get start.
type Configuration struct {
	ServiceHost   	string
	ServicePort    	string
	ContainerDBhost string
	ContainerDBport string
	HostName 		string
}


// Initializer is a function that gets a pointer to a configuration struct variable OR an error opening file.
// It reads the .env file and returns back a Configuration struct variable.
func Initializer() (*Configuration, error) {
	// open .env file in the local directory
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	hostname, errHostname := os.Hostname()
	if errHostname != nil {
		return nil, errHostname
	}

	return &Configuration{
		ServiceHost:     os.Getenv("SERVICE_HOST"),
		ServicePort:     os.Getenv("SERVICE_PORT"),
		ContainerDBhost: os.Getenv("CONTAINER_DB_HOST"),
		ContainerDBport: os.Getenv("CONTAINER_DB_PORT"),
		HostName: 		 hostname,
	}, nil
}
