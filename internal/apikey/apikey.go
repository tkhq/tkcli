package apikey

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
)

// Struct to hold both serialized and ecdsa lib friendly version of a public/private key pair
type ApiKey struct {
	TkPrivateKey string
	TkPublicKey  string

	// Underlying ECDSA keypair
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

const TURNKEY_API_SIGNATURE_SCHEME = "SIGNATURE_SCHEME_TK_API_P256"

type ApiStamp struct {
	// API public key, hex-encoded
	PublicKey string `json:"publicKey"`

	// Signature is the P-256 signature bytes, hex-encoded
	Signature string `json:"signature"`

	// Signature scheme. Must be set to "SIGNATURE_SCHEME_TK_API_P256"
	Scheme string `json:"scheme"`
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
func EncodePublicKey(publicKey *ecdsa.PublicKey) string {
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
	return fmt.Sprintf("%s%064x", prefix, publicKey.X)
}

// Takes an ECDSA private key and create a new TkApiKey.
// Assumes that privateKey.PublicKey has been derived.
func FromEcdsaPrivateKey(privateKey *ecdsa.PrivateKey) (*ApiKey, error) {
	publicKey := &privateKey.PublicKey

	return &ApiKey{
		TkPrivateKey: EncodePrivateKey(privateKey),
		TkPublicKey:  EncodePublicKey(publicKey),
		publicKey:    publicKey,
		privateKey:   privateKey,
	}, nil
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

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}, nil
}

// Stamp generates a signing stamp for the given message with the given API key.
// The resulting stamp should be added as the "X-Stamp" header of an API request.
func Stamp(message []byte, apiKey *ApiKey) (out string, err error) {
	hash := sha256.Sum256(message)

	sigBytes, err := ecdsa.SignASN1(rand.Reader, apiKey.privateKey, hash[:])
	if err != nil {
		return
	}

	stamp := ApiStamp{
		PublicKey: apiKey.TkPublicKey,
		Signature: hex.EncodeToString(sigBytes),
		Scheme:    TURNKEY_API_SIGNATURE_SCHEME,
	}

	jsonStamp, err := json.Marshal(stamp)
	if err != nil {
		return "", errors.Wrap(err, "failed to encode API stamp as JSON")
	}

	return base64.RawURLEncoding.EncodeToString(jsonStamp), nil
}
