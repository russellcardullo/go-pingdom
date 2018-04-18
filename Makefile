all: test install

install:
	go install ./...
	go install github.com/stretchr/testify

test:
	go test ./...

cov:
	go test github.com/russellcardullo/go-pingdom/pingdom -coverprofile=coverage.out
	go tool cover -func=coverage.out
	rm coverage.out
