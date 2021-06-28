FROM golang:1.16.3 as modules
ADD app/go.mod app/go.sum /m/
RUN cd /m && go mod download

FROM golang:1.16.3 as builder
COPY --from=modules /go/pkg /go/pkg
RUN mkdir -p /src
ADD . /src
WORKDIR /src/app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
go build -o /shortener ./cmd/shortener

FROM node:alpine AS node_builder
COPY --from=builder /shortener /shortener
COPY --from=builder /src/app/cmd/shortener/config_docker.json /config.json
COPY --from=builder /src/ui/thapp/my-react-tutorial-app ./
RUN npm install
RUN npm install -g serve
RUN npm run build

FROM alpine:latest
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /shortener /shortener
COPY --from=builder /src/app/cmd/shortener/config_docker.json /config.json
COPY --from=node_builder /build ./web
# COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
RUN chmod +x ./shortener
EXPOSE 8080
CMD ["/shortener"]