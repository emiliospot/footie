package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// VerifyHMACSignature verifies the HMAC SHA256 signature of the webhook request.
// This is a generic implementation that providers can use or override.
func VerifyHMACSignature(payload []byte, signature string, secret string) bool {
	if secret == "" {
		return true // No secret configured, allow in development
	}
	if signature == "" {
		return false
	}

	// Compute expected signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures (constant-time comparison)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
