package agent

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func (a *agent) GetPublicKeyFromPem(path string) (pub *rsa.PublicKey, err error) {
	if path == "" {
		return nil, nil
	}
	var b []byte
	if b, err = os.ReadFile(path); err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block.Type != "RSA PUBLIC KEY" || block == nil {
		return nil, errors.New("agent: Bad public key")
	}
	pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}

func (a *agent) EncryptWithPublicKey(msg []byte) (cText []byte, err error) {
	hash := sha512.New()
	cText, err = rsa.EncryptOAEP(hash, rand.Reader, a.publicKey, msg, nil)
	if err != nil {
		return nil, err
	}
	return cText, nil
}
