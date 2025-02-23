package serialise

import "testing"

func TestToBytesI64(t *testing.T) {

	tests := []int64{
		42, -42, 0, 2165127351287351239, -2165127351287351239,
	}

	for _, test := range tests {

		b, err := ToBytesI64(test)
		if err != nil {
			t.Fatalf("Unexpected error when serialising int64: %v", test)
		}

		v, err := FromBytesI64(b)
		if err != nil {
			t.Fatalf("Unexpected error when deserialising int64: %v", test)
		}

		if v != test {
			t.Fatalf("Mismatch in values: expected: %d, got: %d", test, v)
		}
	}
}

func TestSizeOfI64(t *testing.T) {

	b, _ := ToBytesI64(int64(0))

	if int64(len(b)) != SizeOfI64() {
		t.Fatal("SizeOfI64 returning incorrect size!")
	}
}
