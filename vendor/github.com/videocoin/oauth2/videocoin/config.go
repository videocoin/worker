package videocoin

import (
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2/jwt"
)

const serviceAccountKey = "service_account"

// credentialsFile is the unmarshalled representation of a credentials file.
type credentialsFile struct {
	Type string `json:"type"` // serviceAccountKey

	// Service Account fields
	ClientID     string `json:"client_id"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
}

func (f *credentialsFile) jwtConfig(scopes []string) *jwt.Config {
	cfg := &jwt.Config{
		Subject:      f.ClientID,
		PrivateKey:   []byte(f.PrivateKey),
		PrivateKeyID: f.PrivateKeyID,
		Scopes:       scopes,
	}
	return cfg
}

// JWTConfigFromJSON uses a VideoCoin service account JSON key file to read
// the credentials that authorize and authenticate the requests.
func JWTConfigFromJSON(jsonKey []byte, scope ...string) (*jwt.Config, error) {
	var f credentialsFile
	if err := json.Unmarshal(jsonKey, &f); err != nil {
		return nil, err
	}
	if f.Type != serviceAccountKey {
		return nil, fmt.Errorf("oauth2: read JWT from JSON credentials: 'type' field is %q (expected %q)", f.Type, serviceAccountKey)
	}
	scope = append([]string(nil), scope...) // copy
	return f.jwtConfig(scope), nil
}
