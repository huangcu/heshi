package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"errors"
)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDAMqgyO/g84N8GDjiWidVBJpft
+kIZehoXsm7dx29yloGlg863Imdl9Uwjf1th6h84It940Y6OI0/YiYHeO8/ye1uc
pdlzfmRGOi9+wOiuWrNKJmG4cEmUkze/uLTkFSUnsgD+anpJUUAQ6gOnNq6xUrxN
K34mDvrawxtDP3f2UQIDAQAB
-----END PUBLIC KEY-----
`)

func main() {
	s, _ := New()
	bs, err := rsaEncrypt([]byte(s))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(bs))
}

func New() (string, error) {
	data, err := rsaEncrypt([]byte("heshiservice"))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func rsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}
