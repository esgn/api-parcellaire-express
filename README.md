# api-parcellaire-express

API REST minimale pour diffusion simple du produit IGN Parcellaire Express : https://geoservices.ign.fr/ressources_documentaires/Espace_documentaire/PARCELLAIRE_CADASTRAL/Parcellaire_Express_PCI/DL_Parcellaire_Express_PCI.pdf

Seules les parcelles sont diffusées pour le moment.

## Préalables

1. Changer le mot de passe par défaut de la base (`password`) dans le fichier docker-compose.yml
2. Adapter les options de configuration de PostgreSQL dans le fichier docker-compose.yml si nécessaire

## Déploiement

1. Téléchargement du produit Parcellaire Express

    `python3 download-dataset.py`

2. Lancement des containers

    `docker-compose up -d`

3. Import des données en base (opération longue pouvant être lancée dans un `screen` par exemple)

    `docker exec -ti parcellaire-express-postgis /bin/bash /tmp/scripts/import-data.sh`

## Utilisation

### Routes

* **GET** `/parcelle/{idu}` : Récupération d'une parcelle à partir de son identifiant
* **GET** `/parcelle?lon={lon}&lat={lat}` *ou* `/parcelle?pos={pos}` : Recherche des parcelles intersectant une position donnée en coordonnées géographiques
* **GET** `/parcelle?lon_min={lon}&lat_min={lat}&lon_max={lon}&lat_max={lat}` *ou* `/parcelle?bbox={bbox}` : Recherche des parcelles intersectant une bounding box donnée en coordonnées géographiques

### Formats

#### Paramètres

* `{idu}` : Identifiant de parcelle
  ```goregexp
  [0-9A-Z]{14}
  ```
* `{lon}` : Longitude (décimal de -180 à 180)
  ```goregexp
  -?0*(?:180(?:\\.0+)?|1[0-7][0-9](?:\.[0-9]+)?|[0-9]{1,2}(?:\.[0-9]+)?)
  ```
* `{lat}` : Latitude (décimal de -90 à 90)
  ```goregexp
  -?0*(?:90(?:\\.0+)?|[0-8]?[0-9](?:\.[0-9]+)?)
  ```
* `{pos}` : Position géographique composé de 2 coordonnées (`lon,lat`)
  ```goregexp
  -?0*(?:180(?:\\.0+)?|1[0-7][0-9](?:\.[0-9]+)?|[0-9]{1,2}(?:\.[0-9]+)?),-?0*(?:90(?:\\.0+)?|[0-8]?[0-9](?:\.[0-9]+)?)
  ```
* `{bbox}` : Bounding box composée de 4 coordonnées (`lon_min,lat_min,lon_max,lat_max`)
  ```goregexp
  -?0*(?:180(?:\\.0+)?|1[0-7][0-9](?:\.[0-9]+)?|[0-9]{1,2}(?:\.[0-9]+)?),-?0*(?:90(?:\\.0+)?|[0-8]?[0-9](?:\.[0-9]+)?),-?0*(?:180(?:\\.0+)?|1[0-7][0-9](?:\.[0-9]+)?|[0-9]{1,2}(?:\.[0-9]+)?),-?0*(?:90(?:\\.0+)?|[0-8]?[0-9](?:\.[0-9]+)?)
  ```

#### Résultats

Les résultats sont fournis au format [GeoJSON](https://fr.wikipedia.org/wiki/GeoJSON).

## Arrêt du service

* Sans suppression des données en base : `docker-compose down`
* Avec suppression des données en base : `docker-compose down -v`
