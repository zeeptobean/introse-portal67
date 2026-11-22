package authn

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const IdentityHeader = "X-Portal-Identity"

type Identity struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

type Session struct {
	Identity
	Version   int64     `json:"version"`
	ExpiresAt time.Time `json:"expires_at"`
}

type claims struct {
	Identity
	Audience  string `json:"aud"`
	ExpiresAt int64  `json:"exp"`
}

var ErrInvalidToken = errors.New("invalid service identity token")
var encoding = base64.RawURLEncoding

func ParseSigningKey(value string) (ed25519.PrivateKey, error) {
	key, err := encoding.DecodeString(value)
	if err != nil || len(key) != ed25519.SeedSize {
		return nil, errors.New("IDENTITY_SIGNING_KEY must contain a base64url Ed25519 seed")
	}
	return ed25519.NewKeyFromSeed(key), nil
}

func ParseVerificationKeys(value string) ([]ed25519.PublicKey, error) {
	var keys []ed25519.PublicKey
	for _, part := range strings.Split(value, ",") {
		key, err := encoding.DecodeString(strings.TrimSpace(part))
		if err != nil || len(key) != ed25519.PublicKeySize {
			return nil, errors.New("IDENTITY_PUBLIC_KEYS must contain base64url Ed25519 public keys")
		}
		keys = append(keys, ed25519.PublicKey(key))
	}
	return keys, nil
}

func Sign(key ed25519.PrivateKey, identity Identity, audience string) string {
	// Claims contain only primitive values and always marshal successfully.
	payload, _ := json.Marshal(claims{identity, audience, time.Now().Add(time.Minute).Unix()})
	body := encoding.EncodeToString(payload)
	return body + "." + encoding.EncodeToString(ed25519.Sign(key, []byte(body)))
}

func Verify(keys []ed25519.PublicKey, token, audience string) (Identity, error) {
	if len(token) > 1024 {
		return Identity{}, ErrInvalidToken
	}
	body, sig, ok := strings.Cut(token, ".")
	if !ok {
		return Identity{}, ErrInvalidToken
	}
	signature, err := encoding.DecodeString(sig)
	if err != nil {
		return Identity{}, ErrInvalidToken
	}
	for _, key := range keys {
		if ed25519.Verify(key, []byte(body), signature) {
			data, err := encoding.DecodeString(body)
			var payload claims
			if err != nil || json.Unmarshal(data, &payload) != nil ||
				payload.Audience != audience || payload.ExpiresAt <= time.Now().Unix() {
				return Identity{}, ErrInvalidToken
			}
			return payload.Identity, nil
		}
	}
	return Identity{}, ErrInvalidToken
}
