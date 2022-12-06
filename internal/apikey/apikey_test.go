package apikey_test

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkhq/tkcli/internal/apikey"
)

func Test_FromTkPrivateKey(t *testing.T) {
	// This private key is taken from an openSSL-generated PEM key:
	// 	$ openssl ec -in docs/fixtures/private_key.pem -noout -text
	// 	Private-Key: (256 bit)
	// 	priv:
	// 		48:7f:36:1d:df:d7:34:40:e7:07:f4:da:a6:77:5b:
	// 		37:68:59:e8:a3:c9:f2:9b:3b:b6:94:a1:29:27:c0:
	// 		21:3c
	// 	pub:
	// 		04:f7:39:f8:c7:7b:32:f4:d5:f1:32:65:86:1f:eb:
	// 		d7:6e:7a:9c:61:a1:14:0d:29:6b:8c:16:30:25:08:
	// 		87:03:16:c2:49:70:ad:78:11:cc:d9:da:7f:1b:88:
	// 		f2:02:be:ba:c7:70:66:3e:f5:8b:a6:83:46:18:6d:
	// 		d7:78:20:0d:d4
	// 	ASN1 OID: prime256v1
	// 	NIST CURVE: P-256
	privateKeyFromOpenSSL := "487f361ddfd73440e707f4daa6775b376859e8a3c9f29b3bb694a12927c0213c"
	apiKey, err := apikey.FromTkPrivateKey(privateKeyFromOpenSSL)
	assert.Nil(t, err)

	// This value was computed based on an openssl-generated PEM file:
	//   $ openssl ec -in docs/fixtures/private_key.pem -pubout -conv_form compressed -outform der | tail -c 33 | xxd -p -c 33
	//   read EC key
	//   writing EC key
	//   02f739f8c77b32f4d5f13265861febd76e7a9c61a1140d296b8c16302508870316
	expectedPublicKey := "02f739f8c77b32f4d5f13265861febd76e7a9c61a1140d296b8c16302508870316"
	assert.Equal(t, expectedPublicKey, apiKey.TkPublicKey)
}

func Test_Sign(t *testing.T) {
	tkPrivateKey := "487f361ddfd73440e707f4daa6775b376859e8a3c9f29b3bb694a12927c0213c"
	tkPublicKey := "02f739f8c77b32f4d5f13265861febd76e7a9c61a1140d296b8c16302508870316"

	apiKey, err := apikey.FromTkPrivateKey(tkPrivateKey)
	assert.Nil(t, err)

	sig, err := apikey.Sign("hello", apiKey)
	assert.Nil(t, err)
	sigBytes, err := hex.DecodeString(sig)
	assert.Nil(t, err)

	publicKey, err := apikey.DecodeTKPublicKey(tkPublicKey)
	assert.Nil(t, err)

	// Verify the soundness of the hash:
	//   $ echo -n 'hello' | shasum -a256
	//   2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824  -
	msgHash := sha256.Sum256([]byte("hello"))
	assert.Equal(t, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", fmt.Sprintf("%x", msgHash))

	// Finally, check the signature itself
	verifiedSig := ecdsa.VerifyASN1(publicKey, msgHash[:], sigBytes)
	assert.True(t, verifiedSig)
}

func Test_EncodedKeySizeIsFixed(t *testing.T) {
	for i := 0; i < 1000; i++ {
		apiKey, err := apikey.NewTkApiKey()
		assert.Nil(t, err)
		assert.Equal(t, 66, len(apiKey.TkPublicKey), "attempt %d: expected 66 characters for public key %s", i, apiKey.TkPublicKey)
		assert.Equal(t, 64, len(apiKey.TkPrivateKey), "attempt %d: expected 64 characters for private key %s", i, apiKey.TkPrivateKey)
	}
}
