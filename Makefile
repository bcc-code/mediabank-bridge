.PHONY: proto build

all: build

proto: proto/main.pb.go

proto/%.pb.go: proto/%.proto
	protoc --go_out=./ --go_opt=paths=source_relative \
		--go-grpc_out=./ --go-grpc_opt=paths=source_relative \
		$^

build: proto
	env GOOS=linux GOARCH=amd64 go build -o ./bin/mediabankbridge-linux-amd64 ./cmd/server

clean:
	rm -v ./bin/mediabankbridge-linux-amd64
