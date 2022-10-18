package util

import "math/rand"

/*
	Helper function
	reference: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
*/
func RandStringBytes(letterBytes string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateLicense(username string) string {
	license := username
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	randomString := RandStringBytes(characters, 7)
	return license + randomString
}

func GenerateSecretCode() (secretCode []string) {
	uppercaseCharacters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 9; i++ {
		secretCode = append(secretCode, RandStringBytes(uppercaseCharacters, 5))
	}
	return secretCode
}
