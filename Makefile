default: test

vendor: Gopkg.toml Gopkg.lock
	dep ensure

vendor_update:
	dep ensure -update

install:
	go install ./...

lint:
	golint github.com/russellcardullo/go-pingdom/pingdom

test:
	go test github.com/russellcardullo/go-pingdom/pingdom

acceptance:
	PINGDOM_ACCEPTANCE=1 go test github.com/russellcardullo/go-pingdom/acceptance

cov:
	go test github.com/russellcardullo/go-pingdom/pingdom -coverprofile=coverage.out
	go tool cover -func=coverage.out
	rm coverage.out

.PHONY: default vendor vendor_update install test acceptance cov
