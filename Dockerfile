FROM golang:1.17-alpine as build

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -tags timetzdata -o ambient-exporter

FROM scratch

COPY --from=build /go/src/app/ambient-exporter /ambient-exporter
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/ambient-exporter"]
