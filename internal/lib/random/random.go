package random

import "math/rand/v2"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateAlias(AliasLength int) string {
	res := make([]rune, AliasLength)
	for i := range res {
		res[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	
	return string(res)
}
