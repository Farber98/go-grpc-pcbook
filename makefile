# Generates output file in /pb. Loads proto file from proto/processor_message.go.
gen:
	protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative proto/*.proto

# Removes files under /pb
clean:
	rm pb/*.go

# runs server.
server:
	go run cmd/server/main.go -port 3033

# runs client.
client:
	go run cmd/client/main.go -addr 0.0.0.0:3033

# run all the tests
test: 
	go test -cover -race ./...