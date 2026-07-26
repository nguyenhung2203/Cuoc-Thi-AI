package storage

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// ErrInvalidKey is returned for storage keys that could escape the uploads
// directory or address hidden files.
var ErrInvalidKey = errors.New("storage: invalid key")

// LocalStore serves and writes files inside a single directory. Containment
// is enforced twice: a cheap key pre-check (clean 404s, defense in depth) and
// os.Root, which blocks path traversal AND symlink escapes at the syscall
// level on every platform, including Windows.
type LocalStore struct {
	dir  string
	root *os.Root
}

// NewLocalStore creates dir if needed and opens it as an os.Root jail.
func NewLocalStore(dir string) (*LocalStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("storage: create dir %q: %w", dir, err)
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("storage: open root %q: %w", dir, err)
	}
	return &LocalStore{dir: dir, root: root}, nil
}

// Dir returns the store's base directory (for startup logging).
func (s *LocalStore) Dir() string { return s.dir }

// validKey rejects keys that are empty, hidden, or contain path separators /
// traversal sequences. It deliberately tolerates everything legacy keys may
// contain (dots, spaces, unicode) so files uploaded before this change stay
// reachable.
func validKey(key string) bool {
	if key == "" || strings.HasPrefix(key, ".") {
		return false
	}
	if strings.ContainsAny(key, `/\`) || strings.Contains(key, "..") {
		return false
	}
	return true
}

// Open opens key for reading. Returns ErrInvalidKey for malicious keys and
// fs.ErrNotExist-wrapping errors for missing files.
func (s *LocalStore) Open(key string) (*os.File, error) {
	if !validKey(key) {
		return nil, ErrInvalidKey
	}
	return s.root.Open(key)
}

// Create creates/truncates key for writing.
func (s *LocalStore) Create(key string) (*os.File, error) {
	if !validKey(key) {
		return nil, ErrInvalidKey
	}
	return s.root.Create(key)
}

// Rename renames oldKey to newKey inside the jail (atomic publish of .part files).
func (s *LocalStore) Rename(oldKey, newKey string) error {
	if !validKey(oldKey) || !validKey(newKey) {
		return ErrInvalidKey
	}
	return s.root.Rename(oldKey, newKey)
}

// Remove deletes key. Best-effort callers may ignore the error.
func (s *LocalStore) Remove(key string) error {
	if !validKey(key) {
		return ErrInvalidKey
	}
	return s.root.Remove(key)
}
