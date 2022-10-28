package atomicfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAbort(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test.txt")
	f, err := New(path, 0644)
	require.NoError(t, err, "new")
	_, err = f.Write([]byte("foo"))
	require.NoError(t, err, "write")
	err = f.Abort()
	require.NoError(t, err, "abort")
	require.NoFileExists(t, path, "real")
	require.NoFileExists(t, f.File.Name(), "temp")
}

func TestClose(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test.txt")
	defer func() { _ = os.Remove(path) }()
	f, err := New(path, 0644)
	require.NoError(t, err, "new")
	_, err = f.Write([]byte("foo"))
	require.NoError(t, err, "write")
	err = f.Close()
	require.NoError(t, err, "close")
	require.FileExists(t, path, "real")
	require.NoFileExists(t, f.File.Name(), "temp")
	content, err := os.ReadFile(path)
	require.NoError(t, err, "read")
	require.Equal(t, "foo", string(content), "content")
}

func TestReplace(t *testing.T) {
	src, err := os.CreateTemp("", "")
	defer func() { _ = os.Remove(src.Name()) }()
	require.NoError(t, err, "create src")
	_, err = src.Write([]byte("I'm the source"))
	require.NoError(t, err, "write src")
	err = src.Close()
	require.NoError(t, err, "close src")

	dest, err := os.CreateTemp("", "")
	defer func() { _ = os.Remove(dest.Name()) }()
	require.NoError(t, err, "create dest")
	_, err = dest.Write([]byte("I'm the destination"))
	require.NoError(t, err, "write dest")
	err = dest.Close()
	require.NoError(t, err, "close dest")

	err = ReplaceFile(src.Name(), dest.Name())
	require.NoError(t, err, "replace")
	require.FileExists(t, dest.Name(), "dest")
	require.NoFileExists(t, src.Name(), "src")
	content, err := os.ReadFile(dest.Name())
	require.NoError(t, err, "read dest")
	require.Equal(t, "I'm the source", string(content))
}

func TestWriteFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err, "create temp dir")
	defer func() { _ = os.RemoveAll(dir) }()
	existingDest := filepath.Join(dir, "dest.txt")
	err = os.WriteFile(existingDest, []byte("existing"), 0644)
	require.NoError(t, err, "write dest")
	err = WriteFile(existingDest, []byte("new content"))
	require.NoError(t, err, "WriteFile")
	content, err := os.ReadFile(existingDest)
	require.NoError(t, err, "read")
	require.Equal(t, "new content", string(content), "content")
	entries, err := os.ReadDir(dir)
	require.NoError(t, err, "read dir")
	require.Equal(t, 1, len(entries))
	require.Equal(t, "dest.txt", filepath.Base(entries[0].Name()))
}
