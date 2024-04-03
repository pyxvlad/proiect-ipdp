FROM docker.io/library/golang:alpine

RUN apk add build-base
RUN go install github.com/a-h/templ/cmd/templ@latest

RUN mkdir /gomod

COPY go.mod /gomod

COPY go.sum /gomod

WORKDIR /gomod
RUN go mod download


RUN mkdir /app
COPY . /app

WORKDIR /app
RUN templ generate
ENV CGO_ENABLED=1
RUN go build .

VOLUME [ "/data" ]
ENV IPDP_DB=/data/ipdp.db
CMD ./proiect-ipdp