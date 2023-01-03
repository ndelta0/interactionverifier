// Package interactionverifier Middleware plugin
package interactionverifier

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	PublicKey string `json:"publicKey,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		PublicKey: "",
	}
}

// Middleware a Middleware plugin.
type Middleware struct {
	next      http.Handler
	publicKey string
	name      string
}

// New created a new Middleware plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.PublicKey) == 0 {
		return nil, fmt.Errorf("publicKey cannot be empty")
	}

	return &Middleware{
		publicKey: config.PublicKey,
		next:      next,
		name:      name,
	}, nil
}

func (a *Middleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Get signature and timestamp from header
	signature := req.Header.Get("X-Signature-Ed25519")
	timestamp := req.Header.Get("X-Signature-Timestamp")

	// Reject if signature or timestamp is empty or doesn't exist
	if len(signature) == 0 || len(timestamp) == 0 {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get body
	body, err := req.GetBody()

	// Reject if body is empty or doesn't exist
	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Read body as string
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(body)
	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}
	bodyString := buf.String()

	// Create public key var
	key := ed25519.PublicKey(a.publicKey)

	// Get signature bytes
	sig, err := hex.DecodeString(signature)
	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Construct message
	message := timestamp + bodyString

	// Verify signature
	if !ed25519.Verify(key, []byte(message), sig) {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Call next middleware
	a.next.ServeHTTP(rw, req)
}
