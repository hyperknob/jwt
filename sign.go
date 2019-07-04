package jwt

import (
	"encoding/base64"
	"encoding/json"
)

// SignOption is a functional option for signing.
type SignOption func(*Header)

// Sign signs a payload with alg.
func Sign(payload interface{}, alg Algorithm, opts ...SignOption) ([]byte, error) {
	var hd Header
	for _, opt := range opts {
		opt(&hd)
	}
	// Override some values or set them if empty.
	hd.Algorithm = alg.Name()
	hd.Type = "JWT"
	// Marshal the header part of the JWT.
	hb, err := json.Marshal(hd)
	if err != nil {
		return nil, err
	}
	// Marshal the claims part of the JWT.
	pb, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	enc := base64.RawURLEncoding
	h64len := enc.EncodedLen(len(hb))
	p64len := enc.EncodedLen(len(pb))
	sig64len := enc.EncodedLen(alg.Size())
	token := make([]byte, h64len+1+p64len+1+sig64len)

	enc.Encode(token, hb)
	token[h64len] = '.'
	enc.Encode(token[h64len+1:], pb)
	sig, err := alg.Sign(token[:h64len+1+p64len])
	if err != nil {
		return nil, err
	}
	token[h64len+1+p64len] = '.'
	enc.Encode(token[h64len+1+p64len+1:], sig)
	return token, nil
}

// ContentType sets the "cty" claim for a Header before signing.
func ContentType(cty string) SignOption {
	return func(hd *Header) {
		hd.ContentType = cty
	}
}

// KeyID sets the "kid" claim for a Header before signing.
func KeyID(kid string) SignOption {
	return func(hd *Header) {
		hd.KeyID = kid
	}
}
