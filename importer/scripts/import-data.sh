#!/bin/bash

# Export des variables d'environnement pour la commande psql
export PGUSER=$POSTGRES_USER
export PGDATABASE=$POSTGRES_DB
export PGPASSWORD=$POSTGRES_PASSWORD
export PGHOST=$POSTGRES_HOST

# Suppression du schéma et des données existantes et recréation
psql -c "DROP SCHEMA IF EXISTS $POSTGRES_SCHEMA CASCADE;"
psql -c "CREATE SCHEMA $POSTGRES_SCHEMA;"
psql -c "CREATE EXTENSION IF NOT EXISTS postgis;"

# Initialisation
cd /tmp/parcellaire-express/
append=''
src_epsg=2154
dst_epsg=4326

# Import des shapefiles dans postgis avec gestion des différentes projections
for f in *.7z; do
    if [[ $f =~ "RGAF09UTM20" ]]
    then
        src_epsg=5490
    elif [[ $f =~ "RGM04UTM38S" ]]
    then
        src_epsg=4471
    elif [[ $f =~ "RGR92UTM40S" ]]
    then
        src_epsg=2975
    elif [[ $f =~ "UTM22RGFG95" ]]
    then
        src_epsg=2972
    fi
    7z x $f
    xdir=`basename $f .7z`
    cd "$(find $xdir -name "PARCELLE.SHP" -printf '%h' -quit)"
    # Utilisation du flag -D (Use postgresql dump format) pour de meilleures performances d'import
    shp2pgsql -s $src_epsg:$dst_epsg -D $append PARCELLE.SHP $POSTGRES_SCHEMA.parcelle | psql
    cd /tmp/parcellaire-express/
    # Suppression du dossier d'extraction
    rm -rf $xdir
    # Suppression de l'archive
    rm -rf $f
    append='-a'
done

# Création des index
psql -c "CREATE INDEX parcelle_geom_idx ON $POSTGRES_SCHEMA.parcelle USING GIST (geom)"
psql -c "CREATE INDEX parcelle_idu_idx ON $POSTGRES_SCHEMA.parcelle (idu)"
psql -c "VACUUM ANALYZE"
