FROM golang:1.24.3-bookworm AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go env -w CGO_ENABLED=1

RUN go build -o simple-rest-server

FROM ubuntu:rolling

WORKDIR /app

COPY --from=build /app/simple-rest-server simple-rest-server

EXPOSE 8080

ENTRYPOINT ["/app/simple-rest-server"]
