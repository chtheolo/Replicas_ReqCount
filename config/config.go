package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Service_host      string
	Service_port      string
	Container_db_host string
	Container_db_port string
}

/*
@Returns: a pointer to a configuration struct variable OR an error opening file
@Functionality: reads the .env file
*/
func ConfigService() (*Configuration, error) {
	// open .env file in the local directory
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file!")
		return nil, err
	}

	return &Configuration{
		Service_host:      os.Getenv("SERVICE_HOST"),
		Service_port:      os.Getenv("SERVICE_PORT"),
		Container_db_host: os.Getenv("CONTAINER_DB_HOST"),
		Container_db_port: os.Getenv("CONTAINER_DB_PORT"),
	}, nil
}
