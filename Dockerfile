FROM golang:1.24-alpine AS base


RUN apk update 
RUN apk add --no-cache git ca-certificates

WORKDIR /upsm
ADD . /upsm/


COPY go.mod go.sum ./

# BUILD
RUN go mod download
RUN go env -w CGO_ENABLED=0 && go build -o /go/bin/upsmapi .

# On copie les certificats dans un layer séparé
FROM scratch AS app
EXPOSE 9695
WORKDIR /upsm

# Copier le binaire
COPY --from=base /go/bin/upsmapi /upsm/upsmapi

# Copier les certificats SSL
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/upsm/upsmapi"]