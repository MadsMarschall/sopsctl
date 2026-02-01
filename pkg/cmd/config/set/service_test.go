package set

import (
	"errors"
	"sopsctl/pkg/domain"
	"strings"
	"testing"
)

// MockKeyStorage is a mock implementation of domain.KeyStorage for testing
type MockKeyStorage struct {
	storageMode    domain.StorageMode
	storageModeErr error
	setModeErr     error
}

func (m *MockKeyStorage) GetStorageMode() (domain.StorageMode, error) {
	if m.storageModeErr != nil {
		return "", m.storageModeErr
	}
	return m.storageMode, nil
}

func (m *MockKeyStorage) SetStorageMode(mode domain.StorageMode) error {
	if m.setModeErr != nil {
		return m.setModeErr
	}
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

func TestNewConfigSetCmd(t *testing.T) {
	storage := &MockKeyStorage{}
	cmd := NewConfigSetCmd(storage)
	if cmd == nil {
		t.Fatal("expected non-nil ConfigSetCmd")
	}
}

func TestConfigSetCmd_UseOptions_RequiresTwoArgs(t *testing.T) {
	// Arrange
	storage := &MockKeyStorage{}
	cmd := NewConfigSetCmd(storage)

	testCases := []struct {
		name string
		args []string
	}{
		{"no args", []string{}},
		{"one arg", []string{"storage-mode"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			_, err := cmd.UseOptions(nil, tc.args)

			// Assert
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), "usage:") {
				t.Fatalf("expected usage message in error, got: %v", err)
			}
		})
	}
}

func TestConfigSetCmd_UseOptions_AcceptsTwoArgs(t *testing.T) {
	// Arrange
	storage := &MockKeyStorage{}
	cmd := NewConfigSetCmd(storage)

	// Act
	executor, err := cmd.UseOptions(nil, []string{"storage-mode", "local"})

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if executor == nil {
		t.Fatal("expected non-nil executor")
	}
}

func TestConfigSetCmd_Execute_SetsStorageModeLocal(t *testing.T) {
	// Arrange
	storage := &MockKeyStorage{storageMode: domain.InClusterStorageMode}
	cmd := NewConfigSetCmd(storage)
	cmd.options = &configSetOptions{Key: "storage-mode", Value: "local"}

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result != "storage-mode set to local" {
		t.Fatalf("expected 'storage-mode set to local', got %q", result)
	}
	if storage.storageMode != domain.LocalStorageMode {
		t.Fatalf("expected storage mode to be local, got %v", storage.storageMode)
	}
}

func TestConfigSetCmd_Execute_SetsStorageModeCluster(t *testing.T) {
	// Arrange
	storage := &MockKeyStorage{storageMode: domain.LocalStorageMode}
	cmd := NewConfigSetCmd(storage)
	cmd.options = &configSetOptions{Key: "storage-mode", Value: "cluster"}

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result != "storage-mode set to cluster" {
		t.Fatalf("expected 'storage-mode set to cluster', got %q", result)
	}
	if storage.storageMode != domain.InClusterStorageMode {
		t.Fatalf("expected storage mode to be cluster, got %v", storage.storageMode)
	}
}

func TestConfigSetCmd_Execute_AlreadySetReturnsMessage(t *testing.T) {
	// Arrange
	storage := &MockKeyStorage{storageMode: domain.LocalStorageMode}
	cmd := NewConfigSetCmd(storage)
	cmd.options = &configSetOptions{Key: "storage-mode", Value: "local"}

	// Act
	result, err := cmd.Execute()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	expected := "storage-mode is already set to local"
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestConfigSetCmd_Execute_InvalidStorageMode(t *testing.T) {
	// Arrange
	storage := &MockKeyStorage{storageMode: domain.LocalStorageMode}
	cmd := NewConfigSetCmd(storage)
	cmd.options = &configSetOptions{Key: "storage-mode", Value: "invalid"}

	// Act
	_, err := cmd.Execute()

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "invalid storage mode") {
		t.Fatalf("expected 'invalid storage mode' in error, got: %v", err)
	}
}

func TestConfigSetCmd_Execute_UnknownKey(t *testing.T) {
	// Arrange
	storage := &MockKeyStorage{}
	cmd := NewConfigSetCmd(storage)
	cmd.options = &configSetOptions{Key: "unknown-key", Value: "value"}

	// Act
	_, err := cmd.Execute()

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown configuration key") {
		t.Fatalf("expected 'unknown configuration key' in error, got: %v", err)
	}
}

func TestConfigSetCmd_Execute_GetStorageModeError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("get storage mode error")
	storage := &MockKeyStorage{storageModeErr: expectedErr}
	cmd := NewConfigSetCmd(storage)
	cmd.options = &configSetOptions{Key: "storage-mode", Value: "local"}

	// Act
	_, err := cmd.Execute()

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestConfigSetCmd_Execute_SetStorageModeError(t *testing.T) {
	// Arrange
	expectedErr := errors.New("set storage mode error")
	storage := &MockKeyStorage{
		storageMode: domain.InClusterStorageMode, // Different from what we're setting
		setModeErr:  expectedErr,
	}
	cmd := NewConfigSetCmd(storage)
	cmd.options = &configSetOptions{Key: "storage-mode", Value: "local"}

	// Act
	_, err := cmd.Execute()

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
