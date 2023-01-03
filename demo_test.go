package interactionplugin_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ndelta0/traefik-discord-interaction-verifier"
)

func TestDemo(t *testing.T) {
	cfg := interactionplugin.CreateConfig()
	cfg.PublicKey = "7c53f11c27f19405c8ad63df2a034c3b89f57fdbfb0bd26b7114d25b675c9ae9"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := interactionplugin.New(ctx, next, cfg, "interaction-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	stringJson := `{"type": 1}`

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost", bytes.NewBufferString(stringJson))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)
}
