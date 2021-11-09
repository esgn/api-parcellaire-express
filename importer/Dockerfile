FROM debian:stable-slim

# Création d'un utilisateur et groupe dédié
ENV USER=importer
ENV UID=10001
ENV GROUP=importer
ENV GID=10001

RUN addgroup --gid ${GID} ${GROUP}
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \   
    --ingroup "${GROUP}" \
    "${USER}"

# Copie des scripts d'import de données
RUN mkdir /scripts/
COPY --chown=importer:importer scripts/ /scripts/
# When copying to a volume, data is not updated after the first copy.
# you have to delete the volume to import the new data, 
# or copy/replace the data manually from /scripts/
COPY --chown=importer:importer scripts/ /tmp/


# Installation des dépendances
RUN set -ex \
    && EXEC_DEPS=" \
    python3-pip \
    python3 \
    p7zip* \
    postgresql-client \
    postgis \
    " \
    && apt update && apt install -y --no-install-recommends $EXEC_DEPS \
    && pip3 install --no-cache-dir  -r /scripts/requirements.txt \
    && rm -rf /var/lib/apt/lists/*    

# Configuration
WORKDIR /tmp
USER importer
