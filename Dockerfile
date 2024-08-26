FROM golang:1.22  AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -v -ldflags="-s -w" -o /cyberpower_exporter ./...


FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /cyberpower_exporter /cyberpower_exporter

EXPOSE 9300

USER nonroot:nonroot

ENTRYPOINT ["/cyberpower_exporter"]
