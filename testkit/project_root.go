package testkit

import (
	"os"
	"path/filepath"
)

func GetProjectRoot() string {
	searchLevel := []string{".", "..", "../.."}

	for _, l := range searchLevel {
		fs, err := os.Stat(filepath.Join(l, "go.mod"))
		if os.IsNotExist(err) || fs.IsDir() {
			continue
		}

		return l
	}

	return "."
}
