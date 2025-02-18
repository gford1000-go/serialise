package serialise

import (
	"fmt"
	"testing"
	"time"
)

func TestGetType(t *testing.T) {

	type testData struct {
		V        any
		TypeName string
	}

	var i int = 42
	var i64 int64 = 42
	var s string = "Hello World"
	var td testData
	var tdp *testData = &testData{}
	var tm = time.Now()

	tests := []testData{
		{
			i,
			"int",
		},
		{
			&i,
			"*int",
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
			s,
			"string",
		},
		{
			&s,
			"*string",
		},
		{
			tm,
			"time.Time",
		},
		{
			&tm,
			"*time.Time",
		},
		{
			td,
			"serialise.testData",
		},
		{
			&td,
			"*serialise.testData",
		},
		{
			tdp,
			"*serialise.testData",
		},
		{
			&tdp,
			"**serialise.testData",
		},
	}

	for _, test := range tests {
		name := getTypeName(test.V)
		if name != test.TypeName {
			t.Fatalf("Unexpected mismatch in type name: expected: %s, got: %s", test.TypeName, name)
		}
	}

	testRegistry := NewTypeRegistry()

	f := func(o *TypeRegistryOptions) {
		o.Registry = testRegistry
	}

	for _, test := range tests {
		RegisterType(test.V, f)
	}

	for _, test := range tests {
		ty, err := GetRegisteredType(test.TypeName, f)
		if err != nil {
			t.Fatalf("Unexpected error: %v for type: %s", err, test.TypeName)
		}
		if fmt.Sprintf("%v", ty) != test.TypeName {
			t.Fatalf("Mismatch in type: expected: %s, got: %s", test.TypeName, fmt.Sprintf("%v", ty))
		}
	}

	for _, test := range tests {
		v, err := CreateInstance(test.TypeName, f)
		if err != nil {
			t.Fatalf("Unexpected error: %v for type: %s", err, test.TypeName)
		}
		if getTypeName(v) != test.TypeName {
			t.Fatalf("Mismatch in type: expected: %s, got: %s", test.TypeName, fmt.Sprintf("%T", v))
		}
	}
}

func testCompareValue[T comparable](a, b any, name string, t *testing.T) {
	switch v := b.(type) {
	case T:
		if aa, ok := a.(T); ok {
			if aa != v {
				t.Fatalf("Data mismatch: expected %v, got: %v", v, aa)
			}
		} else {
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

		switch v.(type) {
		case int8:
			testCompareValue[int8](v, i8, test.T, t)
		case *int8:
			testComparePtrValue[int8](v, &i8, test.T, t)
		case int16:
			testCompareValue[int16](v, i16, test.T, t)
		case *int16:
			testComparePtrValue[int16](v, &i16, test.T, t)
		case int32:
			testCompareValue[int32](v, i32, test.T, t)
		case *int32:
			testComparePtrValue[int32](v, &i32, test.T, t)
		case int64:
			testCompareValue[int64](v, i64, test.T, t)
		case *int64:
			testComparePtrValue[int64](v, &i64, test.T, t)
		case uint8:
			testCompareValue[uint8](v, u8, test.T, t)
		case *uint8:
			testComparePtrValue[uint8](v, &u8, test.T, t)
		case uint16:
			testCompareValue[uint16](v, u16, test.T, t)
		case *uint16:
			testComparePtrValue[uint16](v, &u16, test.T, t)
		case uint32:
			testCompareValue[uint32](v, u32, test.T, t)
		case *uint32:
			testComparePtrValue[uint32](v, &u32, test.T, t)
		case uint64:
			testCompareValue[uint64](v, u64, test.T, t)
		case *uint64:
			testComparePtrValue[uint64](v, &u64, test.T, t)
		case float32:
			testCompareValue[float32](v, f32, test.T, t)
		case *float32:
			testComparePtrValue[float32](v, &f32, test.T, t)
		case float64:
			testCompareValue[float64](v, f64, test.T, t)
		case *float64:
			testComparePtrValue[float64](v, &f64, test.T, t)
		case string:
			testCompareValue[string](v, s, test.T, t)
		case *string:
			testComparePtrValue[string](v, &s, test.T, t)
		case []byte:
			testCompareSliceValue[byte](v, bs, test.T, t)
		default:
			t.Fatalf("Unexpected type for %s: %T", test.T, v)
		}
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

func runSliceTest[T comparable](st int8, data []T, eleSize int64, t *testing.T) {

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

	runSliceTest(int8SliceType, []int8{43, 21, 54}, 1, t)
	runSliceTest(int16SliceType, []int16{43, 21, 54}, 2, t)
	runSliceTest(int32SliceType, []int32{43, 21, 54}, 4, t)
	runSliceTest(int64SliceType, []int64{43, 21, 54}, 8, t)
	runSliceTest(uint16SliceType, []uint16{43, 21, 54}, 2, t)
	runSliceTest(uint32SliceType, []uint32{43, 21, 54}, 4, t)
	runSliceTest(uint64SliceType, []uint64{43, 21, 54}, 8, t)
	runSliceTest(float32SliceType, []float32{43, 21, 54}, 4, t)
	runSliceTest(float64SliceType, []float64{43, 21, 54}, 8, t)
	runSliceTest(boolSliceType, []bool{true, true, false}, 1, t)

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

	compareValue := func(a, b any, name string) {
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
		case []string:
			testCompareSliceValue[string](a, b, name, t)
		default:
			t.Fatalf("No test available for type: %s (%s)", fmt.Sprintf("%T", b), name)
		}

	}

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
	}

	approach := NewMinDataApproach()

	testRegistry := NewTypeRegistry()
	testRegistry.AddTypeOf(ss)
	testRegistry.AddTypeOf(is8)

	f := func(o *TypeRegistryOptions) {
		o.Registry = testRegistry
	}

	for _, test := range tests {

		b, _, err := ToBytes(test.V, WithSerialisationApproach(approach))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		v, err := FromBytes(b, approach, WithTypeRegistryOptions(f))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		compareValue(v, test.V, test.TypeName)
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

	compareValue := func(a, b any, name string) {
		if b == nil {
			if a != nil {
				t.Fatalf("Mismatch in <nil>")
			}
			return
		}

		switch b.(type) {
		case int8:
			testCompareValue[int8](a, b, name, t)
		case []string:
			testCompareSliceValue[string](a, b, name, t)
		default:
			t.Fatalf("No test available for type: %s (%s)", fmt.Sprintf("%T", b), name)
		}

	}

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

	testRegistry := NewTypeRegistry()
	testRegistry.AddTypeOf(ss)

	f := func(o *TypeRegistryOptions) {
		o.Registry = testRegistry
	}

	key := []byte("01234567890123456789012345678912")

	for _, test := range tests {

		b, _, err := ToBytes(test.V, WithSerialisationApproach(approach), WithAESGCMEncryption(key))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		v, err := FromBytes(b, approach, WithTypeRegistryOptions(f), WithAESGCMEncryption(key))
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		compareValue(v, test.V, test.TypeName)
	}
}

func TestToBytes_2(t *testing.T) {

	type testData struct {
		V        any
		TypeName string
	}

	var i8 int8 = 42

	compareValue := func(a, b any, name string) {
		if b == nil {
			if a != nil {
				t.Fatalf("Mismatch in <nil>")
			}
			return
		}

		switch b.(type) {
		case int8:
			testCompareValue[int8](a, b, name, t)
		default:
			t.Fatalf("No test available for type: %s (%s)", fmt.Sprintf("%T", b), name)
		}

	}

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
		compareValue(v, test.V, test.TypeName)
	}
}
