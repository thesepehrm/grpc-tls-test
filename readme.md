# GRPC mTLS Authentication

This repo is a simple implementation of mutual TLS authentication in GRPC microservices.

## Requirements

- [Protoc](https://grpc.io/docs/protoc-installation/) 
- [Vault](https://learn.hashicorp.com/vault/)

## Commands

- `make gen` and `make clean` to generate or clean protobuf codes
- `make cert` to generate certificates for the client, the server and Vault
- `make client` to run the client
- `make server` to run the server
- `make vault-init` to initialize the vault
- `make run-vault` to run the vault server


## Vault Setup

1. First, [setup the vault server](https://learn.hashicorp.com/vault/) and
2. Run `make vault-init` to initialize the operator. 
3. [Unseal the vault server](https://www.vaultproject.io/docs/concepts/seal).

*You can test the vault using this command:*
```bash
$ curl --cacert cert/ca-cert.pem \
    -i https://localhost:8200/v1/sys/health

HTTP/1.1 200 OK
```
4. Enable [Vault PKI Secrets Engine backend] using [these instructions](https://www.vaultproject.io/docs/secrets/pki/index.html).

Finally you can run the server and the client. 🎉

