package serialise

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

// MinDataApproach creates an Approach instance
// that uses minimum byte storage during serialisation
// Only basic built-in types are supported.
func NewMinDataApproach() Approach {
	return &minData{}
}

type minData struct {
}

// Name of the approach
func (m *minData) Name() string {
	return "MinData"
}

// IsSerialisable returns true if an instance of the specified type
// can be serialised
func (m *minData) IsSerialisable(v any) (ok bool) {
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
func (m *minData) Pack(data any) ([]byte, error) {

	pack := func(t int8, data any) ([]byte, error) {
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

	packByteSliceSlice := func(t int8, data [][]byte) ([]byte, error) {
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

	if data == nil {
		return pack(nilType, nil)
	}

	switch v := data.(type) {
	case int8:
		return pack(int8Type, data)
	case *int8:
		return pack(pint8Type, data)
	case []int8:
		return packSimpleSliceMD(int8SliceType, v)
	case int16:
		return pack(int16Type, data)
	case *int16:
		return pack(pint16Type, data)
	case []int16:
		return packSimpleSliceMD(int16SliceType, v)
	case int32:
		return pack(int32Type, data)
	case *int32:
		return pack(pint32Type, data)
	case []int32:
		return packSimpleSliceMD(int32SliceType, v)
	case int64:
		return pack(int64Type, data)
	case *int64:
		return pack(pint64Type, data)
	case []int64:
		return packSimpleSliceMD(int64SliceType, v)
	case uint8:
		return pack(uint8Type, data)
	case *uint8:
		return pack(puint8Type, data)
	case uint16:
		return pack(uint16Type, data)
	case *uint16:
		return pack(puint16Type, data)
	case []uint16:
		return packSimpleSliceMD(uint16SliceType, v)
	case uint32:
		return pack(uint32Type, data)
	case *uint32:
		return pack(puint32Type, data)
	case []uint32:
		return packSimpleSliceMD(uint32SliceType, v)
	case uint64:
		return pack(uint64Type, data)
	case *uint64:
		return pack(puint64Type, data)
	case []uint64:
		return packSimpleSliceMD(uint64SliceType, v)
	case float32:
		return pack(float32Type, data)
	case *float32:
		return pack(pfloat32Type, data)
	case []float32:
		return packSimpleSliceMD(float32SliceType, v)
	case float64:
		return pack(float64Type, data)
	case *float64:
		return pack(pfloat64Type, data)
	case []float64:
		return packSimpleSliceMD(float64SliceType, v)
	case bool:
		return pack(boolType, data)
	case *bool:
		return pack(pboolType, data)
	case []bool:
		return packSimpleSliceMD(boolSliceType, v)
	case time.Duration:
		return pack(durationType, data)
	case *time.Duration:
		return pack(pdurationType, data)
	case []time.Duration:
		return packSimpleSliceMD(durationSliceType, v)
	case string:
		return pack(stringType, []byte(v))
	case *string:
		return pack(pstringType, []byte(*v))
	case []string:
		var bss [][]byte = make([][]byte, len(v))
		for i := 0; i < len(v); i++ {
			bss[i] = []byte(v[i])
		}
		return packByteSliceSlice(stringSliceType, bss)
	case []byte:
		return pack(byteSliceType, v)
	case [][]byte:
		return packByteSliceSlice(byteSliceSliceType, v)
	default:
		return nil, ErrMinDataTypeNotSerialisable
	}
}

// ErrMinDataTypeNotDeserialisable is raised if a variable type is not deserialisable by the MinData approach
var ErrMinDataTypeNotDeserialisable = errors.New("type specified within the data is not deserialisable")

// ErrUnexpectedDeserialisationError is raised if a panic is generated during deserialisation
var ErrUnexpectedDeserialisationError = errors.New("unexpected deserialisation failure - possible corrupted data provided")

// Unpack deserialises an instance from the byte slice
func (m *minData) Unpack(data []byte, opts ...func(opt *TypeRegistryOptions)) (output any, e error) {

	defer func() {
		if r := recover(); r != nil {
			output = nil
			e = ErrUnexpectedDeserialisationError
		}
	}()

	var t int8
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

	switch t {
	case nilType:
		return nil, nil
	case int8Type:
		return unpackMD[int8](data)
	case pint8Type:
		return unpackPtr(new(int8), data)
	case int8SliceType:
		return unpackSimpleSliceMD[int8](data[1:], 1)
	case int16Type:
		return unpackMD[int16](data)
	case pint16Type:
		return unpackPtr(new(int16), data)
	case int16SliceType:
		return unpackSimpleSliceMD[int16](data[1:], 2)
	case int32Type:
		return unpackMD[int32](data)
	case pint32Type:
		return unpackPtr(new(int32), data)
	case int32SliceType:
		return unpackSimpleSliceMD[int32](data[1:], 4)
	case int64Type:
		return unpackMD[int64](data)
	case pint64Type:
		return unpackPtr(new(int64), data)
	case int64SliceType:
		return unpackSimpleSliceMD[int64](data[1:], 8)
	case uint8Type:
		return unpackMD[uint8](data)
	case puint8Type:
		return unpackPtr(new(uint8), data)
	case uint16Type:
		return unpackMD[uint16](data)
	case puint16Type:
		return unpackPtr(new(uint16), data)
	case uint16SliceType:
		return unpackSimpleSliceMD[uint16](data[1:], 2)
	case uint32Type:
		return unpackMD[uint32](data)
	case puint32Type:
		return unpackPtr(new(uint32), data)
	case uint32SliceType:
		return unpackSimpleSliceMD[uint32](data[1:], 4)
	case uint64Type:
		return unpackMD[uint64](data)
	case puint64Type:
		return unpackPtr(new(uint64), data)
	case uint64SliceType:
		return unpackSimpleSliceMD[uint64](data[1:], 8)
	case float32Type:
		return unpackMD[float32](data)
	case pfloat32Type:
		return unpackPtr(new(float32), data)
	case float32SliceType:
		return unpackSimpleSliceMD[float32](data[1:], 4)
	case float64Type:
		return unpackMD[float64](data)
	case pfloat64Type:
		return unpackPtr(new(float64), data)
	case float64SliceType:
		return unpackSimpleSliceMD[float64](data[1:], 8)
	case boolType:
		return unpackMD[bool](data)
	case pboolType:
		return unpackPtr(new(bool), data)
	case boolSliceType:
		return unpackSimpleSliceMD[bool](data[1:], 1)
	case durationType:
		return unpackMD[time.Duration](data)
	case pdurationType:
		return unpackPtr(new(time.Duration), data)
	case durationSliceType:
		return unpackSimpleSliceMD[time.Duration](data[1:], 8)
	case stringType:
		return string(data[1:]), nil
	case pstringType:
		s := string(data[1:])
		return &s, nil
	case byteSliceType:
		return data[1:], nil
	case byteSliceSliceType:
		return unpackByteSliceSlice(data[1:])
	case stringSliceType:
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

func packSimpleSliceMD[T any](t int8, data []T) ([]byte, error) {
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
