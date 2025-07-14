package core

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"lan-cloud/internal/shared"
	"time"
	"github.com/patrickmn/go-cache"
)

var (
	// Using auto cleanup map to prevent against resource exhaustion attack
	NonceCache = cache.New(
		time.Duration(shared.ServerPassConfig.NonceExpiry)*time.Millisecond, 
		time.Duration(shared.ServerPassConfig.NonceExpiry * 2)*time.Minute,
	)

	SessionKeyCache = cache.New(
		time.Duration(20)*time.Minute, 
		time.Duration(30)*time.Minute,
	)
)

func GenerateNonce() (sessionId string, nonceStr string, err error) {

	sessionId, err = retry(50, func() (string, error) {
		return generateSecureRandomString(12)
	})
	if err != nil { return "", "", fmt.Errorf("failed to generate session id")}

	nonceStr, err = retry(50, func() (string, error) {
		return generateSecureRandomString(16)
	})
	if err != nil { return "", "", fmt.Errorf("failed to generate nonceStr") }

	NonceCache.Set(sessionId, nonceStr, cache.DefaultExpiration) // save nonce

	return sessionId, nonceStr, nil
}

func generateSecureRandomString(n int) (string, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(b), nil
}

func retry(times int, fn func() (string, error)) (string, error) {
	for i := 0; i < times; i++ {
		if value, err := fn(); err == nil {
			return value, nil
		}
	}
	return "", fmt.Errorf("failed after %d attempts", times)
}

func GenerateSessionKey() (sessionId string, sessionKey string, err error) {

	sessionId, err = retry(50, func() (string, error) {
		return generateSecureRandomString(12)
	})
	if err != nil { return "", "", fmt.Errorf("failed to generate session id")}

	sessionKey, err = retry(50, func() (string, error) {
		return generateSecureRandomString(32)
	})
	if err != nil { return "", "", fmt.Errorf("failed to generate sessionKey") }

	// convert base64 string to byte array
	sessionKeyByte, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil { return "", "", fmt.Errorf("failed to decode sessionKey") }

	SessionKeyCache.Set(sessionId, sessionKeyByte, cache.DefaultExpiration)

	return sessionId, sessionKey, nil
}