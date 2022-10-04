package util

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
