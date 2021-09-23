ARG POSTGRES_VERSION=13

FROM postgres:${POSTGRES_VERSION}

ARG POSTGIS_VERSION=3

# Installation de postgis
RUN apt-get update \
 && apt-get install -y postgresql-${PG_MAJOR}-postgis-${POSTGIS_VERSION} \
 && apt-get install -y postgresql-${PG_MAJOR}-postgis-${POSTGIS_VERSION}-scripts \
 && apt-get install -y postgresql-${PG_MAJOR}-pgrouting \
 && apt-get install -y --no-install-recommends postgis \
 && rm -rf /var/lib/apt/lists/*

# Utilisateur dédié
USER postgres:postgres
