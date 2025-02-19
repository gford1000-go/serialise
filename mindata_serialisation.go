package serialise

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// MinDataVersion describes a version of a MinData serialisation implementation
// All breaking changes to serialisation will trigger an increment, to ensure
// backwards compatibility to any consumers of existing versions.
type MinDataVersion int8

const (
	UnknownVersion MinDataVersion = iota
	V1
	OutOfRange
)

var defaultVersion MinDataVersion = V1

// NewMinDataApproach creates an instance of the
// current default version of the MinData serialisation
func NewMinDataApproach() Approach {
	return NewMinDataApproachWithVersion(defaultVersion)
}

// NewMinDataApproachWithVersion creates an instance of the
// specified version of MinData serialisation
func NewMinDataApproachWithVersion(version MinDataVersion) Approach {
	switch version {
	case V1:
		name := "MD1"
		return &minDataV1{name: name}
	default:
		panic(fmt.Sprintf("Illegal MinDataVersion passed to NewMinDataApproach (%d)", version))
	}
}

type minDataV1 struct {
	name string
}

// Name of the approach
func (m *minDataV1) Name() string {
	return m.name
}

// IsSerialisable returns true if an instance of the specified type
// can be serialised
func (m *minDataV1) IsSerialisable(v any) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	_, err := m.Pack(v)
	return err == nil
}

// ErrMinDataTypeNotSerialisable is raised if a variable type is not serialisable by the MinData approach
var ErrMinDataTypeNotSerialisable = errors.New("type of argument is not serialisable")

// Pack serialises the instance to a byte slice
func (m *minDataV1) Pack(data any) ([]byte, error) {

	pack := func(t TypeID, data any) ([]byte, error) {
		var tbuf bytes.Buffer
		var vbuf bytes.Buffer

		err := binary.Write(&tbuf, binary.LittleEndian, t)
		if err != nil {
			return nil, err
		}
		if data != nil {
			if b, ok := data.([]byte); ok {
				return append(tbuf.Bytes(), b...), nil
			}
			err = binary.Write(&tbuf, binary.LittleEndian, data)
			if err != nil {
				return nil, err
			}
		}
		return append(tbuf.Bytes(), vbuf.Bytes()...), nil
	}

	packByteSliceSlice := func(t TypeID, data [][]byte) ([]byte, error) {
		var buf bytes.Buffer

		err := binary.Write(&buf, binary.LittleEndian, t)
		if err != nil {
			return nil, err
		}
		if data != nil {
			var size int64 = int64(len(data))
			err = binary.Write(&buf, binary.LittleEndian, size)
			if err != nil {
				return nil, err
			}
			for _, b := range data {
				size = int64(len(b))
				err = binary.Write(&buf, binary.LittleEndian, size)
				if err != nil {
					return nil, err
				}
				if size > 0 {
					err = binary.Write(&buf, binary.LittleEndian, b)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		return buf.Bytes(), nil
	}

	packTime := func(t TypeID, tm *time.Time) ([]byte, error) {
		var buf bytes.Buffer

		err := binary.Write(&buf, binary.LittleEndian, t)
		if err != nil {
			return nil, err
		}

		b, err := tm.MarshalBinary()
		if err != nil {
			return nil, err
		}

		return append(buf.Bytes(), b...), nil
	}

	if data == nil {
		return pack(NilType, nil)
	}

	switch v := data.(type) {
	case int8:
		return pack(Int8Type, data)
	case *int8:
		return pack(Pint8Type, data)
	case []int8:
		return packSimpleSliceMD(Int8SliceType, v)
	case int16:
		return pack(Int16Type, data)
	case *int16:
		return pack(Pint16Type, data)
	case []int16:
		return packSimpleSliceMD(Int16SliceType, v)
	case int32:
		return pack(Int32Type, data)
	case *int32:
		return pack(Pint32Type, data)
	case []int32:
		return packSimpleSliceMD(Int32SliceType, v)
	case int64:
		return pack(Int64Type, data)
	case *int64:
		return pack(Pint64Type, data)
	case []int64:
		return packSimpleSliceMD(Int64SliceType, v)
	case uint8:
		return pack(Uint8Type, data)
	case *uint8:
		return pack(Puint8Type, data)
	case uint16:
		return pack(Uint16Type, data)
	case *uint16:
		return pack(Puint16Type, data)
	case []uint16:
		return packSimpleSliceMD(Uint16SliceType, v)
	case uint32:
		return pack(Uint32Type, data)
	case *uint32:
		return pack(Puint32Type, data)
	case []uint32:
		return packSimpleSliceMD(Uint32SliceType, v)
	case uint64:
		return pack(Uint64Type, data)
	case *uint64:
		return pack(Puint64Type, data)
	case []uint64:
		return packSimpleSliceMD(Uint64SliceType, v)
	case float32:
		return pack(Float32Type, data)
	case *float32:
		return pack(Pfloat32Type, data)
	case []float32:
		return packSimpleSliceMD(Float32SliceType, v)
	case float64:
		return pack(Float64Type, data)
	case *float64:
		return pack(Pfloat64Type, data)
	case []float64:
		return packSimpleSliceMD(Float64SliceType, v)
	case bool:
		return pack(BoolType, data)
	case *bool:
		return pack(PboolType, data)
	case []bool:
		return packSimpleSliceMD(BoolSliceType, v)
	case time.Duration:
		return pack(DurationType, data)
	case *time.Duration:
		return pack(PdurationType, data)
	case []time.Duration:
		return packSimpleSliceMD(DurationSliceType, v)
	case time.Time:
		return packTime(TimeType, &v)
	case *time.Time:
		return packTime(PtimeType, v)
	case string:
		return pack(StringType, []byte(v))
	case *string:
		return pack(PstringType, []byte(*v))
	case []string:
		var bss [][]byte = make([][]byte, len(v))
		for i := 0; i < len(v); i++ {
			bss[i] = []byte(v[i])
		}
		return packByteSliceSlice(StringSliceType, bss)
	case []byte:
		return pack(ByteSliceType, v)
	case [][]byte:
		return packByteSliceSlice(ByteSliceSliceType, v)
	default:
		return nil, ErrMinDataTypeNotSerialisable
	}
}

// ErrMinDataTypeNotDeserialisable is raised if a variable type is not deserialisable by the MinData approach
var ErrMinDataTypeNotDeserialisable = errors.New("type specified within the data is not deserialisable")

// ErrUnexpectedDeserialisationError is raised if a panic is generated during deserialisation
var ErrUnexpectedDeserialisationError = errors.New("unexpected deserialisation failure - possible corrupted data provided")

// Unpack deserialises an instance from the byte slice
// func (m *minData) Unpack(data []byte, opts ...func(opt *TypeRegistryOptions)) (output any, e error) {
func (m *minDataV1) Unpack(data []byte) (output any, e error) {

	defer func() {
		if r := recover(); r != nil {
			output = nil
			e = ErrUnexpectedDeserialisationError
		}
	}()

	var t TypeID
	err := binary.Read(bytes.NewBuffer(data[0:1]), binary.LittleEndian, &t)
	if err != nil {
		return nil, err
	}

	unpackPtr := func(v any, data []byte) (any, error) {
		err := binary.Read(bytes.NewBuffer(data[1:]), binary.LittleEndian, v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	unpackByteSliceSlice := func(data []byte) ([][]byte, error) {
		var size int64
		err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &size)
		if err != nil {
			return nil, err
		}

		bss := make([][]byte, size)
		var i int64
		var offset int64 = 8
		for i = 0; i < size; i++ {
			var itemSize int64
			err := binary.Read(bytes.NewBuffer(data[offset:]), binary.LittleEndian, &itemSize)
			if err != nil {
				return nil, err
			}
			bss[i] = data[offset+8 : offset+8+itemSize]
			offset += 8 + itemSize
		}
		return bss, nil
	}

	unpackTime := func(data []byte) (any, error) {
		tm := new(time.Time)
		err := tm.UnmarshalBinary(data)
		if err != nil {
			return nil, err
		}
		return tm, nil
	}

	switch t {
	case NilType:
		return nil, nil
	case Int8Type:
		return unpackMD[int8](data)
	case Pint8Type:
		return unpackPtr(new(int8), data)
	case Int8SliceType:
		return unpackSimpleSliceMD[int8](data[1:], 1)
	case Int16Type:
		return unpackMD[int16](data)
	case Pint16Type:
		return unpackPtr(new(int16), data)
	case Int16SliceType:
		return unpackSimpleSliceMD[int16](data[1:], 2)
	case Int32Type:
		return unpackMD[int32](data)
	case Pint32Type:
		return unpackPtr(new(int32), data)
	case Int32SliceType:
		return unpackSimpleSliceMD[int32](data[1:], 4)
	case Int64Type:
		return unpackMD[int64](data)
	case Pint64Type:
		return unpackPtr(new(int64), data)
	case Int64SliceType:
		return unpackSimpleSliceMD[int64](data[1:], 8)
	case Uint8Type:
		return unpackMD[uint8](data)
	case Puint8Type:
		return unpackPtr(new(uint8), data)
	case Uint16Type:
		return unpackMD[uint16](data)
	case Puint16Type:
		return unpackPtr(new(uint16), data)
	case Uint16SliceType:
		return unpackSimpleSliceMD[uint16](data[1:], 2)
	case Uint32Type:
		return unpackMD[uint32](data)
	case Puint32Type:
		return unpackPtr(new(uint32), data)
	case Uint32SliceType:
		return unpackSimpleSliceMD[uint32](data[1:], 4)
	case Uint64Type:
		return unpackMD[uint64](data)
	case Puint64Type:
		return unpackPtr(new(uint64), data)
	case Uint64SliceType:
		return unpackSimpleSliceMD[uint64](data[1:], 8)
	case Float32Type:
		return unpackMD[float32](data)
	case Pfloat32Type:
		return unpackPtr(new(float32), data)
	case Float32SliceType:
		return unpackSimpleSliceMD[float32](data[1:], 4)
	case Float64Type:
		return unpackMD[float64](data)
	case Pfloat64Type:
		return unpackPtr(new(float64), data)
	case Float64SliceType:
		return unpackSimpleSliceMD[float64](data[1:], 8)
	case BoolType:
		return unpackMD[bool](data)
	case PboolType:
		return unpackPtr(new(bool), data)
	case BoolSliceType:
		return unpackSimpleSliceMD[bool](data[1:], 1)
	case DurationType:
		return unpackMD[time.Duration](data)
	case PdurationType:
		return unpackPtr(new(time.Duration), data)
	case DurationSliceType:
		return unpackSimpleSliceMD[time.Duration](data[1:], 8)
	case TimeType:
		tm, err := unpackTime(data[1:])
		return *(tm.(*time.Time)), err
	case PtimeType:
		return unpackTime(data[1:])
	case StringType:
		return string(data[1:]), nil
	case PstringType:
		s := string(data[1:])
		return &s, nil
	case ByteSliceType:
		return data[1:], nil
	case ByteSliceSliceType:
		return unpackByteSliceSlice(data[1:])
	case StringSliceType:
		bss, err := unpackByteSliceSlice(data[1:])
		if err != nil {
			return nil, err
		}
		ss := make([]string, len(bss))
		for i := 0; i < len(bss); i++ {
			ss[i] = string(bss[i])
		}
		return ss, nil
	default:
		return nil, ErrMinDataTypeNotDeserialisable
	}
}

func packSimpleSliceMD[T any](t TypeID, data []T) ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, t)
	if err != nil {
		return nil, err
	}

	err = binary.Write(&buf, binary.LittleEndian, int64(len(data)))
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		err = binary.Write(&buf, binary.LittleEndian, d)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func unpackSimpleSliceMD[T any](data []byte, eleSize int64) (any, error) {
	var size int64
	err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &size)
	if err != nil {
		return nil, err
	}

	var v = make([]T, size)

	var i int64
	for i = 0; i < size; i++ {
		var vv T
		err := binary.Read(bytes.NewBuffer(data[8+i*eleSize:]), binary.LittleEndian, &vv)
		if err != nil {
			return nil, err
		}
		v[i] = vv
	}

	return v, nil
}

func unpackMD[T any](data []byte) (T, error) {
	unpack := func(v any, data []byte) error {
		err := binary.Read(bytes.NewBuffer(data[1:]), binary.LittleEndian, v)
		if err != nil {
			return err
		}
		return nil
	}

	var v T
	err := unpack(&v, data)
	return v, err
}
