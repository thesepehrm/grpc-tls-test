gen: 
	protoc --proto_path=proto proto/*.proto --go_opt=module=github.com/thesepehrm/grpc-tls-test --go_out=plugins=grpc:.
clean:
	rm -rf pb/*

cert:
	cd ./cert; bash gen.sh; cd ..

client:
	go run client/cmd/main.go --address localhost:5000

server:
	go run server/cmd/main.go --port 5000

run-vault:
	vault server -config=vault/vault_config.hcl

vault-init:
	export VAULT_ADDR=https://localhost:8200
	export VAULT_CACERT=cert/ca-cert.pem
	vault operator init

.PHONY: gen clean server client cert 