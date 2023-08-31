package crypto

import "golang.org/x/crypto/bcrypt"

// Encrypt encrypts a string
func Encrypt(source string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(source), bcrypt.DefaultCost)
	return string(hashed), err
}

// Compare compares clear and hashed texts
func Compare(hashed, clear string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(clear))
}
