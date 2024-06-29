package utils

import (
	"fmt"
	"log"

	hashers "github.com/meehow/go-django-hashers"
)

const DefaultHashingIteration = 1500

func GetEncodedPassword(plainPass string) string {
	hashers.Iter = DefaultHashingIteration
	hashPass, err := hashers.MakePassword(plainPass)
	if err != nil {
		return ""
	}
	fmt.Println()

	return hashPass
}

func VerifyPassword(plainPass, encodedPass string) bool {
	hashers.Iter = DefaultHashingIteration
	bl, err := hashers.CheckPassword(plainPass, encodedPass)
	if err != nil {
		log.Println(err)
		return false
	}

	return bl
}
