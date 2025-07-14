package shared

import (
	"crypto/sha1"
	"encoding/hex"
	"path/filepath"
	"strings"
)

// Note: use number based id instead of file path for file tracking
func HashPath(path string) string {
	h := sha1.New()
	h.Write([]byte(path))
	return hex.EncodeToString(h.Sum(nil))
}

func GetRelativePath(path string) string {
	relativePath := strings.TrimPrefix(path, ActiveStorage.Path)
	return strings.TrimPrefix(relativePath, string(filepath.Separator)) // remove leading "/"
}