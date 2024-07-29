FROM golang:1.22-alpine

WORKDIR /app

# Installer git et bash
RUN apk add --no-cache git bash

# Copier go.mod et go.sum
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier le reste des fichiers source
COPY . .
# Afficher le contenu du répertoire pour le débogage
RUN ls -la && ls -la handlers

# Construire l'application
RUN go build -o go-cert-api .

CMD ["./go-cert-api"]