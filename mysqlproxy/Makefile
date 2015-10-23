all: build

build: kingshard

kingshard:
	go build -o bin/sqlproxy ./cmd/proxy

clean:
	@rm -rf bin

test:
	go test ./go/... -race
