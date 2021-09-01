package shortener

import (
	"math/rand"
	"strconv"
	"time"
)

func Shorten() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	lenLetters := len(letters)
	lenLink := 10

	var isLetter bool
	var link string

	for i := 0; i < lenLink; i++ {
		if r.Intn(2) != 0 {
			isLetter = true
		} else {
			isLetter = false
		}

		if isLetter {
			index := r.Intn(lenLetters)
			link = link + string(letters[index])
		} else {
			link = link + strconv.Itoa(r.Intn(10))
		}
	}

	return link
}
