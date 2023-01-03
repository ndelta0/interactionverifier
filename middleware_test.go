// Package interactionverifier_test Testing
package interactionverifier_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ndelta0/interactionverifier"
)

func TestDemo(t *testing.T) {
	cfg := interactionverifier.CreateConfig()
	cfg.PublicKey = "7c53f11c27f19405abcde3df2a034c3b89f57abcde0bd26b7114d25b675c9ae9"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := interactionverifier.New(ctx, next, cfg, "interaction-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	stringJSON := `{"type": 1}`

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost", bytes.NewBufferString(stringJSON))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)
}
