package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"zincsearch.com/proxy/override/godotenv"
	"zincsearch.com/proxy/service"
)

var servicio service.ZincSearch

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error en la carga de variables de ambiente: %s", err)
	}

	//inicia cliente de zincsearch
	servicio.Inicia()
}

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//timeout de 40 segundos
	r.Use(middleware.Timeout(40 * time.Second))

	r.Get("/", getRecordsFromZinc)

	//iniciar escucha
	http.ListenAndServe(":8001", r)
}

func getRecordsFromZinc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola mundo!"))

	rec, err := servicio.GetRecords()
}
