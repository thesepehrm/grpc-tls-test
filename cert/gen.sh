rm *.pem

# Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 3650 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=CA/ST=Ontorio/L=Toronto/O=Test, Inc./OU=Labs/CN=localhost"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=CA/ST=Ontorio/L=Toronto/O=Test, Inc./OU=Labs/CN=server.com"

# Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 3650 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.conf

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text

# Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/C=CA/ST=Ontorio/L=Toronto/O=Test, Inc./OU=Labs/CN=specificclient.com"

# Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client-req.pem -days 3650 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.conf

echo "Client's signed certificate"
openssl x509 -in client-cert.pem -noout -text
