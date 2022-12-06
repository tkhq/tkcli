package apikey

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Struct to hold both serialized and ecdsa lib friendly version of a public/private key pair
type ApiKey struct {
	TkPrivateKey string
	TkPublicKey  string
	// Underlying ECDSA keypair
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// Create a new Turnkey API key
func NewTkApiKey() (*ApiKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	apiKey, err := FromEcdsaPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}

// Encode an ECDSA private key into the Turnkey format
// For now, "Turnkey format" = raw DER form
func EncodePrivateKey(privateKey *ecdsa.PrivateKey) string {
	return fmt.Sprintf("%064x", privateKey.D)
}

// Encode an ECDSA public key into the Turnkey format.
// For now, "Turnkey format" = standard compressed form for ECDSA keys
func EncodePublicKey(publicKey *ecdsa.PublicKey) (string, error) {
	// ANSI X9.62 point encoding
	var prefix string
	if publicKey.Y.Bit(0) == 0 {
		// Even Y
		prefix = "02"
	} else {
		// Odd Y
		prefix = "03"
	}

	// Encode the public key X coordinate as 64 hexadecimal characters, padded with zeroes as necessary
	encodedX := fmt.Sprintf("%064x", publicKey.X)
	compressedPubKey := prefix + encodedX

	return compressedPubKey, nil
}

// Takes an ECDSA private key and create a new TkApiKey.
// Assumes that privateKey.PublicKey has been derived.
func FromEcdsaPrivateKey(privateKey *ecdsa.PrivateKey) (*ApiKey, error) {
	publicKey := &privateKey.PublicKey

	tkPrivateKey := EncodePrivateKey(privateKey)

	tkPublicKey, err := EncodePublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	tkApiKey := ApiKey{
		TkPublicKey:  tkPublicKey,
		TkPrivateKey: tkPrivateKey,
		publicKey:    publicKey,
		privateKey:   privateKey,
	}

	return &tkApiKey, nil
}

// Takes a TK-encoded private key and creates an ECDSA private key
func FromTkPrivateKey(encodedPrivateKey string) (*ApiKey, error) {
	bytes, err := hex.DecodeString(encodedPrivateKey)
	if err != nil {
		return nil, err
	}

	dValue := new(big.Int).SetBytes(bytes)

	publicKey := new(ecdsa.PublicKey)
	privateKey := ecdsa.PrivateKey{
		PublicKey: *publicKey,
		D:         dValue,
	}

	// Derive the public key
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKey.D.Bytes())

	apiKey, err := FromEcdsaPrivateKey(&privateKey)
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}

// Takes a TK-encoded public key and creates an ECDSA public key
func DecodeTKPublicKey(encodedPublicKey string) (*ecdsa.PublicKey, error) {
	bytes, err := hex.DecodeString(encodedPublicKey)
	if err != nil {
		return nil, err
	}

	if len(bytes) != 33 {
		return nil, fmt.Errorf("expected a 33-bytes-long public key (compressed). Got %d bytes", len(bytes))
	}

	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), bytes)
	publicKey := new(ecdsa.PublicKey)
	publicKey.Curve = elliptic.P256()
	publicKey.X = x
	publicKey.Y = y

	return publicKey, nil
}

func Sign(message string, apiKey *ApiKey) (string, error) {
	hash := sha256.Sum256([]byte(message))

	sigBytes, err := ecdsa.SignASN1(rand.Reader, apiKey.privateKey, hash[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sigBytes), nil
}
