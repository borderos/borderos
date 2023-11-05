package sysctl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const basePath = "/proc/sys"

// Read returns the specified sysctl value.
func Read(key string) (string, error) {
	path := strings.ReplaceAll(key, ".", "/")
	b, err := os.ReadFile(basePath + path)
	if err != nil {
		return "", fmt.Errorf("sysctl read: %w", err)
	}
	return string(b), nil
}

// Write writes value to the specified sysctl key.
func Write(key, value string) error {
	path := strings.ReplaceAll(key, ".", "/")
	return os.WriteFile(filepath.Join(basePath, path), []byte(value), 0644)
}
