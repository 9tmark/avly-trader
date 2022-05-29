.PHONY: clean test build
define ReadTxt
$(shell cat $(1))
endef
DCKR_IMG:=9tmark/avly-trader
VERSION:=$(call ReadTxt,VERSION)

clean:
	rm -rf ./build/*; \
	rm ./*.out; \
	go clean

test:
	go test -v ./...

testcov:
	go test ./... -coverpkg=./... -coverprofile=./coverage.out && \
	go tool cover -func=./coverage.out; \
	go tool cover -html=./coverage.out

build:
	mkdir -p ./build && \
	go build -o=./build ./cmd/avly

docker_build:
	docker build -t ${DCKR_IMG}:${VERSION} .

docker_build_latest:
	docker build -t ${DCKR_IMG}:${VERSION} -t ${DCKR_IMG}:latest .

docker_push:
	docker push ${DCKR_IMG}:${VERSION}

docker_push_latest:
	docker push ${DCKR_IMG}:${VERSION}; \
	docker push ${DCKR_IMG}:latest
