#!/usr/bin/env bash

# Dezippage et import des donnees en base
function import() {
    xdir=$(basename "$1" .7z)
    base_dir=$(pwd)
    7z -bso0 -bsp0 -y x "$1" -o"$xdir"
    shp_dir="$(find "$xdir" -name "PARCELLE.SHP" -printf '%h' -quit)"
    cd "$shp_dir" || exit
    shp2pgsql -s "$2":"$3" -D "$4" PARCELLE.SHP "$POSTGRES_SCHEMA".parcelle | psql
    cd "$base_dir" || exit
    rm -rf "$xdir"
    rm -rf "$1"
}
export -f import

# Export des variables d'environnement pour psql
export PGUSER=$POSTGRES_USER
export PGDATABASE=$POSTGRES_DB
export PGPASSWORD=$POSTGRES_PASSWORD
export PGHOST=$POSTGRES_HOST

# Suppression du schéma et des données existantes et recréation
psql -c "DROP SCHEMA IF EXISTS $POSTGRES_SCHEMA CASCADE;"
psql -c "CREATE SCHEMA $POSTGRES_SCHEMA;"
psql -c "CREATE EXTENSION IF NOT EXISTS postgis;"

# Initialisation
cd "/tmp/parcellaire-express/" || exit
append=''
src_epsg=2154
dst_epsg=4326

# Import des shapefiles dans postgis avec gestion des différentes projections
params=()
for f in *.7z; do
    if [[ $f =~ "RGAF09UTM20" ]]; then
        src_epsg=5490
    elif [[ $f =~ "RGM04UTM38S" ]]; then
        src_epsg=4471
    elif [[ $f =~ "RGR92UTM40S" ]]; then
        src_epsg=2975
    elif [[ $f =~ "UTM22RGFG95" ]]; then
        src_epsg=2972
    fi
    params+=("$f $src_epsg $dst_epsg $append")
    append="-a"
done

parallel --verbose --colsep ' ' -j 1 import ::: "${params[0]}"
params=("${params[@]:1}")
parallel --verbose --colsep ' ' -j+0 import ::: "${params[@]}"

# Création des index
psql -c "CREATE INDEX parcelle_geom_idx ON $POSTGRES_SCHEMA.parcelle USING GIST (geom)"
psql -c "CREATE INDEX parcelle_idu_idx ON $POSTGRES_SCHEMA.parcelle (idu)"
psql -c "VACUUM ANALYZE"
