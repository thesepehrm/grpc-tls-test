package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/thesepehrm/grpc-tls-test/pb/greeter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {

	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	transportOption := grpc.WithTransportCredentials(tlsCredentials)
	cc, err := grpc.Dial(*serverAddress, transportOption)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	greeterClient := greeter.NewGreeterClient(cc)
	res, err := greeterClient.Hello(context.Background(), &greeter.HelloRequest{
		Greeting: "hey",
	})
	if err != nil {
		log.Fatalf("cannot say hello to the server %s", err.Error())
	}

	log.Printf("[Server]: %s", res.GetReply())

}
