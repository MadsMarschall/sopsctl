package remove

import (
	"errors"
	"strings"
	"testing"

	"filippo.io/age"
)

// MockSopsKeyManager is a mock implementation of domain.SopsKeyManager for testing
type MockSopsKeyManager struct {
	keys         []string
	listKeysErr  error
	removeKeyErr error
	removedKeys  []string
}

func (m *MockSopsKeyManager) GetIdentityCurrentCtx() (age.Identity, error) { return nil, nil }
func (m *MockSopsKeyManager) GetPrivateKey(ctxName string) (string, error) { return "", nil }
func (m *MockSopsKeyManager) GetPublicKey(ctxName string) (string, error)  { return "", nil }
func (m *MockSopsKeyManager) AddKeyFromCluster(ctxName string, namespace string, secretName string, secretKey string) (string, error) {
	return "", nil
}

func (m *MockSopsKeyManager) ListContextsWithKeys() ([]string, error) {
	if m.listKeysErr != nil {
		return nil, m.listKeysErr
	}
	return m.keys, nil
}

func (m *MockSopsKeyManager) RemoveKeyForContext(ctx string) error {
	if m.removeKeyErr != nil {
		return m.removeKeyErr
	}
	m.removedKeys = append(m.removedKeys, ctx)
	return nil
}

func TestNewKeyRemoveCmd(t *testing.T) {
	skm := &MockSopsKeyManager{}
	cmd := NewKeyRemoveCmd(skm)
	if cmd == nil {
		t.Fatal("expected non-nil KeyRemoveCmd")
	}
}

func TestKeyRemoveCmd_Execute_RemovesSingleKey(t *testing.T) {
	// Arrange
	skm := &MockSopsKeyManager{keys: []string{"production", "staging", "development"}}
	cmd := NewKeyRemoveCmd(skm)
	cmd.options = NewKeyRemoveCmdOptions(false, []string{"production"})

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(skm.removedKeys) != 1 {
		t.Fatalf("expected 1 key removed, got %d", len(skm.removedKeys))
	}
	if skm.removedKeys[0] != "production" {
		t.Fatalf("expected 'production' to be removed, got %q", skm.removedKeys[0])
	}
	if !strings.Contains(result, "production") {
		t.Fatalf("expected result to contain 'production', got %q", result)
	}
}

func TestKeyRemoveCmd_Execute_RemovesMultipleSpecificKeys(t *testing.T) {
	// Arrange
	skm := &MockSopsKeyManager{keys: []string{"production", "staging", "development"}}
	cmd := NewKeyRemoveCmd(skm)
	cmd.options = NewKeyRemoveCmdOptions(false, []string{"production", "staging"})

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(skm.removedKeys) != 2 {
		t.Fatalf("expected 2 keys removed, got %d", len(skm.removedKeys))
	}
	if !strings.Contains(result, "production") || !strings.Contains(result, "staging") {
		t.Fatalf("expected result to contain removed keys, got %q", result)
	}
}

func TestKeyRemoveCmd_Execute_RemovesAllKeysWithFlag(t *testing.T) {
	// Arrange
	skm := &MockSopsKeyManager{keys: []string{"production", "staging", "development"}}
	cmd := NewKeyRemoveCmd(skm)
	cmd.options = NewKeyRemoveCmdOptions(true, []string{}) // RemoveAll = true

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(skm.removedKeys) != 3 {
		t.Fatalf("expected 3 keys removed, got %d", len(skm.removedKeys))
	}
	if !strings.Contains(result, "production") {
		t.Fatalf("expected result to contain 'production', got %q", result)
	}
	if !strings.Contains(result, "staging") {
		t.Fatalf("expected result to contain 'staging', got %q", result)
	}
	if !strings.Contains(result, "development") {
		t.Fatalf("expected result to contain 'development', got %q", result)
	}
}

func TestKeyRemoveCmd_Execute_NoKeysFound(t *testing.T) {
	// Arrange
	skm := &MockSopsKeyManager{keys: []string{}}
	cmd := NewKeyRemoveCmd(skm)
	cmd.options = NewKeyRemoveCmdOptions(true, []string{})

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(result, "No SOPS keys found") {
		t.Fatalf("expected 'No SOPS keys found' message, got %q", result)
	}
}

func TestKeyRemoveCmd_Execute_ListKeysError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("list keys error")
	skm := &MockSopsKeyManager{listKeysErr: expectedErr}
	cmd := NewKeyRemoveCmd(skm)
	cmd.options = NewKeyRemoveCmdOptions(true, []string{})

	// Act
	_, err := cmd.Execute()

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "list keys") {
		t.Fatalf("expected 'list keys' in error, got: %v", err)
	}
}

func TestKeyRemoveCmd_Execute_RemoveKeyError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("remove key error")
	skm := &MockSopsKeyManager{
		keys:         []string{"production"},
		removeKeyErr: expectedErr,
	}
	cmd := NewKeyRemoveCmd(skm)
	cmd.options = NewKeyRemoveCmdOptions(false, []string{"production"})

	// Act
	_, err := cmd.Execute()

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "failed to remove SOPS key") {
		t.Fatalf("expected 'failed to remove SOPS key' in error, got: %v", err)
	}
}

func TestKeyRemoveCmd_Execute_KeyNotFoundMessage(t *testing.T) {
	// Arrange
	skm := &MockSopsKeyManager{keys: []string{"production", "staging"}}
	cmd := NewKeyRemoveCmd(skm)
	cmd.options = NewKeyRemoveCmdOptions(false, []string{"nonexistent"})

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(skm.removedKeys) != 0 {
		t.Fatalf("expected 0 keys removed, got %d", len(skm.removedKeys))
	}
	if !strings.Contains(result, "Key not found: nonexistent") {
		t.Fatalf("expected 'Key not found: nonexistent' message, got %q", result)
	}
}

func TestKeyRemoveCmd_Execute_MixedFoundAndNotFound(t *testing.T) {
	// Arrange
	skm := &MockSopsKeyManager{keys: []string{"production", "staging"}}
	cmd := NewKeyRemoveCmd(skm)
	cmd.options = NewKeyRemoveCmdOptions(false, []string{"production", "nonexistent"})

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(skm.removedKeys) != 1 {
		t.Fatalf("expected 1 key removed, got %d", len(skm.removedKeys))
	}
	if skm.removedKeys[0] != "production" {
		t.Fatalf("expected 'production' to be removed, got %q", skm.removedKeys[0])
	}
	if !strings.Contains(result, "Removed SOPS key for context: ") {
		t.Fatalf("expected removal message in result, got %q", result)
	}
	if !strings.Contains(result, "Key not found: nonexistent") {
		t.Fatalf("expected 'Key not found: nonexistent' message, got %q", result)
	}
}
