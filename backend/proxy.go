package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	w.Header().Set("Access-Control-Allow-Origin", "*")

	query := r.URL.Query().Get("query")
	pageTxt := r.URL.Query().Get("page")
	if pageTxt == "" {
		pageTxt = "1"
	}

	page, err := strconv.Atoi(pageTxt)

	if err != nil {
		http.Error(w, string("{El valor del argumento 'page' debe ser un numero entero}"), http.StatusBadRequest)
	}
	if page <= 0 {
		http.Error(w, string("{El valor del argumento 'page' debe ser un numero mayor a 0}"), http.StatusBadRequest)
	}

	rec, errorLocal := servicio.GetRecords(query, page)
	if errorLocal.Code != 0 {
		errorJson, _ := json.Marshal(errorLocal)
		http.Error(w, string(errorJson), http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(rec))
	}
}
