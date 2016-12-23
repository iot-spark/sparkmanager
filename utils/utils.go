package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	bStr := make([]byte, n)
	_, err := rand.Read(bStr)
	if err != nil {
		return nil, err
	}

	return bStr, nil
}

func GeneratePskKey(s int) (string, error) {
	data, err := GenerateRandomBytes(s)

	hash := md5.New()
	hash.Write([]byte(data))

	return hex.EncodeToString(hash.Sum(nil)), err
}

func GeneratePskIdentity(name string) string {
	psk := strings.Replace(strings.Trim(name, " "), " ", "_", -1)

	return strings.ToLower(psk)
}
