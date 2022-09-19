package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
)

func main() {
	filename := "key"
	bitSize := 4096
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Fatal()
	}
	pub := key.Public()
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)
	if err := ioutil.WriteFile(filename+".rsa", keyPEM, 0700); err != nil {
		log.Fatal()
	}
	if err := ioutil.WriteFile(filename+".rsa.pub", pubPEM, 0700); err != nil {
		log.Fatal()
	}
}
