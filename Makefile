build:
	cd ./cmd/shortener && go build -o shortener shortener.go

fmt:
	gofmt -w -s .

test:
	go test -race -cover -coverprofile=coverage.out ./...

check:
	go vet ./...
	golangci-lint run ./...

run:
	cd ./cmd/shortener && go run -race shortener.go