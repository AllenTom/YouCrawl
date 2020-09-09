package youcrawl

import (
	"bufio"
	"os"
)

func ReadListFile(listFilePath string) ([]string, error) {
	file, err := os.Open(listFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, err
}
