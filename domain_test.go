package lshensemble

import (
	"bytes"
	"math/rand"
	"testing"
)

func Test_Serialization(t *testing.T) {
	key := "abcd"
	for i := 0; i < 500; i++ {
		size := rand.Int()
		sig := randomSignature(128, int64(42))
		record := &DomainRecord{
			Key:       key,
			Size:      size,
			Signature: sig,
		}
		var data bytes.Buffer
		n, err := record.Write(&data, func(s string) ([]byte, error) {
			return []byte(s), nil
		})
		if err != nil {
			t.Error(err)
		}
		t.Log("Number of bytes written to buffer", n)
		t.Log(record)

		var record2 DomainRecord
		n, err = record2.Read(&data, 4, func(b []byte) (string, error) {
			return string(b), nil
		})
		if err != nil {
			t.Error(err)
		}
		t.Log("Number of bytes read", n)
		t.Log(record2)
		if record2.Key != record.Key || record2.Size != record.Size {
			t.Error("Incorrect record after deserialization")
		}
		for i := range record.Signature {
			if record.Signature[i] != record2.Signature[i] {
				t.Error("Signature does not match after serialization")
			}
		}
	}
}
