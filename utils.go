package youcrawl

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
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

func RandomIntRangeWithStringSeed(min int, max int,seedString string) int {
	h := md5.New()
	io.WriteString(h, seedString)
	var seed uint64 = binary.BigEndian.Uint64(h.Sum(nil))
	fmt.Println(seed)
	rand.Seed(int64(seed))
	return rand.Intn(max - min) + min
}