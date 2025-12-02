package util

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomReceiveCode() (int32, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(90000000))
	if err != nil {
		return 0, err
	}
	return int32(n.Int64() + 10000000), err
}
