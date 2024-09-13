run:
	go run cmd/main.go
prot-exp:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	export PATH="$PATH:$(go env GOPATH)/bin"

gen-proto:
	protoc --go_out=./ \
    --go-grpc_out=./ \
	submodule/order_service/*.proto