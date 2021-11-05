# API Parcellaire Express

Exemple d'API REST minimale pour une diffusion simple du produit IGN [Parcellaire Express](https://geoservices.ign.fr/sites/default/files/2021-07/DC_Parcellaire_Express_%28PCI%29_1-0.pdf).

Seules les parcelles cadastrales sont diffusées pour le moment.

## Outils utilisés

* [Go](https://golang.org/)
* [PostgreSQL](https://www.postgresql.org/) / [PostGIS](https://postgis.net/)
* [Adminer](https://www.adminer.org/)
* [Python](https://www.python.org/)

## Prérequis

Les outils suivants doivent être installés sur la machine hôte :
* [Docker](https://docs.docker.com/get-started/overview/)
* [Docker Compose](https://docs.docker.com/compose/)

## Architecture

Trois containers sont utilisés : 
* `parcellaire-express-importer` : Container en charge du téléchargement des donnnées et de leur import en base
* `parcellaire-express-postgis` : Base de données PostGIS
* `parcellaire-express-api` : API en Go

Il est possible de décommenter le service `adminer` dans `docker-compose.yml` pour ajouter une interface web d'exploration de la base de données.

## Préalables

0. 🚨 Copier le fichier `.env.example` vers le fichier `.env` avant toute opération. Les valeurs par défauts devraient être suffisantes, mais il vous est possible de l'adapter à votre environnement.

1. Le fichier `.env` regroupe l'ensemble des valeurs de configuration. On liste ci-dessous les options les plus utiles :
    * Configuration de l'importer
      * `MAX_PARALLEL_DL` : Nombre de téléchargement d'archives de données simultanés. Fixé à `4` par défaut.
      * `TEST_IMPORTER` : A passer à `1` pour tester l'import de données sur une seule archive départementale.
    * Configuration de la base de données
      * `POSTGRES_PASSWORD` :  Mot de passe de la base de données. **A modifier**.
    * Configuration de l'API
      * `API_PORT` : Port d'écoute de l'API. Fixé à `8010` par défaut.
      * `MAX_FEATURE` : Nombre maximal d'objets retournés par l'API. Fixé à `1000` par défaut. `0` pour désactiver la limite.
    * Configuration du viewer
      * `VIEWER_URL` : (Optionel) Url d'accès à une page de consultation des parcelles. Laisser vide ou non défini pour désactiver.
2. Des options de configuration de PostgreSQL sont définies dans le fichier `docker-compose.yml`. Utiliser [PGTune](https://pgtune.leopard.in.ua/#/) pour les adapter aux caractéristiques de la machine hôte.

## Déploiement

1. Construction des images

    `docker-compose build`

2. Lancement des containers via docker-compose

    `docker-compose up -d`

3. Import des données en base (opérations longues pouvant être lancées dans un `screen` ou en utilisant l'option `-d` de `docker-compose run`)

   * Téléchargement des données du produit :

      `docker-compose run parcellaire-importer python3 /tmp/download-dataset.py`

   * Mise en base des données du produit :

      `docker-compose run parcellaire-importer bash /tmp/import-data.sh`


## Utilisation

### Routes

* **GET** `/parcelle/{idu}` : Récupération d'une parcelle à partir de son identifiant
  * Exemple : http://localhost:8010/parcelle/01053000BE0095
* **GET** `/parcelle?pos={pos}` *ou* `/parcelle?lon={lon}&lat={lat}` : Recherche des parcelles intersectant une position donnée en coordonnées géographiques (WGS84)
  * Exemple : http://localhost:8010/parcelle?pos=5.2709,44.6247
* **GET** `/parcelle?bbox={bbox}` *ou* `/parcelle?lon_min={lon}&lat_min={lat}&lon_max={lon}&lat_max={lat}` : Recherche des parcelles intersectant une bounding box donnée en coordonnées géographiques (WGS84)
  * Exemple : http://localhost:8010/parcelle?bbox=5.2135,44.5719,5.2709,44.6247

### Formats

#### Paramètres

* `{idu}` : Identifiant unique de parcelle (ex: `01053000BE0095`)
* `{lon}` : Longitude (décimal entre -180 et 180)
* `{lat}` : Latitude (décimal entre -90 et 90)
* `{pos}` : Position géographique composé de 2 coordonnées (`lon,lat`)
* `{bbox}` : Bounding box composée de 4 coordonnées (`lon_min,lat_min,lon_max,lat_max`)

#### Résultats

Les résultats sont fournis au format [GeoJSON](https://geojson.org/).

## Arrêt du service

* Sans suppression des données importées en base : `docker-compose down`
* Avec suppression des données importées en base : `docker-compose down -v`

## TODO

- [ ] Paging des résultats
- [ ] Ajout de test unitaires
- [ ] Eventuel rapprochement avec [OGC API Feature](https://www.ogc.org/standards/ogcapi-features)
- [ ] Gestion des autres classes du produit Parcellaire Express
- [ ] Fiabilisation
