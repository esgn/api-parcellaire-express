ARG postgres_version=13

FROM postgres:${postgres_version}

ARG postgis_version=3

RUN apt-get update \
 && apt-get install -y postgresql-${PG_MAJOR}-postgis-${postgis_version} \
 && apt-get install -y postgresql-${PG_MAJOR}-postgis-${postgis_version}-scripts \
 && apt-get install -y postgresql-${PG_MAJOR}-pgrouting \
 && apt-get install -y --no-install-recommends postgis \
 && apt-get install -y p7zip* \
 && rm -rf /var/lib/apt/lists/*
