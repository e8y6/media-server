package cloudflare

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	jose "gopkg.in/square/go-jose.v2"
	jwt "gopkg.in/square/go-jose.v2/jwt"
)

type claims struct {
	KeyID     string          `json:"kid,omitempty"`
	VideoID   string          `json:"sub,omitempty"`
	Expiry    jwt.NumericDate `json:"exp,omitempty"`
	NotBefore jwt.NumericDate `json:"nbf,omitempty"`
}

// obtain this thing on boot
const keyID = "908b0c1ec99edeecc27523c407b3e07b"
const privateKey = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBdkdQNjkza2VPWjZocnNVUzZQTGlLdDBzNG9RTmpGWU9jT25hWU5aREZqMlNGUGs4CmZBZ0hNR0pVSHRoRU8zeFdGY3pKb1ZMdmpVM2RmMTVrZGlQRWdJK3ZFTlcrTnVsSGp0TzlTc2ZWUTdJcmdSYVYKNldibzF0OWlhY1ZhQ2FFQ2JJVzlmTU9UelJhK09RL0RldU8wY2dQRXZSb2h4TGdmOW1UZ25BVEpLUkwyaWplRgpaY3ZUK1RORVZOMldXdmJ2REJiOXlIQklSSVdTN0xtUmJSclYvMTBkNmRXWTRFUDkxTHNpTm9SNUJBaGd3TVFzCmtFcEFTdVVSRG1qSGVGdjMxVzEzRXJHV2JFNTR3Y3pPWXNqV1VBUmpLcm9EL3JlRFJMWEFEQlRVM0JvSm1yZ0gKMW8yRm9wbmxiRHluMDFpaE4yQ1RHMHVaZWFSbnVoekV5VWZwaVFJREFRQUJBb0lCQUdJN2kvaVF5a2JuUlkxNQpLNVFXV2dKWjAvYkZQcjlIZkQ0NlltbU9MK3NmN2RWTDVOTVQ2Sk85SWZuM2NSVEhqNmZNWWZMaDZSRjRZWi94Cm0xYlM3YnJQc1V5STk2ZHdXcVRLR2ZFdFpESHBiSy9pRkFkaFp0WHNJMGZkNVVZU0U4NThxa0t0Ukp0eldYc3QKa2hPNU9qVWRhKy9pK1dxM1M4dFI3S0RPQ2dxRTk1MkVPSHBBSUZDa1RYQjZEYWs5MUlQaFZzUUJUUlo1UisxaApHaXdiRDJHbktIVGtMdE1qeGVtTEFldzRjcWtvdUxnb3FYS1RWSkJTeGYzbDRMTWN4eU9rUWYxSDZZZjNobmVXCk84Q2F1SnMzWUtLd2g5dHV4dFRCdmUxVmpQcVVYRVFuNmFHT2QrcXhoamVoRXpUMlI0VjJXcWdCSnFxVHFMcFUKUys5b3ZaRUNnWUVBeG1DTXNCYkZ2SzZLT2VsSlNwNUtZTUloc0M4eDd3TGUzY1E4RGVaUnVGaHhhdHBWcW5WcwpqUk96cWlHcTE2ZDdMWXdmRzBEaXBkQitUL212c2oxSmt1djd2amM3SXRKNW92V2EydEdwZDhvQW03a3ZqQnlxCmYvTnBwWFAwRHlDZGRSV3hmSEVFQXd3QS9sTVAyRklEMXpvTXI5bGpobmxYL0VjMWVtaWJrTlVDZ1lFQTh4elMKSzRtV0VjVzFLeUxvK0FyRVlVUzFPYWxrNzRTOWo0SnNUVmtRZ3Y4VTZtMm5MRCtCNTRxRjZ5c2lTZzJ1SzR6cApiNlZ4b2h0YlY3WDdSMjRUbFZyZDc4V2NnRkp3WGJ2YjNjNU8wVWg1VVk5Y3pCbTAxcVNhaWtRekVGUm5aVnVxCkw1b3Zhby8rZ0JXYlpBclVqS2NjZndkYkpKbXFWQXZPMzlxSGIrVUNnWUFRSmhYdDA3eS9FbHBRUW94ajFhVHEKWVlOS2kxeitQdmFUaVFEMmhMUk1WRzdQS3Z6a1JuRFN2ekxWKzYyanBvK2hjcEdwcjB1RUFnZjJUTFlmeFZ0eQo2V0R4NkI4WlE1Y0JUQXNTR3hVM21pc1lnaWU3dVMyc0FzMnIyVmVaejZiaHZDVlpvdjJYbmVlS3pJb3lxdUtECko4ZVduUlM4QXE3RjY2b1B6K295RFFLQmdRQ2VEbTdWdjIzQndEZzVBMGxUZWl5UzBJakNKRTlyS1hIVWk0YSsKQXRtcFVRM1lHOUpFQWtZQ3N4eTkrQjBpNkNJRVRtaTJIV256YXJYSUlKaHRxRE16TnhCemdwWnRGeXZXeFB2OQpDTlJWWERrRHZhRW5VTTh4ZlhLNzBmc2c4cjhHTnNJK2tJK2pTRnErQzA0cmpKOERrdUlEOUZSZFIwcGx1Z2VWCm9nM0x3UUtCZ0RqV2FJM2kvNENoRCs2eU54UHZrRU9jSVRzL2J1QWo2dmhoSTZOT3cwZDNERVk0YnR0Y2V0bXIKZmd4YldNZXBFTVJJeXN4UkkyalBKVWJzUDFaOGR6cExEb081R0ZuKzZaSFJRRkNtOHdoUEF4amZMOWhub09LbgprZGZrWG55Zzk0SkM2QTF6VnJtNGJlZlhHczZocjRhZVh3ZU83SzlvSVc4ejR3RmtPSGtBCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg=="
const expiresIn = time.Hour

func GetSignedURL(videoID string) string {
	// Decode privateKey
	keyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		fmt.Printf("failed to generate key: %s\n", err)
		panic(err)
	}
	block, _ := pem.Decode(keyBytes)
	if err != nil {
		fmt.Printf("failed to decode pem: %s\n", err)
		panic(err)
	}
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Printf("failed to parse key: %s\n", err)
		panic(err)
	}

	// Prepare to sign
	var options jose.SignerOptions
	options.WithType("JWT").WithHeader("kid", keyID)
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: rsaPrivateKey},
		&options)
	if err != nil {
		fmt.Printf("failed to initialize signer: %s\n", err)
		panic(err)
	}

	// Sign a JWT
	builder := jwt.Signed(signer)
	builder = builder.Claims(claims{
		KeyID:   keyID,
		VideoID: videoID,
		Expiry:  *jwt.NewNumericDate(time.Now().Add(expiresIn)),
	})
	token, err := builder.CompactSerialize()
	if err != nil {
		fmt.Printf("failed to get token: %s\n", err)
		panic(err)
	}

	return "https://videodelivery.net/" + token + "/manifest/video.m3u8"
}
