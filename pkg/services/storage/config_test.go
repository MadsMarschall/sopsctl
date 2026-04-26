package storage

import (
	"sopsctl/pkg/domain"
	"strings"
	"testing"
)

func TestConfigFile_GetPrivateKey_ReturnsErrorForUnknownContext(t *testing.T) {
	// Arrange
	uut := &ConfigFile{Contexts: map[string]domain.CTX{}}

	// Act
	key, err := uut.GetPrivateKey("missing")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if key != "" {
		t.Fatalf("expected empty key, got %q", key)
	}
	if !strings.Contains(err.Error(), `"missing"`) {
		t.Fatalf("expected error to name the context %q, got: %v", "missing", err)
	}
	if !strings.Contains(err.Error(), "sopsctl create key") {
		t.Fatalf("expected error to point at 'sopsctl create key', got: %v", err)
	}
}

func TestConfigFile_GetPrivateKey_ReturnsErrorWhenStoredKeyIsEmpty(t *testing.T) {
	// Arrange
	uut := &ConfigFile{Contexts: map[string]domain.CTX{
		"prod": {Namespace: "flux-system", SecretName: "sops-age", KeyName: "age.agekey"},
	}}

	// Act
	key, err := uut.GetPrivateKey("prod")

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if key != "" {
		t.Fatalf("expected empty key, got %q", key)
	}
	if !strings.Contains(err.Error(), "sopsctl create key") {
		t.Fatalf("expected error to point at 'sopsctl create key', got: %v", err)
	}
}

func TestConfigFile_GetPrivateKey_ReturnsKeyWhenStored(t *testing.T) {
	// Arrange
	const stored = "AGE-SECRET-KEY-1EXAMPLE"
	uut := &ConfigFile{Contexts: map[string]domain.CTX{
		"prod": {PrivateKey: stored},
	}}

	// Act
	key, err := uut.GetPrivateKey("prod")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if key != stored {
		t.Fatalf("expected key %q, got %q", stored, key)
	}
}
