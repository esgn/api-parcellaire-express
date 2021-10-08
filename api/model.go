package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

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
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	fc := geojson.NewFeatureCollection()
	for rows.Next() {
		var a parcelle
		if err := rows.Scan(&a.idu, &a.numero, &a.feuille, &a.section, &a.nom_com, &a.code_com, &a.com_abs, &a.code_arr, &a.geometry); err != nil {
			log.Println(err.Error())
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

	rerr := rows.Close()
	if rerr != nil {
		log.Println(rerr.Error())
		return nil, rerr
	}

	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return fc, nil
}

func bboxToPolygon(bbox string) string {
	c := strings.Split(bbox, ",")
	return fmt.Sprintf("POLYGON ((%[1]s %[2]s, %[1]s %[4]s, %[3]s %[4]s, %[3]s %[2]s, %[1]s %[2]s))", c[0], c[1], c[2], c[3])
}

func getParcelle(db *sql.DB, key, value string) (*geojson.FeatureCollection, error) {
	return getGeoJSON(db, fmt.Sprintf("SELECT idu, numero, feuille, section, nom_com, code_com, com_abs, code_arr, ST_AsGeoJSON(geom) FROM %s.parcelle WHERE %s=$1", os.Getenv("POSTGRES_SCHEMA"), key), value)
}

func getParcelleIntersects(db *sql.DB, pos string) (*geojson.FeatureCollection, error) {
	return getGeoJSON(db, fmt.Sprintf("SELECT idu, numero, feuille, section, nom_com, code_com, com_abs, code_arr, ST_AsGeoJSON(geom) FROM %s.parcelle WHERE ST_Intersects(geom,ST_SetSRID(ST_MakePoint(%s),4326))", os.Getenv("POSTGRES_SCHEMA"), pos))
}

func getParcelleBbox(db *sql.DB, bbox string) (*geojson.FeatureCollection, error) {
	_sql := "SELECT idu, numero, feuille, section, nom_com, code_com, com_abs, code_arr, ST_AsGeoJSON(geom) FROM %s.parcelle WHERE geom && ST_GeomFromText('%s', 4326)"

	if os.Getenv("LIMIT_FEATURE") == "1" {
		return getGeoJSON(db, fmt.Sprintf(_sql+" LIMIT %s", os.Getenv("POSTGRES_SCHEMA"), bboxToPolygon(bbox), os.Getenv("MAX_FEATURE")))
	}

	return getGeoJSON(db, fmt.Sprintf(_sql, os.Getenv("POSTGRES_SCHEMA"), bboxToPolygon(bbox)))
}
