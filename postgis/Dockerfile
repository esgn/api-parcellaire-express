ARG POSTGRES_VERSION=13

FROM postgres:${POSTGRES_VERSION}

ARG POSTGIS_VERSION=3

# Installation de postgis
RUN set -ex \
    && EXEC_DEPS=" \
    postgresql-${PG_MAJOR}-postgis-${POSTGIS_VERSION} \
    postgresql-${PG_MAJOR}-postgis-${POSTGIS_VERSION}-scripts \
    postgresql-${PG_MAJOR}-pgrouting \
    postgis \
    " \
    && apt update && apt install -y --no-install-recommends $EXEC_DEPS \
    && rm -rf /var/lib/apt/lists/*

# Utilisateur dédié
USER postgres:postgres
