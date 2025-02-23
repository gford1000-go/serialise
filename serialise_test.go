package serialise

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func timeNeq(x any, y any) (time.Time, bool, bool) {
	switch v := x.(type) {
	case time.Time:
		return v, v != (y.(time.Time)).Truncate(0), false
	default:
		return *new(time.Time), false, true
	}
}

func testCompareValue[T comparable](a, b any, name string, t *testing.T, opts ...func(any, any) (T, bool, bool)) {

	neq := func(x any, y any) (T, bool, bool) {
		if xx, ok := x.(T); ok {
			return xx, xx != y.(T), false
		}
		return *new(T), false, true
	}

	if len(opts) > 0 {
		neq = opts[0]
	}

	switch v := b.(type) {
	case T:
		if aa, test, bad := neq(a, v); test {
			t.Fatalf("Data mismatch: expected %v, got: %v", v, aa)
		} else if bad {
			t.Fatalf("Type mismatch: expected: %s, got: %s", name, fmt.Sprintf("%T", a))
		}
	default:
		t.Fatalf("Unexpected error: b was the wrong type: %s", fmt.Sprintf("%T", b))
	}
}

func testCompareSliceValue[T comparable](a, b any, name string, t *testing.T) {
	switch v := b.(type) {
	case []T:
		if aa, ok := a.([]T); ok {
			if len(aa) != len(v) {
				t.Fatalf("Data size mismatch: expected %v, got: %v", len(v), len(aa))
			}
			for i, vv := range aa {
				if v[i] != vv {
					t.Fatalf("Data mismatch at %d: expected %v, got: %v", i, vv, aa)
				}
			}
		} else {
			t.Fatalf("Type mismatch: expected: %s, got: %s", name, fmt.Sprintf("%T", a))
		}
	default:
		t.Fatalf("Unexpected error: b was the wrong type: %s", fmt.Sprintf("%T", b))
	}
}

func testComparePtrValue[T comparable](a, b any, name string, t *testing.T) {

	switch v := b.(type) {
	case *T:
		if aa, ok := a.(*T); ok {
			if (aa == nil && v != nil) || (aa != nil && v == nil) {
				t.Fatalf("Pointer mismatch: expected %v, got: %v", v, aa)
			}
			if *aa != *v {
				t.Fatalf("Data mismatch: expected %v, got: %v", *v, *aa)
			}
		} else {
			t.Fatalf("Type mismatch: expected: %s, got: %s", name, fmt.Sprintf("%T", a))
		}
	default:
		t.Fatalf("Unexpected error: b was the wrong type: %s", fmt.Sprintf("%T", b))
	}
}

func compareValue(a, b any, name string, t *testing.T) {
	if b == nil {
		if a != nil {
			t.Fatalf("Mismatch in <nil>")
		}
		return
	}

	switch b.(type) {
	case []byte:
		testCompareSliceValue[byte](a, b, name, t)
	case int8:
		testCompareValue[int8](a, b, name, t)
	case *int8:
		testComparePtrValue[int8](a, b, name, t)
	case []int8:
		testCompareSliceValue[int8](a, b, name, t)
	case int16:
		testCompareValue[int16](a, b, name, t)
	case *int16:
		testComparePtrValue[int16](a, b, name, t)
	case []int16:
		testCompareSliceValue[int16](a, b, name, t)
	case int32:
		testCompareValue[int32](a, b, name, t)
	case *int32:
		testComparePtrValue[int32](a, b, name, t)
	case []int32:
		testCompareSliceValue[int32](a, b, name, t)
	case int64:
		testCompareValue[int64](a, b, name, t)
	case *int64:
		testComparePtrValue[int64](a, b, name, t)
	case []int64:
		testCompareSliceValue[int64](a, b, name, t)
	case uint8:
		testCompareValue[uint8](a, b, name, t)
	case *uint8:
		testComparePtrValue[uint8](a, b, name, t)
	case uint16:
		testCompareValue[uint16](a, b, name, t)
	case *uint16:
		testComparePtrValue[uint16](a, b, name, t)
	case []uint16:
		testCompareSliceValue[uint16](a, b, name, t)
	case uint32:
		testCompareValue[uint32](a, b, name, t)
	case *uint32:
		testComparePtrValue[uint32](a, b, name, t)
	case []uint32:
		testCompareSliceValue[uint32](a, b, name, t)
	case uint64:
		testCompareValue[uint64](a, b, name, t)
	case *uint64:
		testComparePtrValue[uint64](a, b, name, t)
	case []uint64:
		testCompareSliceValue[uint64](a, b, name, t)
	case float32:
		testCompareValue[float32](a, b, name, t)
	case *float32:
		testComparePtrValue[float32](a, b, name, t)
	case []float32:
		testCompareSliceValue[float32](a, b, name, t)
	case float64:
		testCompareValue[float64](a, b, name, t)
	case *float64:
		testComparePtrValue[float64](a, b, name, t)
	case []float64:
		testCompareSliceValue[float64](a, b, name, t)
	case bool:
		testCompareValue[bool](a, b, name, t)
	case *bool:
		testComparePtrValue[bool](a, b, name, t)
	case []bool:
		testCompareSliceValue[bool](a, b, name, t)
	case time.Duration:
		testCompareValue[time.Duration](a, b, name, t)
	case *time.Duration:
		testComparePtrValue[time.Duration](a, b, name, t)
	case []time.Duration:
		testCompareSliceValue[time.Duration](a, b, name, t)
	case time.Time:
		testCompareValue[time.Time](a, b, name, t, timeNeq)
	case string:
		testCompareValue[string](a, b, name, t)
	case *string:
		testComparePtrValue[string](a, b, name, t)
	case []string:
		testCompareSliceValue[string](a, b, name, t)
	default:
		t.Fatalf("No test available for type: %s (%s)", fmt.Sprintf("%T", b), name)
	}

}

func TestToMinBytes(t *testing.T) {

	var i8 int8 = 42
	var i16 int16 = 84
	var i32 int32 = 126
	var i64 int64 = 168
	var u8 uint8 = 42
	var u16 uint16 = 84
	var u32 uint32 = 126
	var u64 uint64 = 168
	var f32 float32 = -42.42
	var f64 float64 = -84.84
	var s string = "Hello World"
	var bs []byte = []byte(s)

	type testData struct {
		V any
		T string
	}

	tests := []testData{
		{
			i8,
			"int8",
		},
		{
			&i8,
			"*int8",
		},
		{
			i16,
			"int16",
		},
		{
			&i16,
			"*int16",
		},
		{
			i32,
			"int32",
		},
		{
			&i32,
			"*int32",
		},
		{
			i64,
			"int64",
		},
		{
			&i64,
			"*int64",
		},
		{
			u8,
			"uint8",
		},
		{
			&u8,
			"*uint8",
		},
		{
			u16,
			"uint16",
		},
		{
			&u16,
			"*uint16",
		},
		{
			u32,
			"uint32",
		},
		{
			&u32,
			"*uint32",
		},
		{
			u64,
			"uint64",
		},
		{
			&u64,
			"*uint64",
		},
		{
			f32,
			"float32",
		},
		{
			&f32,
			"*float32",
		},
		{
			f64,
			"float64",
		},
		{
			&f64,
			"*float64",
		},
		{
			s,
			"string",
		},
		{
			&s,
			"*string",
		},
		{
			bs,
			"[]byte",
		},
	}

	approach := NewMinDataApproach()

	for _, test := range tests {
		b, err := approach.Pack(test.V)
		if err != nil {
			t.Fatalf("Unexpected error for %s: %v", test.T, err)
		}

		v, err := approach.Unpack(b)
		if err != nil {
			t.Fatalf("Unexpected error for %s: %v", test.T, err)
		}

		compareValue(v, test.V, test.T, t)
	}

	b, err := approach.Pack(nil)
	if err != nil {
		t.Fatalf("Unexpected error when packing nil: %v", err)
	}
	vNil, err := approach.Unpack(b)
	if err != nil {
		t.Fatalf("Unexpected error when unpacking nil: %v", err)
	}
	if vNil != nil {
		t.Fatalf("Expected nil; got: %v (%T)", vNil, vNil)
	}
}

func runSliceTest[T comparable](st TypeID, data []T, eleSize int64, t *testing.T) {

	b, err := packSimpleSliceMD(st, data)
	if err != nil {
		t.Fatalf("Unexpected error when packing slice: %v", err)
	}

	v, err := unpackSimpleSliceMD[T](b[1:], eleSize)
	if err != nil {
		t.Fatalf("Unexpected error when unpacking slice: %v", err)
	}

	vv, ok := v.([]T)
	if !ok {
		t.Fatalf("Unexpected type when unpacking slice: %T", v)
	}

	if len(data) != len(vv) {
		t.Fatalf("Unexpected size difference: expected: %d, got: %d", len(data), len(vv))
	}

	for i, vvv := range vv {
		if data[i] != vvv {
			t.Fatalf("Unexpected value difference at %d: expected: %v, got: %v", i, data[i], vvv)
		}
	}

}

func TestPackSlice(t *testing.T) {

	runSliceTest(Int8SliceType, []int8{43, 21, 54}, 1, t)
	runSliceTest(Int16SliceType, []int16{43, 21, 54}, 2, t)
	runSliceTest(Int32SliceType, []int32{43, 21, 54}, 4, t)
	runSliceTest(Int64SliceType, []int64{43, 21, 54}, 8, t)
	runSliceTest(Uint16SliceType, []uint16{43, 21, 54}, 2, t)
	runSliceTest(Uint32SliceType, []uint32{43, 21, 54}, 4, t)
	runSliceTest(Uint64SliceType, []uint64{43, 21, 54}, 8, t)
	runSliceTest(Float32SliceType, []float32{43, 21, 54}, 4, t)
	runSliceTest(Float64SliceType, []float64{43, 21, 54}, 8, t)
	runSliceTest(BoolSliceType, []bool{true, true, false}, 1, t)

}

func TestToBytes(t *testing.T) {

	type testData struct {
		V        any
		TypeName string
	}

	var i8 int8 = 42
	var i16 int16 = 42
	var i32 int32 = 42
	var i64 int64 = 42
	var u8 uint8 = 42
	var u16 uint16 = 42
	var u32 uint32 = 42
	var u64 uint64 = 42
	var f32 float32 = 42.42
	var f64 float64 = -42.42
	var bl bool = true
	var td time.Duration = 1234
	var bs []byte = []byte("Hello World")
	var ss []string = []string{"Hello", "World"}
	var is8 []int8 = []int8{1, 2, 3, 4}
	var is16 []int16 = []int16{1, 2, 3, 4}
	var is32 []int32 = []int32{1, 2, 3, 4}
	var is64 []int64 = []int64{1, 2, 3, 4}
	var uis16 []uint16 = []uint16{1, 2, 3, 4}
	var uis32 []uint32 = []uint32{1, 2, 3, 4}
	var uis64 []uint64 = []uint64{1, 2, 3, 4}
	var fs32 []float32 = []float32{1, 2, 3, 4}
	var fs64 []float64 = []float64{1, 2, 3, 4}
	var bbs []bool = []bool{false, true, true, false}
	var tds []time.Duration = []time.Duration{1, 2, 3, 4}
	var tm time.Time = time.Now()

	tests := []testData{
		{
			nil,
			"blah",
		},
		{
			i8,
			"int8",
		},
		{
			&i8,
			"*int8",
		},
		{
			i16,
			"int16",
		},
		{
			&i16,
			"*int16",
		},
		{
			i32,
			"int32",
		},
		{
			&i32,
			"*int32",
		},
		{
			i64,
			"int64",
		},
		{
			&i64,
			"*int64",
		},
		{
			u8,
			"uint8",
		},
		{
			&u8,
			"*uint8",
		},
		{
			u16,
			"uint16",
		},
		{
			&u16,
			"*uint16",
		},
		{
			u32,
			"uint32",
		},
		{
			&u32,
			"*uint32",
		},
		{
			u64,
			"uint64",
		},
		{
			&u64,
			"*uint64",
		},
		{
			f32,
			"float32",
		},
		{
			&f32,
			"*float32",
		},
		{
			f64,
			"float64",
		},
		{
			&f64,
			"*float64",
		},
		{
			bl,
			"bool",
		},
		{
			&bl,
			"*bool",
		},
		{
			td,
			"time.Duration",
		},
		{
			&td,
			"*time.Duration",
		},
		{
			bs,
			"[]byte",
		},
		{
			is8,
			"[]int8",
		},
		{
			is16,
			"[]int16",
		},
		{
			is32,
			"[]int32",
		},
		{
			is64,
			"[]int64",
		},
		{
			uis16,
			"[]uint16",
		},
		{
			uis32,
			"[]uint32",
		},
		{
			uis64,
			"[]uint64",
		},
		{
			fs32,
			"[]float32",
		},
		{
			fs64,
			"[]float64",
		},
		{
			bbs,
			"[]bool",
		},
		{
			tds,
			"[]time.Duration",
		},
		{
			ss,
			"[]string",
		},
		{
			tm,
			"time.Time",
		},
	}

	approach := NewMinDataApproach()

	for _, test := range tests {

		b, _, err := ToBytes(test.V, WithSerialisationApproach(approach))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		v, err := FromBytes(b, approach)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		compareValue(v, test.V, test.TypeName, t)
	}
}

func TestToBytes_1(t *testing.T) {

	type testData struct {
		V        any
		TypeName string
	}

	var i8 int8 = 42
	var ss []string = []string{"This", "is", "not", "encrypted"}
	var ss2 []string = []string{"01234567890123456789012345678901234567890123456789"}

	tests := []testData{
		{
			ss,
			"[]string",
		},
		{
			ss2,
			"[]string",
		},
		{
			i8,
			"int8",
		},
	}

	approach := NewMinDataApproach()

	key := []byte("01234567890123456789012345678912")

	for _, test := range tests {

		b, _, err := ToBytes(test.V, WithSerialisationApproach(approach), WithAESGCMEncryption(key))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		v, err := FromBytes(b, approach, WithAESGCMEncryption(key))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		compareValue(v, test.V, test.TypeName, t)
	}
}

func TestToBytes_2(t *testing.T) {

	type testData struct {
		V        any
		TypeName string
	}

	var i8 int8 = 42

	tests := []testData{
		{
			i8,
			"int8",
		},
	}

	for _, test := range tests {

		b, name, err := ToBytes(test.V)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		a, err := GetApproach(name)
		if err != nil {
			t.Fatalf("Unexpected error retrieving Approach for name '%s': %v", name, err)
		}
		v, err := FromBytes(b, a)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		compareValue(v, test.V, test.TypeName, t)
	}
}

func TestToBytes_3(t *testing.T) {

	var i8 int8 = 42

	b1, name1, _ := ToBytes(i8)
	b2, name2, _ := ToBytes(i8, WithSerialisationApproach(defaultSerialisationApproach))

	if defaultSerialisationApproach.Name() != name2 {
		t.Fatalf("Unexpected name returned: expected: %s, got: %s", defaultSerialisationApproach.Name(), name2)
	}

	if name1 != name2 {
		t.Fatalf("Unexpected difference in names: got: %s rather than: %s", name1, name2)
	}

	if !bytes.Equal(b1, b2) {
		t.Fatalf("Unexpected variation in byte slices")
	}
}

func TestToBytesMany(t *testing.T) {

	tests := [][]any{
		{
			int64(0), float32(-42),
		},
		{
			int64(2), []string{"Hello", "World"},
		},
		{},
		{nil},
		{
			int64(-1), "", "Hello", nil, true,
		},
	}

	for i, test := range tests {

		b, name, err := ToBytesMany(test, WithSerialisationApproach(defaultSerialisationApproach))
		if err != nil {
			t.Fatalf("(%d) Unexpected error when serialising: %v", i, err)
		}
		if name != defaultSerialisationApproach.Name() {
			t.Fatalf("(%d) Unexpected difference in Approach name: expected: %s, got: %s", i, defaultSerialisationApproach.Name(), name)
		}

		v, err := FromBytesMany(b, defaultSerialisationApproach)
		if err != nil {
			t.Fatalf("(%d) Unexpected error when deserialising: %v", i, err)
		}

		if len(v) != len(test) {
			t.Fatalf("(%d) Unexpected error in output length: expected: %d, got: %d", i, len(test), len(v))
		}

		for j := 0; j < len(test); j++ {
			compareValue(v[j], test[j], fmt.Sprintf("%T", test[j]), t)

		}
	}
}

func TestToBytesMany_1(t *testing.T) {

	tests := [][]any{
		{
			int64(0), float32(-42),
		},
		{
			int64(2), []string{"Hello", "World"},
		},
		{
			[]byte("Hello"), []byte("World"),
		},
	}

	key := []byte("01234567890123456789012345678901")

	for i, test := range tests {

		b, name, err := ToBytesMany(test, WithSerialisationApproach(defaultSerialisationApproach), WithAESGCMEncryption(key))
		if err != nil {
			t.Fatalf("(%d) Unexpected error when serialising: %v", i, err)
		}
		if name != defaultSerialisationApproach.Name() {
			t.Fatalf("(%d) Unexpected difference in Approach name: expected: %s, got: %s", i, defaultSerialisationApproach.Name(), name)
		}

		v, err := FromBytesMany(b, defaultSerialisationApproach, WithAESGCMEncryption(key))
		if err != nil {
			t.Fatalf("(%d) Unexpected error when deserialising: %v", i, err)
		}

		if len(v) != len(test) {
			t.Fatalf("(%d) Unexpected error in output length: expected: %d, got: %d", i, len(test), len(v))
		}

		for j := 0; j < len(test); j++ {
			compareValue(v[j], test[j], fmt.Sprintf("%T", test[j]), t)

		}
	}
}
