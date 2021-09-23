FROM golang:1.17-alpine

# Port d'écoute par défaut de l'API
ARG api_port=8010

# Création d'un utilisateur et groupe dédié
ENV USER=api
ENV UID=10001
ENV GROUP=api
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

# Copie du code source de l'API
COPY . /go/src/app
WORKDIR /go/src/app

# Build des sources de l'API
RUN go mod download
RUN go build -o api.go

# Exposition du port de l'API
EXPOSE ${api_port}

# Execution avec l'utilisateur dédié
USER api:api

CMD ["./api.go"]
