package videocoin

import "os"

const envVar = "VIDEOCOIN_APPLICATION_CREDENTIALS"

// FindDefaultCredentials searches for "Application Default Credentials".
//
// It looks for credentials in the following places,
// preferring the first location found:
//
//   1. A JSON file whose path is specified by the
//      VIDEOCOIN_APPLICATION_CREDENTIALS environment variable.
//
func FindDefaultCredentials() string {
	// First, try the environment variable.
	return os.Getenv(envVar)
}
