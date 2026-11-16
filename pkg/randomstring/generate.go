package randomstring

import (
	"crypto/rand"
	"math/big"
)

const letters = "0123456789abcdefghijklmnopqrstuvwxyz"

// Generate returns a cryptographically random base36 string of the given
// length. It panics if the system random source fails, which is unrecoverable.
func Generate(length int) string {
	n := make([]byte, length)
	max := big.NewInt(int64(len(letters)))
	for i := range n {
		c, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		n[i] = letters[c.Int64()]
	}
	return string(n)
}
