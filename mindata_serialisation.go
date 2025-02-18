package serialise

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
		return pack(stringType, []byte(*v))
	case []byte:
		return pack(byteSliceType, v)
	default:
		panic(fmt.Sprintf("Bums! (%T)", v))
	}
}

// Unpack deserialises an instance from the byte slice
func (m *minData) Unpack(data []byte, opts ...func(opt *TypeRegistryOptions)) (any, error) {

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
	case int32Type:
		return unpackMD[int32](data)
	case pint32Type:
		return unpackPtr(new(int32), data)
	case int64Type:
		return unpackMD[int64](data)
	case pint64Type:
		return unpackPtr(new(int64), data)
	case uint8Type:
		return unpackMD[uint8](data)
	case puint8Type:
		return unpackPtr(new(uint8), data)
	case uint16Type:
		return unpackMD[uint16](data)
	case puint16Type:
		return unpackPtr(new(uint16), data)
	case uint32Type:
		return unpackMD[uint32](data)
	case puint32Type:
		return unpackPtr(new(uint32), data)
	case uint64Type:
		return unpackMD[uint64](data)
	case puint64Type:
		return unpackPtr(new(uint64), data)
	case float32Type:
		return unpackMD[float32](data)
	case pfloat32Type:
		return unpackPtr(new(float32), data)
	case float64Type:
		return unpackMD[float64](data)
	case pfloat64Type:
		return unpackPtr(new(float64), data)
	case boolType:
		return unpackMD[bool](data)
	case pboolType:
		return unpackPtr(new(bool), data)
	case durationType:
		return unpackMD[time.Duration](data)
	case pdurationType:
		return unpackPtr(new(time.Duration), data)
	case stringType:
		return string(data[1:]), nil
	case pstringType:
		s := string(data[1:])
		return &s, nil
	case byteSliceType:
		return data[1:], nil

	default:
		panic(fmt.Sprintf("Bums Again! (%d)", t))
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
