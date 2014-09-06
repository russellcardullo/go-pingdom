all: test install

install:
	go install

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	rm coverage.out
