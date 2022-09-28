package security

import "crypto/sha256"

type Hasher interface {
	CalculateHash(hashTarget string) (string, error)
}

var HashMaker Hasher

func init() {
	HashMaker = hasherImpl{}
}

type hasherImpl struct{}

func (c hasherImpl) CalculateHash(hashTarget string) (string, error) {
	sum256 := sha256.Sum256([]byte(hashTarget))
	return string(sum256[:]), nil
}
