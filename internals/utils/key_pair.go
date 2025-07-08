package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
)

type EncodedKeyPair struct {
	PrivateKey string `json:"private_key"` // base64
	PublicKey  string `json:"public_key"`  // base64
}

func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func EncodeKeyPair(priv *ecdsa.PrivateKey) (*EncodedKeyPair, error) {
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, err
	}

	return &EncodedKeyPair{
		PrivateKey: base64.StdEncoding.EncodeToString(privBytes),
		PublicKey:  base64.StdEncoding.EncodeToString(pubBytes),
	}, nil
}

func DecodePublicKey(encoded string) (*ecdsa.PublicKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	pub, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return nil, err
	}

	return pub.(*ecdsa.PublicKey), nil
}

func DecodePrivateKey(encoded string) (*ecdsa.PrivateKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	privKey, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		return nil, err
	}
	return privKey, nil
}
