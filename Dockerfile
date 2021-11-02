FROM golang:latest as builder
WORKDIR /app
ADD . /app
RUN go mod download
RUN go install github.com/traefik/yaegi/internal/cmd/extract
RUN go generate ./...
RUN CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o befe .

FROM scratch AS app
WORKDIR /
COPY --from=builder /app/befe /befe
COPY --from=builder /etc/mime.types /etc/mime.types
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["./befe"]