package homedir

import (
	"os"
	"path/filepath"
	"strings"
)

// Expand expands a leading "~/" in path to the current user's home directory.
// It returns path unchanged if it does not start with "~/" or the home
// directory cannot be determined.
func Expand(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
