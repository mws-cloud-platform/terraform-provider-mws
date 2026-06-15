package examples

import (
	"os"
	"strings"
)

func listExampleResources(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "mws_") {
			continue
		}

		names = append(names, entry.Name())
	}

	return names, nil
}
