package decoder

import (
	"encoding/base64"
	"fmt"
	"sopsctl/pkg/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestBase64Decoder_EditDecodedFile_Success(t *testing.T) {
	decoder := Base64Decoder{}

	// Create a valid Kubernetes secret with a single data key
	// Note: In YAML, the data values are plain text, yaml.Unmarshal converts them to []byte
	const yamlContent = "TestContent"
	secretYAML := `apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  config.yaml: ` + base64.StdEncoding.EncodeToString([]byte(yamlContent)) + `
type: Opaque`

	// Call EditDecodedFile
	decoded, restoreFunc, err := decoder.EditDecodedFile([]byte(secretYAML), "config.yaml")

	// Assert no error
	require.NoError(t, err)
	assert.NotNil(t, restoreFunc)

	// Assert the decoded content is correct (should be base64 decoded)
	assert.Equal(t, yamlContent, string(decoded))
}

func TestBase64Decoder_EditDecodedFile_NoDataKeys(t *testing.T) {
	decoder := Base64Decoder{}

	// Secret with no data keys
	secretYAML := `apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data: {}
type: Opaque`

	// Call EditDecodedFile
	decoded, restoreFunc, err := decoder.EditDecodedFile([]byte(secretYAML), "config.yaml")

	// Assert error occurred
	assert.Error(t, err)
	assert.Nil(t, decoded)
	assert.Nil(t, restoreFunc)
	assert.Contains(t, err.Error(), "did not find data for key")
}

func TestBase64Decoder_RestoreEncodedFile_Success(t *testing.T) {
	decoder := Base64Decoder{}

	// Create a valid secret with properly padded base64
	secretYAML := `apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  config.yaml: b3JpZ2luYWw6IGNvbnRlbnQ=
type: Opaque`

	// Call EditDecodedFile to get the restore function
	_, restoreFunc, err := decoder.EditDecodedFile([]byte(secretYAML), "config.yaml")
	require.NoError(t, err)
	require.NotNil(t, restoreFunc)

	// Modify the content and restore it
	modifiedContent := []byte("modified: content\nnewkey: newvalue")
	restored, err := restoreFunc(modifiedContent)

	// Assert no error
	require.NoError(t, err)
	assert.NotNil(t, restored)

	// Parse the restored YAML and verify the structure
	var restoredSecret = &domain.RawSecret{}
	err = yaml.Unmarshal(restored, &restoredSecret)
	require.NoError(t, err)

	// Verify the secret has the expected structure
	assert.Equal(t, "v1", restoredSecret.APIVersion)
	assert.Equal(t, "Secret", restoredSecret.Kind)

	// Verify the data is base64 encoded
	data := restoredSecret.Data
	encodedValue := data["config.yaml"]

	// Decode and verify the content
	decodedValue, err := base64.StdEncoding.DecodeString(encodedValue)
	require.NoError(t, err)
	assert.Equal(t, string(modifiedContent), string(decodedValue))
}

func TestBase64Decoder_RestoreEncodedFile_EmptyContent(t *testing.T) {
	decoder := Base64Decoder{}

	// Create a valid secret with properly padded base64
	secretYAML := `apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  config.yaml: b3JpZ2luYWw=
type: Opaque`

	// Call EditDecodedFile to get the restore function
	_, restoreFunc, err := decoder.EditDecodedFile([]byte(secretYAML), "config.yaml")
	require.NoError(t, err)

	// Restore with empty content
	emptyContent := []byte("")
	restored, err := restoreFunc(emptyContent)

	// Assert no error
	require.NoError(t, err)
	assert.NotNil(t, restored)

	// Verify the restored secret contains empty base64 content
	var restoredSecret map[string]interface{}
	err = yaml.Unmarshal(restored, &restoredSecret)
	require.NoError(t, err)

	data := restoredSecret["data"].(map[string]interface{})
	encodedValue := data["config.yaml"].(string)

	decodedValue, err := base64.StdEncoding.DecodeString(encodedValue)
	require.NoError(t, err)
	assert.Equal(t, "", string(decodedValue))
}

func TestBase64Decoder_RestoreEncodedFile_LargeContent(t *testing.T) {
	decoder := Base64Decoder{}

	// Create a valid secret with properly padded base64
	secretYAML := `apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  config.yaml: c21hbGw=
type: Opaque`

	// Call EditDecodedFile to get the restore function
	_, restoreFunc, err := decoder.EditDecodedFile([]byte(secretYAML), "config.yaml")
	require.NoError(t, err)

	// Create large content
	largeContent := make([]byte, 10000)
	for i := range largeContent {
		largeContent[i] = byte('a' + (i % 26))
	}

	// Restore with large content
	restored, err := restoreFunc(largeContent)

	// Assert no error
	require.NoError(t, err)
	assert.NotNil(t, restored)

	// Verify the content is correctly restored
	var restoredSecret map[string]interface{}
	err = yaml.Unmarshal(restored, &restoredSecret)
	require.NoError(t, err)

	data := restoredSecret["data"].(map[string]interface{})
	encodedValue := data["config.yaml"].(string)

	decodedValue, err := base64.StdEncoding.DecodeString(encodedValue)
	require.NoError(t, err)
	assert.Equal(t, largeContent, decodedValue)
}

func TestBase64Decoder_RoundTrip(t *testing.T) {
	decoder := Base64Decoder{}

	originalContent := "key: value\nfoo: bar\nnested:\n  item: test"
	// Properly padded base64
	secretYAML := `apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  config.yaml: a2V5OiB2YWx1ZQpmb286IGJhcgpuZXN0ZWQ6CiAgaXRlbTogdGVzdA==
type: Opaque`

	// Decode
	decoded, restoreFunc, err := decoder.EditDecodedFile([]byte(secretYAML), "config.yaml")
	require.NoError(t, err)
	assert.Equal(t, originalContent, string(decoded))

	// Modify
	modifiedContent := decoded
	modifiedContent = append(modifiedContent, []byte("\nadded: field")...)

	// Restore
	restored, err := restoreFunc(modifiedContent)
	require.NoError(t, err)

	// Decode again to verify
	decoded2, _, err := decoder.EditDecodedFile(restored, "config.yaml")
	require.NoError(t, err)
	assert.Equal(t, string(modifiedContent), string(decoded2))
}

// TestBase64Decoder_PreservesContentExactly verifies that content is preserved byte-for-byte
// through the decode -> edit -> encode cycle, including trailing newlines.
func TestBase64Decoder_PreservesContentExactly(t *testing.T) {
	decoder := Base64Decoder{}

	testCases := []struct {
		name    string
		content string
	}{
		{
			name:    "content without trailing newline",
			content: "key: value\nfoo: bar",
		},
		{
			name:    "content with trailing newline",
			content: "key: value\nfoo: bar\n",
		},
		{
			name:    "content with multiple trailing newlines",
			content: "key: value\nfoo: bar\n\n",
		},
		{
			name:    "single line without newline",
			content: "simple-value",
		},
		{
			name:    "single line with newline",
			content: "simple-value\n",
		},
		{
			name:    "content with leading and trailing whitespace",
			content: "  leading\ntrailing  \n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a secret with the content base64-encoded using StdEncoding (Kubernetes standard)
			encodedContent := base64.StdEncoding.EncodeToString([]byte(tc.content))
			secretYAML := fmt.Sprintf(`apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  config.yaml: %s
type: Opaque`, encodedContent)

			// Decode the secret
			decoded, restoreFunc, err := decoder.EditDecodedFile([]byte(secretYAML), "config.yaml")
			require.NoError(t, err, "EditDecodedFile should not fail")

			// Verify decoded content matches exactly
			assert.Equal(t, tc.content, string(decoded), "Decoded content should match original exactly")

			// Restore (simulate saving without changes)
			restored, err := restoreFunc(decoded)
			require.NoError(t, err, "restoreFunc should not fail")

			// Decode again to verify round-trip
			decoded2, _, err := decoder.EditDecodedFile(restored, "config.yaml")
			require.NoError(t, err, "Second EditDecodedFile should not fail")

			// Verify content is preserved exactly through the round-trip
			assert.Equal(t, tc.content, string(decoded2), "Content should be preserved exactly through round-trip")
		})
	}
}

// TestBase64Decoder_EditorAddsTrailingNewline simulates a user editing content
// where the editor adds a trailing newline (like nano does by default)
func TestBase64Decoder_EditorAddsTrailingNewline(t *testing.T) {
	decoder := Base64Decoder{}

	// Original content has no trailing newline
	originalContent := "database-password"
	encodedContent := base64.StdEncoding.EncodeToString([]byte(originalContent))
	secretYAML := fmt.Sprintf(`apiVersion: v1
kind: Secret
metadata:
  name: test-secret
  namespace: default
data:
  password: %s
type: Opaque`, encodedContent)

	// Decode
	decoded, restoreFunc, err := decoder.EditDecodedFile([]byte(secretYAML), "password")
	require.NoError(t, err)
	assert.Equal(t, originalContent, string(decoded))

	// Simulate editor adding trailing newline (like nano does)
	editedContent := []byte(originalContent + "\n")

	// Restore with the editor-modified content
	restored, err := restoreFunc(editedContent)
	require.NoError(t, err)

	// Decode again - should now have the trailing newline
	decoded2, _, err := decoder.EditDecodedFile(restored, "password")
	require.NoError(t, err)

	// The saved content should include the newline the user/editor added
	assert.Equal(t, string(editedContent), string(decoded2), "Editor-added newline should be preserved")
}
