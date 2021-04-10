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

## Setup

Just run `make cert` to generate the certificates and it will be ready to run 🎉
