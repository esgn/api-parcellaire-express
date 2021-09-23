package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	// identifiant unique de parcelle
	iduRegex = "[0-9A-Z]{14}"
	// nombre décimal
	decimalRegex = "(?:\\.[0-9]+)?"
	// flottant entre -180 et 180
	lonRegex = "-?0*(?:180(?:\\.0+)?|1[0-7][0-9]" + decimalRegex + "|[0-9]{1,2}" + decimalRegex + ")"
	// flottant entre -90 et 90
	latRegex = "-?0*(?:90(?:\\.0+)?|[0-8]?[0-9]" + decimalRegex + ")"
	// position : lon,lat
	posRegex = lonRegex + "," + latRegex
	// bbox : lon_min,lat_min,lon_max,lat_max
	bboxRegex = posRegex + "," + posRegex
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname, hostname string) {

	log.SetOutput(os.Stderr)

	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", user, password, dbname, hostname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println("JSON marshalling error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
		log.Println("Could not send response")
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithAppropriate(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, data)
}

func (a App) error(code int, message string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, code, message)
	}
}

func (a App) getById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cleabs := vars["idu"]
	parcelle, err := getParcelle(a.DB, "idu", cleabs)
	respondWithAppropriate(w, parcelle, err)
}

func (a App) findByPosition(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	pos := v.Get("pos")
	parcelle, err := getParcelleIntersects(a.DB, pos)
	respondWithAppropriate(w, parcelle, err)
}

func (a App) findByPositionSplit(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	lon, lat := v.Get("lon"), v.Get("lat")
	parcelle, err := getParcelleIntersects(a.DB, fmt.Sprintf("%s,%s", lon, lat))
	respondWithAppropriate(w, parcelle, err)
}

func (a App) findByBbox(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	bbox := v.Get("bbox")
	parcelle, err := getParcelleBbox(a.DB, bbox)
	respondWithAppropriate(w, parcelle, err)
}

func (a App) findByBboxSplit(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	lon_min, lat_min, lon_max, lat_max :=
		v.Get("lon_min"), v.Get("lat_min"), v.Get("lon_max"), v.Get("lat_max")
	parcelle, err := getParcelleBbox(a.DB, fmt.Sprintf("%s,%s,%s,%s", lon_min, lat_min, lon_max, lat_max))
	respondWithAppropriate(w, parcelle, err)
}

func (a App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// TODO : Add database status information
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/parcelle/{idu:"+iduRegex+"}", a.getById).Methods("GET")

	a.Router.HandleFunc("/parcelle", a.findByPosition).Queries(
		"pos", "{pos:"+posRegex+"}").Methods("GET")

	a.Router.HandleFunc("/parcelle", a.findByPositionSplit).Queries(
		"lon", "{lon:"+lonRegex+"}",
		"lat", "{lat:"+latRegex+"}").Methods("GET")

	a.Router.HandleFunc("/parcelle", a.findByBbox).Queries(
		"bbox", "{bbox:"+bboxRegex+"}").Methods("GET")

	a.Router.HandleFunc("/parcelle", a.findByBboxSplit).Queries(
		"lon_min", "{lon_min:"+lonRegex+"}",
		"lat_min", "{lat_min:"+latRegex+"}",
		"lon_max", "{lon_max:"+lonRegex+"}",
		"lat_max", "{lat_max:"+latRegex+"}").Methods("GET")

	a.Router.HandleFunc("/parcelle", a.error(http.StatusBadRequest, "Requête invalide"))

	a.Router.HandleFunc("/status", a.healthCheckHandler).Methods("GET")

	a.Router.PathPrefix("/").HandlerFunc(a.error(http.StatusNotFound, "URL inconnue"))
}
