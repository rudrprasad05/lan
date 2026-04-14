package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
)

func GenerateKeyPair() (pub string, priv string, err error) {
	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", err
	}

	pub = base64.StdEncoding.EncodeToString(public)
	priv = base64.StdEncoding.EncodeToString(private)

	return
}

func DecodePublicKey(pub string) (ed25519.PublicKey, error) {
	return base64.StdEncoding.DecodeString(pub)
}

func DecodePrivateKey(priv string) (ed25519.PrivateKey, error) {
	return base64.StdEncoding.DecodeString(priv)
}

func Sign(message []byte, privKey string) ([]byte, error) {
	priv, err := DecodePrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	signature := ed25519.Sign(priv, message)
	return signature, nil
}

func Verify(message, signature []byte, pubKey string) bool {
	pub, err := DecodePublicKey(pubKey)
	if err != nil {
		return false
	}

	return ed25519.Verify(pub, message, signature)
}
