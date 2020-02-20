package util

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
)

// GetKeyHash can be used to create a hash from a string
func GetKeyHash(key string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SplitArgs(s string) ([]string, error) {
	indexes := regexp.MustCompile("\"").FindAllStringIndex(s, -1)
	if len(indexes)%2 != 0 {
		return []string{}, errors.New("quotes not closed")
	}

	indexes = append([][]int{{0, 0}}, indexes...)

	if indexes[len(indexes)-1][1] < len(s) {
		indexes = append(indexes, [][]int{{len(s), 0}}...)
	}

	var args []string
	for i := 0; i < len(indexes)-1; i++ {
		start := indexes[i][1]
		end := Min(len(s), indexes[i+1][0])

		if i%2 == 0 {
			args = append(args, strings.Split(strings.Trim(s[start:end], " "), " ")...)
		} else {
			args = append(args, s[start:end])
		}

	}

	cleanedArgs := make([]string, len(args))
	count := 0

	for _, arg := range args {
		if arg != "" {
			cleanedArgs[count] = strings.TrimSpace(arg)
			count++
		}
	}

	return cleanedArgs[0:count], nil
}

func Deduplicate(a []string) []string {
	check := make(map[string]int)
	result := make([]string, 0)
	for _, val := range a {
		check[val] = 1
	}

	for key := range check {
		result = append(result, key)
	}

	return result
}
