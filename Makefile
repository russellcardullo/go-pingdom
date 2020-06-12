default: test

vendor:
	go mod vendor

install:
	go install ./...

lint:
	golint github.com/sightmachine/go-pingdom/pingdom

test:
	go test github.com/sightmachine/go-pingdom/pingdom

acceptance:
	PINGDOM_ACCEPTANCE=1 go test github.com/sightmachine/go-pingdom/acceptance

cov:
	go test github.com/sightmachine/go-pingdom/pingdom -coverprofile=coverage.out
	go tool cover -func=coverage.out
	rm coverage.out

.PHONY: default vendor vendor_update install test acceptance cov
