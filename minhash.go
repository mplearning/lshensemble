package lshensemble

import (
	"encoding/binary"
	"hash/fnv"
	"io"
	"math/rand"

	minwise "github.com/dgryski/go-minhash"
)

// The number of byte in a hash value for Minhash
const HashValueSize = 8

// Represents a MinHash object
type Minhash struct {
	mw *minwise.MinWise
}

// Represents a MinHash signature - an array of hash values
type Signature []uint64

// Initialize a MinHash object with a seed and the number of
// hash functions.
func NewMinhash(seed, numHash int) *Minhash {
	r := rand.New(rand.NewSource(int64(seed)))
	b := binary.BigEndian
	b1 := make([]byte, HashValueSize)
	b2 := make([]byte, HashValueSize)
	b.PutUint64(b1, uint64(r.Int63()))
	b.PutUint64(b2, uint64(r.Int63()))
	fnv1 := fnv.New64a()
	fnv2 := fnv.New64a()
	h1 := func(b []byte) uint64 {
		fnv1.Reset()
		fnv1.Write(b1)
		fnv1.Write(b)
		return fnv1.Sum64()
	}
	h2 := func(b []byte) uint64 {
		fnv2.Reset()
		fnv2.Write(b2)
		fnv2.Write(b)
		return fnv2.Sum64()
	}
	return &Minhash{minwise.NewMinWise(h1, h2, numHash)}
}

// Push a new value to the MinHash object.
// The value should be serialized to byte slice.
func (m *Minhash) Push(b []byte) {
	m.mw.Push(b)
}

// Export the MinHash signature.
func (m *Minhash) Signature() Signature {
	return m.mw.Signature()
}

// Serialize the siganture into the writer.
func (sig Signature) Write(w io.Writer) error {
	for i := range sig {
		if err := binary.Write(w, binary.BigEndian, sig[i]); err != nil {
			return err
		}
	}
	return nil
}

// Deserialize the signature from the reader.
func (sig Signature) Read(r io.Reader) error {
	for i := range sig {
		if err := binary.Read(r, binary.BigEndian, &(sig[i])); err != nil {
			return err
		}
	}
	return nil
}

// Compute the serialized length in terms of number of bytes.
func (sig Signature) ByteLen() int {
	return len(sig) * HashValueSize
}
