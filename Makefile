SOURCES := $(shell find * -type f -name "*.go")
COVERAGE_FILE=coverage.out

converter: $(SOURCES)
	go build -o $@ examples/heic2jpg/main.go

vendor: go.mod go.sum
	go mod vendor -v

deps:
	go get -u -v ./...
	go mod tidy -v

generate: clean
	go get -u github.com/golang/mock/mockgen
	go generate ./...

test: vendor
	go test -race -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)

clean:
	rm -vrf converter \
	vendor \
	$(COVERAGE_FILE)

codecov:
	curl https://keybase.io/codecovsecurity/pgp_keys.asc | gpg --import
	curl -Os https://uploader.codecov.io/latest/linux/codecov
	curl -Os https://uploader.codecov.io/latest/linux/codecov.SHA256SUM
	curl -Os https://uploader.codecov.io/latest/linux/codecov.SHA256SUM.sig
	gpg --verify codecov.SHA256SUM.sig codecov.SHA256SUM
	shasum -a 256 -c codecov.SHA256SUM
	chmod +x codecov
	./codecov
