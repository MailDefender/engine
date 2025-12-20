# syntax=docker/dockerfile:1

FROM golang:1.25 AS build

ENV GOROOT=

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY --exclude=go.* . .

# Build
WORKDIR /app/cmd/engine
RUN CGO_ENABLED=0 GOOS=linux go build -o /engine

FROM alpine AS atlas

WORKDIR /inst

RUN apk add curl
RUN curl -sSf https://atlasgo.sh -o atlas.sh
RUN chmod +x atlas.sh
RUN ./atlas.sh --yes --no-install --output atlas
RUN chmod +x atlas/atlas-linux-amd64-latest

FROM alpine

WORKDIR /engine

COPY --from=build /engine engine
RUN chmod +x engine

COPY --from=atlas /inst/atlas/atlas-linux-amd64-latest /utils/atlas
COPY atlas.hcl atlas.hcl

COPY start_engine.sh start.sh
RUN chmod +x start.sh

COPY migrations migrations/

EXPOSE 8080

# Run
CMD ["./start.sh"]