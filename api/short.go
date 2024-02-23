package api

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type shortBody struct {
	LongURL string `json:"long_url"`
}

const baseStr = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYSabcdefghijklmnopqrstuvwxyz"

var (
	check = regexp.MustCompile(baseStr)
	base  = []rune(baseStr)
)

// encode takes an id from the database
// and encode it to a base62 string
// we assume an higher bound of 7 characters max which is 62^7 = 3.52e+12 ids
func encode(id int64) (string, error) {
	if id > int64(pow(62, 7)) {
		return "", errors.New("id overflow")
	}
	res := id
	encoded := make([]rune, 7)
	i := 6
	for res > 0 {
		encoded[i] = base[res%62]
		res /= 62
		i--
	}
	i++
	var b strings.Builder
	for i < 7 {
		_, err := b.WriteRune(encoded[i])
		if err != nil {
			return "", err
		}
		i++
	}
	return b.String(), nil
}

// decode does the opposite of encode, it returns the base10 representation of the input string encoded in base62
func decode(slug string) (int64, error) {
	if len(slug) > 7 {
		return 0, errors.New("overflow: slug must be <= 7 character")
	}
	var id int64
	j := 0
	for i := len(slug) - 1; i >= 0; i-- {
		for n := 0; n < len(base); n++ {
			if base[n] == rune(slug[i]) {
				id += int64(n * pow(62, j))
				break
			}
			// if the character was not found in our base62 character list
			if n == len(base) {
				return id, fmt.Errorf("invalid character %s in the slug", string(slug[i]))
			}
		}
		j++
	}
	return id, nil
}

func pow(n, e int) int {
	if e == 0 {
		return 1
	}
	r := n
	for i := 1; i < e; i++ {
		r *= n
	}
	return r
}
