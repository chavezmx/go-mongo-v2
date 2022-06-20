# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine3.16 AS build
ENV GO111MODULE=on
#ENV GOPATH=/app
#ENV GO111MODULE=off
#ENV CGO_ENABLED 0
#ENV GOOS 'linux'

WORKDIR /app/go-mongo-v2/src

COPY go.mod ./
COPY go.sum ./

COPY config .
COPY controllers .
COPY models .
COPY responses .
COPY routes .
COPY main.go .
#RUN go run main.go

#RUN go mod init go-mongo-v2
#RUN go get "github.com/go-playground/validator/v10"
#RUN go get "github.com/joho/godotenv"
#RUN go get "github.com/labstack/echo/v4"
#RUN go get "github.com/labstack/gommon"
#RUN go get "github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
#RUN go get "go.mongodb.org/mongo-driver"

RUN go build -o /go-mongo-v2 /app/go-mongo-v2/src

##
## Deploy
##
FROM alpine:3.16

WORKDIR /

COPY --from=build /go-mongo-v2 /go-mongo-v2
EXPOSE 6000

USER nonroot:nonroot

ENTRYPOINT ["/go-mongo-v2"]



#FROM alpine:3.16
#
#WORKDIR /app
#ADD go-mongo-v2 .
#CMD chmod +x go-mongo-v2
#
#ENTRYPOINT ["./go-mongo-v2"]
