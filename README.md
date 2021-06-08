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

Les résultats sont fournis au format GeoJSON.

* (GET) `/parcelle/{idu}` : Récupération d'une parcelle à partir de son identifiant
* (GET) `/parcelle?lon=...&lat=...` : Recherche des parcelles intersectant une position donnée en coordonnées géographiques

## Arrêt du service

* Sans suppression des données en base : `docker-compose down`
* Avec suppression des données en base : `docker-compose down -v`

