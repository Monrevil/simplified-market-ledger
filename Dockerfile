# syntax=docker/dockerfile:1

FROM golang:1.16-buster AS build

WORKDIR /app

COPY / ./

RUN go mod download -x
RUN go build -o /go/bin/ledger

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=build /go/bin/ledger /

EXPOSE 50051

USER nonroot:nonroot

CMD ["/ledger"]