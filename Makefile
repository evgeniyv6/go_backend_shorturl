build:
	cd ./app/cmd/shortener && go build -o shortener shortener.go

fmt:
	gofmt -w -s .

test:
	go test -race -cover -coverprofile=coverage.out ./...

check:
	go vet ./...
	golangci-lint run ./...

run:
	cd ./app/cmd/shortener && go run -race shortener.go