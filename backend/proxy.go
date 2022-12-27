package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Printf("Error en la carga de variables de ambiente: %s", err)
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
	r.Get("/search", getRecordsFromZinc)

	//iniciar escucha
	http.ListenAndServe(":8001", r)
}

// consulta API de ZincSearch
func getRecordsFromZinc(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("query")

	rec, err := servicio.GetRecords(query)

	if err.Code != 0 {
		errorJson, _ := json.Marshal(err)
		http.Error(w, string(errorJson), http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(rec))
	}
}
