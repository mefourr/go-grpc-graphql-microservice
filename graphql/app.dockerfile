FROM golang:1.22.5-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/mefourr/go-graphql-microservice
COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./order/cmd/order

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]

#FROM golang:1.13-alpine3.11 AS build
#RUN apk --no-cache add gcc g++ make ca-certificates
#WORKDIR /go/src/github.com/akhilsharma90/go-graphql-microservice
#COPY go.mod go.sum ./
#COPY vendor vendor
#COPY account account
#COPY catalog catalog
#COPY order order
#COPY graphql graphql
#RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./graphql
#
#FROM alpine:3.11
#WORKDIR /usr/bin
#COPY --from=build /go/bin .
#EXPOSE 8080
#CMD ["app"]