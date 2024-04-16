FROM docker.io/library/golang:alpine

RUN apk add build-base
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN mkdir /gomod

COPY go.mod /gomod

COPY go.sum /gomod

WORKDIR /gomod
RUN go mod download


RUN mkdir /app
COPY . /app

WORKDIR /app
RUN templ generate
WORKDIR /app/database
RUN sqlc generate
WORKDIR /app
ENV CGO_ENABLED=1
RUN go build .

VOLUME [ "/data" ]
ENV IPDP_DB=/data/ipdp.db
CMD ./proiect-ipdp
