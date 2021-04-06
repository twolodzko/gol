package parser

import (
	"bufio"
	"os"
	"strings"
)

func ReadFile(name string) (string, error) {
	var lines []string

	file, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
}
