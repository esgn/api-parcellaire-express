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
	idu      sql.NullString
	numero   sql.NullString
	feuille  sql.NullInt64
	section  sql.NullString
	nom_com  sql.NullString
	code_com sql.NullString
	com_abs  sql.NullString
	code_arr sql.NullString
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

		// 🥲 why no ternary operator OR 'if expression' in go, snif 🥲
		// idu
		if a.idu.Valid {
			_val, _ := a.idu.Value()
			f.SetProperty("idu", _val)
		} else {
			f.SetProperty("idu", nil)
		}
		// numero
		if a.numero.Valid {
			_val, _ := a.numero.Value()
			f.SetProperty("numero", _val)
		} else {
			f.SetProperty("numero", nil)
		}
		// feuille
		if a.feuille.Valid {
			_val, _ := a.feuille.Value()
			f.SetProperty("feuille", _val)
		} else {
			f.SetProperty("feuille", nil)
		}
		// section
		if a.section.Valid {
			_val, _ := a.section.Value()
			f.SetProperty("section", _val)
		} else {
			f.SetProperty("section", nil)
		}
		// nom_com
		if a.nom_com.Valid {
			_val, _ := a.nom_com.Value()
			f.SetProperty("nom_com", _val)
		} else {
			f.SetProperty("nom_com", nil)
		}
		// code_com
		if a.nom_com.Valid {
			_val, _ := a.code_com.Value()
			f.SetProperty("code_com", _val)
		} else {
			f.SetProperty("code_com", nil)
		}
		// com_abs
		if a.nom_com.Valid {
			_val, _ := a.com_abs.Value()
			f.SetProperty("com_abs", _val)
		} else {
			f.SetProperty("com_abs", nil)
		}
		// code_arr
		if a.code_arr.Valid {
			_val, _ := a.code_arr.Value()
			f.SetProperty("code_arr", _val)
		} else {
			f.SetProperty("code_arr", nil)
		}

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
	return getGeoJSON(db, fmt.Sprintf("SELECT idu, numero, feuille, section, nom_com, code_com, com_abs, code_arr, ST_AsGeoJSON(geom) FROM %s.parcelle WHERE %s=$1", os.Getenv(ENV_POSTGRES_SCHEMA), key), value)
}

func getParcelleIntersects(db *sql.DB, pos string) (*geojson.FeatureCollection, error) {
	return getGeoJSON(db, fmt.Sprintf("SELECT idu, numero, feuille, section, nom_com, code_com, com_abs, code_arr, ST_AsGeoJSON(geom) FROM %s.parcelle WHERE ST_Intersects(geom,ST_SetSRID(ST_MakePoint(%s),4326))", os.Getenv(ENV_POSTGRES_SCHEMA), pos))
}

func getParcelleBbox(db *sql.DB, bbox string) (*geojson.FeatureCollection, error) {
	_sql := "SELECT idu, numero, feuille, section, nom_com, code_com, com_abs, code_arr, ST_AsGeoJSON(geom) FROM %s.parcelle WHERE geom && ST_GeomFromText('%s', 4326)"

	_limiter := ""
	if os.Getenv(ENV_MAX_FEATURE) != "0" {
		_limiter = fmt.Sprintf(" LIMIT %s", os.Getenv(ENV_MAX_FEATURE))
	}

	return getGeoJSON(db, fmt.Sprintf(_sql+_limiter, os.Getenv(ENV_POSTGRES_SCHEMA), bboxToPolygon(bbox)))
}
