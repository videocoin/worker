package videocoin

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/videocoin/oauth2/helpers"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jws"
)

func JWTAccessTokenSource(filename string, audience string) (oauth2.TokenSource, error) {
	jsonKey, err := readCredentialsFile(filename)
	if err != nil {
		return nil, err
	}
	return JWTAccessTokenSourceFromJSON(jsonKey, audience)
}

func JWTAccessTokenSourceFromJSON(jsonKey []byte, audience string) (oauth2.TokenSource, error) {
	cfg, err := JWTConfigFromJSON(jsonKey)
	if err != nil {
		return nil, fmt.Errorf("jwt: could not parse JSON key: %v", err)
	}
	pk, err := helpers.ParseKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("jwt: could not parse key: %v", err)
	}
	ts := &jwtAccessTokenSource{
		subject:  cfg.Subject,
		audience: audience,
		pk:       pk,
		pkID:     cfg.PrivateKeyID,
	}
	tok, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return oauth2.ReuseTokenSource(tok, ts), nil
}

type jwtAccessTokenSource struct {
	subject, audience string
	pk                *rsa.PrivateKey
	pkID              string
}

func (ts *jwtAccessTokenSource) Token() (*oauth2.Token, error) {
	iat := time.Now()
	exp := iat.Add(time.Hour)
	cs := &jws.ClaimSet{
		Iss: ts.subject,
		Sub: ts.subject,
		Aud: ts.audience,
		Iat: 1,
		Exp: exp.Unix(),
	}
	hdr := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
		KeyID:     string(ts.pkID),
	}
	msg, err := jws.Encode(hdr, cs, ts.pk)
	if err != nil {
		return nil, fmt.Errorf("jwt: could not encode JWT: %v", err)
	}
	return &oauth2.Token{AccessToken: msg, TokenType: "Bearer", Expiry: exp}, nil
}

func readCredentialsFile(filename string) ([]byte, error) {
	if filename == "" {
		if filename = FindDefaultCredentials(); filename == "" {
			fmt.Errorf("jwt: credentials not found")
		}
	}
	return ioutil.ReadFile(filename)
}
