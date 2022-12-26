package main

import (
	"log"

	"zincsearch.com/proxy/override/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error en la carga de variables de ambiente: %s", err)
	}

}

func main() {

}
