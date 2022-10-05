package security

import "crypto/sha256"

type Hasher interface {
	SHA256(hashTarget string) (string, error)
}

var HashMaker Hasher

func init() {
	HashMaker = hasherImpl{}
}

type hasherImpl struct{}

func (c hasherImpl) SHA256(hashTarget string) (string, error) {
	sum256 := sha256.Sum256([]byte(hashTarget))
	return string(sum256[:]), nil
}
