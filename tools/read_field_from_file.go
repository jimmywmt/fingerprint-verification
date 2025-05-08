package tools

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func readFieldFromFile(path, key string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		if strings.HasPrefix(line, key) {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
		if err == io.EOF {
			break
		}
	}
	return "", fmt.Errorf("key '%s' not found in %s", key, path)
}
