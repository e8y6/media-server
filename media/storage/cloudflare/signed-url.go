package cloudflare

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"../../../config"
	"../../../misc/exceptions"
	"../../../misc/log"

	jose "gopkg.in/square/go-jose.v2"
	jwt "gopkg.in/square/go-jose.v2/jwt"
)

type claims struct {
	KeyID     string          `json:"kid,omitempty"`
	VideoID   string          `json:"sub,omitempty"`
	Expiry    jwt.NumericDate `json:"exp,omitempty"`
	NotBefore jwt.NumericDate `json:"nbf,omitempty"`
}

type KeyInfo struct {
	ID            string `json:"id"`
	PrivateKeyPem string `json:"pem"`
	PrivateKeyJWK string `json:"jwk"`
}

var keyInfo KeyInfo

const expiresIn = time.Hour

func ObtainSigningKeys() {

	url := "https://api.cloudflare.com/client/v4/accounts/" + config.CF_ACCOUNT_ID + "/stream/keys"
	log.Info("CF: Obtaining Signing keys from Cloudflare.")

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		panic("CF: error ocurred while creating client for getting signing keys.")
	}
	req.Header.Add("X-Auth-Email", config.CF_EMAIL)
	req.Header.Add("X-Auth-Key", config.CF_STREAM_TOKEN)

	res, err := client.Do(req)
	if err != nil {
		panic("CF: error ocurred while creating sending request for signing keys.")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic("CF: error ocurred while reading body from request for signing keys.")
	}

	var jsonDecoded map[string]json.RawMessage
	err = json.Unmarshal(body, &jsonDecoded)
	if err != nil {
		panic("CF: error ocurred while parsing result for signing keys. (Step 1)")
	}
	err = json.Unmarshal(jsonDecoded["result"], &keyInfo)
	if err != nil {
		panic("CF: error ocurred while parsing result for signing keys. (Step 2)")
	}

	log.Info("CF: Obtaining Signing keys from Cloudflare completed successfully.")
}

func GetSignedURL(videoID string) string {
	// Decode privateKey
	keyBytes, err := base64.StdEncoding.DecodeString(keyInfo.PrivateKeyPem)
	if err != nil {
		panic(exceptions.Exception{
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}
	block, _ := pem.Decode(keyBytes)
	if err != nil {
		panic(exceptions.Exception{
			Cause: fmt.Sprintf("failed to decode pem: %s\n", err),
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(exceptions.Exception{
			Cause: fmt.Sprintf("failed to parse key: %s\n", err),
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	// Prepare to sign
	var options jose.SignerOptions
	options.WithType("JWT").WithHeader("kid", keyInfo.ID)
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: rsaPrivateKey},
		&options)
	if err != nil {
		panic(exceptions.Exception{
			Cause: fmt.Sprintf("failed to initialize signer: %s\n", err),
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	// Sign a JWT
	builder := jwt.Signed(signer)
	builder = builder.Claims(claims{
		KeyID:   keyInfo.ID,
		VideoID: videoID,
		Expiry:  *jwt.NewNumericDate(time.Now().Add(expiresIn)),
	})
	token, err := builder.CompactSerialize()
	if err != nil {
		panic(exceptions.Exception{
			Cause: fmt.Sprintf("failed to get token: %s\n", err),
			Type:  exceptions.TYPE_INTERNAL_ERROR,
			Error: err,
		})
	}

	return "https://videodelivery.net/" + token + "/manifest/video.m3u8"
}
