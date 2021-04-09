gen: 
	protoc --proto_path=proto proto/*.proto --go_opt=module=github.com/thesepehrm/grpc-tls-test --go_out=plugins=grpc:.
clean:
	rm -rf pb/*

run-client:
	go run client/cmd/main.go

run-server:
	go run server/cmd/main.go