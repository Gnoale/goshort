package api

import (
	"strings"
)

type shortBody struct {
	LongURL string `json:"long_url"`
}

/*
	Shortener logic

	We receive https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74

	1- convert this with a hash function so each url.String() is mapped with a short hash value

	2- store the mapping in the database

	3- return the shorten URL


	Redirect logic

	We receive https://<my-domain>/<slug>

	1- inspect the database for such a slug value

		if found, return the associated url value and send an http 301 response with the url in location

		if not found return 404



*/

var base = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYSabcdefghijklmnopqrstuvwxyz")

// encode takes an id from the database
// and encode it to a base62 string
// we assume an higher bound of 7 characters max which is 62^7 = 3.52e+12 ids
func encode(id int64) string {
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
			panic(err)
		}
		i++
	}
	return b.String()
}

func decode(slug string) int64 {
	var id int64
	j := 0
	for i := len(slug) - 1; i >= 0; i-- {
		for n := 0; n < len(base); n++ {
			if base[n] == rune(slug[i]) {
				id += int64(n * pow(62, j))
				break
			}
		}
		j++
	}
	return id
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
