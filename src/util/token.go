package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQDAMqgyO/g84N8GDjiWidVBJpft+kIZehoXsm7dx29yloGlg863
Imdl9Uwjf1th6h84It940Y6OI0/YiYHeO8/ye1ucpdlzfmRGOi9+wOiuWrNKJmG4
cEmUkze/uLTkFSUnsgD+anpJUUAQ6gOnNq6xUrxNK34mDvrawxtDP3f2UQIDAQAB
AoGAc63JqDqKBXI/KbDjhE+/R/BHn1dx802XaM3fhqKjxG8r5wf3IiiV3TsPsYnU
4ZD9a1cp89kFGS3NwAG7ZZvQYxAjHvHA4h/hCajxZgUcCzBiZ4JH0tdoM+u54BJO
mj/2cr6towI3yVF7DS5+xDw8JNAwlz1dRqrBSX3B3JW/gTECQQD3HdgE58EsDp3X
qzISVxVd90F/h/pl9hu4+VXD6wwxyc2fh3HBzgLoMBnub/CrpPs8ciCfx1asgUQT
O9kIFQ/lAkEAxxtow0pj77RKyr14et5cJzo9b/fyPQGurnNFKbXerEwW0krPZaHP
QgJjyOwpRIi4oq+1Qknwdtj2/G78Tawt/QJAPJgeziUd4vW6kpWx83lTDfWBJApt
xe6xIYxSEXZjSRoYx5Nou4MOh2y0Dxl3xD7yNAIwKb2xbR9NWAIG18qCWQJAFH96
4pgW/8eE56hn7eZUgGlbh9pz4tn4fNc7KJcjrINM2it/fIwTBU2vrjC58udMcts6
AvAPxHyDuOtIKErwlQJAI7OkRAqgTSFphu61whp7a4vE8pXuEdm01DIcmbWTw0OR
JwMW3WJKH9YLBlxyGAe/PconqpQ5PW6HeOUmaW/hew==
-----END RSA PRIVATE KEY-----
`)

var priv *rsa.PrivateKey

func init() {
	fmt.Println("init")
	block, _ := pem.Decode(privateKey)
	if block == nil {
		log.Fatalln("no private key")
	}
	var err error
	priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalln("failed to get private key", err)
	}
}

func VerfiyToken(token string) bool {
	if token == "" {
		return false
	}
	ciphertext, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}
	tokenPair, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return false
	}

	if string(tokenPair) == "BEYOU_DIAMOND" {
		return true
	}
	return false
}

func GenerateToken() (string, error) {
	publicKey, err := ioutil.ReadFile("token_pk.pem")
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	bs, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte("BEYOU_DIAMOND"))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bs), nil
}
