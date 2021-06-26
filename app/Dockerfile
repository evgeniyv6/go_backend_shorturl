FROM golang:1.16.3 as modules
ADD go.mod go.sum /m/
RUN cd /m && go mod download
FROM golang:1.16.3 as builder
COPY --from=modules /go/pkg /go/pkg
RUN mkdir -p /src
ADD . /src
WORKDIR /src
RUN useradd -u 10001 gouser
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
go build -o /shortener ./cmd/shortener
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
USER gouser
COPY --from=builder /shortener /shortener
COPY --from=builder /src/cmd/shortener/config_docker.json /config.json
# COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
EXPOSE 6089
CMD ["/shortener"]