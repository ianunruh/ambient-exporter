FROM golang:1.22-alpine as build

RUN addgroup -S app && adduser -S app -G app

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags '-extldflags "-static"' -tags timetzdata -o ambient-exporter

FROM scratch

COPY --from=build /go/src/app/ambient-exporter /ambient-exporter
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd

USER app

ENTRYPOINT ["/ambient-exporter"]
