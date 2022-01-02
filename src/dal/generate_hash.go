package dal

import (
	"crypto/rand"
	"encoding/base64"
)

const _BITES_LENGTH int8 = 4
const _HARDCODED_LENGTH int8 = 6

func generateHash() (string, error) {
	buf := make([]byte, _BITES_LENGTH)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(buf)
	asRunes := []rune(str)
	base64url := string(asRunes[0:_HARDCODED_LENGTH])
	return base64url, nil
}
