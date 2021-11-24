.PHONY: proto

proto: src/proto/main.pb.go

src/proto/%.pb.go: proto/%.proto
	protoc --go_out=./src --go_opt=paths=source_relative \
		--go-grpc_out=./src --go-grpc_opt=paths=source_relative \
		$^

binary: bin/mediabankbridge-linux-amd64 src/proto/*

bin/%: src/** src/proto/*
	cd src && env GOOS=linux GOARCH=amd64 go build -o ../$@ .
