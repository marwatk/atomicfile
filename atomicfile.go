package atomicfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var _ io.WriteCloser = AtomicFile{}

type AtomicFile struct {
	finalPath string
	*os.File
}

func New(path string, mode os.FileMode) (AtomicFile, error) {
	f := AtomicFile{finalPath: path}
	var err error
	f.File, err = os.CreateTemp(filepath.Dir(f.finalPath), filepath.Base(f.finalPath))
	if err != nil {
		return AtomicFile{}, fmt.Errorf("creating temp file for atomicfile: %w", err)
	}
	return f, nil
}

func (f AtomicFile) Abort() error {
	err1 := f.File.Close()
	err2 := os.Remove(f.File.Name())
	if err2 != nil {
		return err2
	}
	if err1 != nil {
		return err1
	}
	return nil
}

// Close closes the atomic file renaming it to it's final name or
// returning an error and deleting the temp file.
func (f AtomicFile) Close() error {
	err := f.File.Close()
	if err != nil {
		_ = f.Abort()
		return fmt.Errorf("closing temp file for atomicfile: %w", err)
	}
	err = ReplaceFile(f.File.Name(), f.finalPath)
	if err != nil {
		_ = f.Abort()
	}
	return err
}

// ReplaceFile atomically replaces the complete contents of destPath
// with the complete contents of srcPath preserving file attributes.
// Files must be on the same filesystem on most (all) systems. On failure
// neither file is modified.
func ReplaceFile(srcPath, destPath string) error {
	return replaceFile(srcPath, destPath)
}
