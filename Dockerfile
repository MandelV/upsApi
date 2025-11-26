FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /upsm

# D'abord les deps pour profiter du cache
COPY go.mod go.sum ./
RUN go mod download

# Puis le reste du code
COPY . .

RUN CGO_ENABLED=0 go build -o /upsm/upsmapi .

########################
# Image finale
########################
FROM alpine:3.20 AS app


RUN apk add --no-cache ca-certificates nut

WORKDIR /upsm
EXPOSE 9695

# Copier le binaire
COPY --from=builder /upsm/upsmapi /upsm/upsmapi

ENTRYPOINT ["/upsm/upsmapi"]
