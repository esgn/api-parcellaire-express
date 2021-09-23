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
COPY --chown=importer:importer scripts/ /tmp/

# Installation des dépendances
RUN apt update && \
    apt install -y python3-pip python3 p7zip* && \
    apt install -y --no-install-recommends postgresql-client postgis
RUN pip3 install -r /tmp/requirements.txt

# Configuration
WORKDIR /tmp
USER importer
