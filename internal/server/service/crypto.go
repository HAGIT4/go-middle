package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func (sv *MetricService) GetPrivateKey() *rsa.PrivateKey {
	return sv.privateKey
}

func GetPrivateKeyFromPem(path string) (prv *rsa.PrivateKey, err error) {
	var b []byte
	if path == "" {
		return nil, nil
	}
	if b, err = os.ReadFile(path); err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block.Type != "RSA PRIVATE KEY" || block == nil {
		return nil, errors.New("api: Bad private key")
	}
	prv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return prv, nil
}

func (sv *MetricService) DecryptWithPrivateKey(cText []byte) (dText []byte, err error) {
	hash := sha512.New()
	dText, err = rsa.DecryptOAEP(hash, rand.Reader, sv.privateKey, cText, nil)
	if err != nil {
		return nil, err
	}
	return dText, nil
}
