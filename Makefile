.PHONY: proto

proto: proto/main.pb.go

proto/%.pb.go: proto/%.proto
	protoc --go_out=./ --go_opt=paths=source_relative \
		--go-grpc_out=./ --go-grpc_opt=paths=source_relative \
		$^

binary: bin/mediabankbridge-linux-amd64 proto/*.pb.go

bin/%: **/*.go proto/*.proto
	env GOOS=linux GOARCH=amd64 go build -o $@ .
