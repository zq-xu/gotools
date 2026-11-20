package format

import (
	"bufio"
	"os"
	"strings"
)

var GoProject string

// InitGoProject reads go.mod and initializes the GoProject variable.
func InitGoProject() error {
	file, err := os.Open("go.mod")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			GoProject = strings.TrimSpace(strings.TrimPrefix(line, "module"))
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
