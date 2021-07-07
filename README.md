# api-parcellaire-express

API REST minimale pour diffusion simple du produit IGN [Parcellaire Express](https://geoservices.ign.fr/ressources_documentaires/Espace_documentaire/PARCELLAIRE_CADASTRAL/Parcellaire_Express_PCI/DL_Parcellaire_Express_PCI.pdf).

Seules les parcelles sont diffusées pour le moment.

## Prérequis

Les outils suivants doivent être installés sur la machine :
* Python 3
* Docker
* docker-compose

## Préalables

1. Changer le mot de passe par défaut de la base (`password`) dans le fichier docker-compose.yml
2. Adapter les options de configuration de PostgreSQL dans le fichier docker-compose.yml si nécessaire
3. Installer les dépendances python nécessaires au téléchargement des données : `pip install -r requirements.txt`

## Déploiement

1. Téléchargement du produit IGN Parcellaire Express (PCI)

    `python download-dataset.py`

2. Lancement des containers via docker-compose

    `docker-compose up -d`

3. Import des données en base (opération longue pouvant être lancée dans un `screen` par exemple)

    `docker exec -ti parcellaire-express-postgis /bin/bash /tmp/scripts/import-data.sh`

## Utilisation

### Routes

* **GET** `/parcelle/{idu}` : Récupération d'une parcelle à partir de son identifiant
  * Exemple : http://localhost:8010/parcelle/69389000CR0048
* **GET** `/parcelle?pos={pos}` *ou* `/parcelle?lon={lon}&lat={lat}` : Recherche des parcelles intersectant une position donnée en coordonnées géographiques (WGS84)
  * Exemple : http://localhost:8010/parcelle?pos=5.2709,44.6247
* **GET** `/parcelle?bbox={bbox}` *ou* `/parcelle?lon_min={lon}&lat_min={lat}&lon_max={lon}&lat_max={lat}` : Recherche des parcelles intersectant une bounding box donnée en coordonnées géographiques (WGS84)
  * Exemple : http://localhost:8010/parcelle?bbox=5.2135,44.5719,5.2709,44.6247

### Formats

#### Paramètres

* `{idu}` : Identifiant unique de parcelle (ex: `69389000CR0048`)
* `{lon}` : Longitude (décimal entre -180 et 180)
* `{lat}` : Latitude (décimal entre -90 et 90)
* `{pos}` : Position géographique composé de 2 coordonnées (`lon,lat`)
* `{bbox}` : Bounding box composée de 4 coordonnées (`lon_min,lat_min,lon_max,lat_max`)

#### Résultats

Les résultats sont fournis au format [GeoJSON](https://fr.wikipedia.org/wiki/GeoJSON).

## Arrêt du service

* Sans suppression des données importées en base : `docker-compose down`
* Avec suppression des données importées en base : `docker-compose down -v`
