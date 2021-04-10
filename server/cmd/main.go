package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/thesepehrm/grpc-tls-test/pb/greeter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	logger = logrus.StandardLogger()
)

const (
	clientCAFile   = "cert/ca-cert.pem"
	serverKeyFile  = "cert/server-key.pem"
	serverCertFile = "cert/server-cert.pem"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := ioutil.ReadFile(clientCAFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

type greeterServer struct {
}

func (s *greeterServer) Hello(ctx context.Context, req *greeter.HelloRequest) (*greeter.HelloResponse, error) {

	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authentication")
	}

	tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
	commonName := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName

	return &greeter.HelloResponse{
		Reply: "Hi!, " + commonName,
	}, nil
}

func main() {

	port := flag.String("port", "", "the server port")
	flag.Parse()

	if port == nil {
		logger.Fatal("port is not set")
	}

	opts := []grpc.ServerOption{}

	creds, err := loadTLSCredentials()

	if err != nil {
		logger.Fatal(err)
	}

	opts = append(opts, grpc.Creds(creds))

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		logger.Fatal(err)
	}
	defer lis.Close()

	logger.Infof("Server is listening on port %s", *port)

	s := grpc.NewServer(opts...)
	logger.Infof("Starting gRPC services")

	greeter.RegisterGreeterServer(s, &greeterServer{})

	if err = s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %s", err)
	}

}
