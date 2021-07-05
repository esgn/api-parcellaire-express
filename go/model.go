package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	geojson "github.com/paulmach/go.geojson"
)

type parcelle struct {
	idu      string
	numero   string
	feuille  int
	section  string
	nom_com  string
	code_com string
	com_abs  string
	code_arr string
	geometry *geojson.Geometry
}

func getGeoJSON(db *sql.DB, query string, args ...interface{}) (*geojson.FeatureCollection, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Println("Cannot close request")
		}
	}(rows)

	fc := geojson.NewFeatureCollection()
	for rows.Next() {
		var a parcelle
		if err := rows.Scan(&a.idu, &a.numero, &a.feuille, &a.section, &a.nom_com, &a.code_com, &a.com_abs, &a.code_arr, &a.geometry); err != nil {
			return nil, err
		}
		f := geojson.NewFeature(a.geometry)
		f.SetProperty("idu", a.idu)
		f.SetProperty("numero", a.numero)
		f.SetProperty("feuille", a.feuille)
		f.SetProperty("section", a.section)
		f.SetProperty("nom_com", a.nom_com)
		f.SetProperty("code_com", a.code_com)
		f.SetProperty("com_abs", a.com_abs)
		f.SetProperty("code_arr", a.code_arr)
		fc.AddFeature(f)
	}

	return fc, nil
}

func getParcelle(db *sql.DB, key, value string) (*geojson.FeatureCollection, error) {
	return getGeoJSON(db, fmt.Sprintf("SELECT idu, numero, feuille, section, nom_com, code_com, com_abs, code_arr, ST_AsGeoJSON(geom) FROM %s.parcelle WHERE %s=$1", os.Getenv("APP_DB_SCHEMA"), key), value)
}

func getParcelleIntersects(db *sql.DB, pos string) (*geojson.FeatureCollection, error) {
	return getGeoJSON(db, fmt.Sprintf("SELECT idu, numero, feuille, section, nom_com, code_com, com_abs, code_arr, ST_AsGeoJSON(geom) FROM %s.parcelle WHERE ST_Intersects(geom,ST_SetSRID(ST_MakePoint(%s),4326))", os.Getenv("APP_DB_SCHEMA"), pos))
}

func getParcelleBbox(db *sql.DB, bbox string) (*geojson.FeatureCollection, error) {
	return getGeoJSON(db, fmt.Sprintf("SELECT idu, numero, feuille, section, nom_com, code_com, com_abs, code_arr, ST_AsGeoJSON(geom) FROM %s.parcelle WHERE ST_Intersects(geom,ST_MakeEnvelope(%s,4326))", os.Getenv("APP_DB_SCHEMA"), bbox))
}
