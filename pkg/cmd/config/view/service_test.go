package view

import (
	"errors"
	"sopsctl/pkg/domain"
	"testing"
)

// MockKeyStorage is a mock implementation of domain.KeyStorage for testing
type MockKeyStorage struct {
	storageMode    domain.StorageMode
	storageModeErr error
}

func (m *MockKeyStorage) GetStorageMode() (domain.StorageMode, error) {
	if m.storageModeErr != nil {
		return "", m.storageModeErr
	}
	return m.storageMode, nil
}

func (m *MockKeyStorage) SetStorageMode(mode domain.StorageMode) error {
	m.storageMode = mode
	return nil
}

// Unused interface methods - implemented to satisfy the interface
func (m *MockKeyStorage) GetCtx(ctxName string) (*domain.CTX, error)      { return nil, nil }
func (m *MockKeyStorage) SavePrivateKey(key string, ctxName string) error { return nil }
func (m *MockKeyStorage) GetPrivateKey(ctxName string) (string, error)    { return "", nil }
func (m *MockKeyStorage) ListContextsWithKeys() ([]string, error)         { return nil, nil }
func (m *MockKeyStorage) RemoveKeyForContext(ctx string) error            { return nil }
func (m *MockKeyStorage) SaveCtxReference(ctxName string, namespace string, secretName string, key string) error {
	return nil
}

func TestNewConfigViewCmd(t *testing.T) {
	storage := &MockKeyStorage{}
	cmd := NewConfigViewCmd(storage)
	if cmd == nil {
		t.Fatal("expected non-nil ConfigViewCmd")
	}
}

func TestConfigViewCmd_Execute_ReturnsLocalMode(t *testing.T) {
	// Arrange
	storage := &MockKeyStorage{storageMode: domain.LocalStorageMode}
	cmd := NewConfigViewCmd(storage)

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	expected := "storage-mode: local"
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestConfigViewCmd_Execute_ReturnsClusterMode(t *testing.T) {
	// Arrange
	storage := &MockKeyStorage{storageMode: domain.InClusterStorageMode}
	cmd := NewConfigViewCmd(storage)

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	expected := "storage-mode: cluster"
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestConfigViewCmd_Execute_ReturnsErrorOnStorageFailure(t *testing.T) {
	// Arrange
	expectedErr := errors.New("storage error")
	storage := &MockKeyStorage{storageModeErr: expectedErr}
	cmd := NewConfigViewCmd(storage)

	// Act
	result, err := cmd.Execute()

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if result != "" {
		t.Fatalf("expected empty result, got %q", result)
	}
}

func TestConfigViewCmd_UseOptions_ReturnsExecutor(t *testing.T) {
	// Arrange
	storage := &MockKeyStorage{storageMode: domain.LocalStorageMode}
	cmd := NewConfigViewCmd(storage)

	// Act
	executor, err := cmd.UseOptions(nil, nil)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if executor == nil {
		t.Fatal("expected non-nil executor")
	}
}
