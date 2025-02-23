package serialise

import (
	"fmt"
)

// ToBytesI64 ensures a standard treatment of int64 serialisation
func ToBytesI64(v int64) ([]byte, error) {
	b, _, err := ToBytes(v, WithSerialisationApproach(NewMinDataApproachWithVersion(V1)))
	if err != nil {
		return nil, err
	}
	return b, nil
}

// SizeOfI64 returns standard length of int64 when serialised
func SizeOfI64() int64 {
	b, err := ToBytesI64(0)
	if err != nil {
		panic(fmt.Sprintf("unexpected error when determining serialised size of int64: %v", err))
	}
	return int64(len(b))
}

// FromBytesI64 deserialises an int64 from material created by ToBytesI64
func FromBytesI64(data []byte) (int64, error) {
	v, err := FromBytes(data, NewMinDataApproachWithVersion(V1))
	if err != nil {
		return 0, err
	}
	if i, ok := v.(int64); ok {
		return i, nil
	}
	panic("Should never have an issue deserialising int64")
}
