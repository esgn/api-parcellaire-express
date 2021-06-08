package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname, hostname string) {

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
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (a *App) getById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cleabs := vars["idu"]

	parcelle, err := getParcelle(a.DB, "idu", cleabs)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Aucun résultat trouvé")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, parcelle)
}

func (a *App) findByPosition(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	lon := v.Get("lon")
	lat := v.Get("lat")

	_lon, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "lon n'est pas un float")
	}

	_lat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "lat n'est pas un float")
	}

	parcelle, err := getParcelleIntersects(a.DB, _lon, _lat)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Aucun résultat trouvé")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, parcelle)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/parcelle/{idu}", a.getById).Methods("GET")
	a.Router.HandleFunc("/parcelle", a.findByPosition).Queries("lon", "{lon}", "lat", "{lat}").Methods("GET")
}
