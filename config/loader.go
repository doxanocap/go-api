package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if line == "" {
			continue
		}

		vars := strings.Split(line, "=")
		if len(vars) < 2 {
			return fmt.Errorf("scanning line %d: invalid format", i)
		}
		key := vars[0]
		val := removeBrackets(line[len(key)+1:])

		temp := os.Getenv(key)
		if temp != "" {
			continue
		}

		err := os.Setenv(key, val)
		if err != nil {
			return fmt.Errorf("scanning line %d: %v", i, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}

	return nil
}

func removeBrackets(str string) string {
	if len(str) < 3 {
		return str
	}

	if str[0] == '"' {
		str = str[1:]

	}
	ln := len(str)
	if str[ln-1] == '"' {
		str = str[:ln-1]
	}
	return str
}
