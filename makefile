# Generates output file in /pb. Loads proto file from proto/processor_message.go.
gen:
	protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative proto/*.proto

# Removes files under /pb
clean:
	rm pb/*.go

# runs main file.
run:
	go run main.go

# run all the tests
test: 
	go test -cover -race ./...