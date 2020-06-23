package utils

import (
	"bufio"
	"os"
	"strings"
)

func Sum(input []int) int {
	sum := 0
	for i := range input {
		sum += input[i]
	}
	return sum
}

func IntIn(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringIn(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringInMap(a string, mapItem map[string]interface{}) bool {
	for _, b := range mapItem {
		if b == a {
			return true
		}
	}
	return false
}

func StringInLike(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}

func StringInLower(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1    //not found.
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func Btoi(val bool) int {
	if val{
		return 1
	}else{
		return 0
	}
}