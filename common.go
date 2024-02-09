package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func padNumberWithZero(value int, expectedLength int) string {
	padded := fmt.Sprintf("%02d", value)
	valLength := countDigits(value)
	if valLength != expectedLength {
		return fmt.Sprintf("%s%d", strings.Repeat("0", expectedLength-valLength), value)
	}
	return padded
}

func checkCacheFolder() {
	if _, err := os.Stat(".cache"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(".cache", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func countDigits(i int) int {
	count := 0
	if i == 0 {
		count = 1
	}
	for i != 0 {
		i /= 10
		count = count + 1
	}
	return count
}

func contains(s []string, e string) bool {

	for _, ss := range s {
		if e == ss {
			return true
		}
	}
	return false
}
