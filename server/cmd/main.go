package main

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/johanbrandhorst/certify"
	"github.com/johanbrandhorst/certify/issuers/vault"
	"github.com/sirupsen/logrus"
	"github.com/thesepehrm/grpc-tls-test/pb/greeter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	logrusadapter "logur.dev/adapter/logrus"
)

var (
	logger = logrus.StandardLogger()
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		logger.Fatalf("%s env variable is not set", name)
	}
	return v
}

type RSA struct {
	bits int
}

func (r RSA) Generate() (crypto.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, r.bits)
}

func vaultCert(f string) (credentials.TransportCredentials, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("vaultCert: problem with input file")
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("vaultCert: failed to append certificates")
	}
	issuer := &vault.Issuer{
		URL: &url.URL{
			Scheme: "https",
			Host:   "localhost:8200",
		},
		TLSConfig: &tls.Config{
			RootCAs: cp,
		},
		Token: getenv("TOKEN"),
		Role:  "certissuer",
	}
	cfg := certify.CertConfig{
		SubjectAlternativeNames: []string{"localhost"},
		IPSubjectAlternativeNames: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},
		KeyGenerator: RSA{bits: 2048},
	}
	c := &certify.Certify{
		CommonName:  "localhost",
		Issuer:      issuer,
		Cache:       certify.NewMemCache(),
		CertConfig:  &cfg,
		RenewBefore: 24 * time.Hour,
		Logger:      logrusadapter.New(logger),
	}
	tlsConfig := &tls.Config{
		GetCertificate: c.GetCertificate,
	}
	return credentials.NewTLS(tlsConfig), nil
}

type greeterServer struct {
}

func (s *greeterServer) Hello(context.Context, *greeter.HelloRequest) (*greeter.HelloResponse, error) {
	return &greeter.HelloResponse{
		Reply: "Hi!",
	}, nil
}

func main() {

	port := flag.String("port", "", "the server port")
	flag.Parse()

	if port == nil {
		logger.Fatal("port is not set")
	}

	opts := []grpc.ServerOption{}

	creds, err := vaultCert("cert/ca-cert.pem")
	if err != nil {
		logger.Fatal(err)
	}
	opts = append(opts, grpc.Creds(creds))

	lis, err := net.Listen("tcp", "127.0.0.1:"+*port)
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
