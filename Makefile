default: test

vendor: Gopkg.toml Gopkg.lock
	dep ensure

vendor_update:
	dep ensure -update

install:
	go install ./...

test:
	go test ./...

cov:
	go test github.com/billtrust/go-pingdom/pingdom -coverprofile=coverage.out
	go tool cover -func=coverage.out
	rm coverage.out
