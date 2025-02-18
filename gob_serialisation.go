package serialise

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"time"
)

// NewGOBApproach creates an Approach instance
// that uses gob serialisation.
func NewGOBApproach() Approach {
	return &gobApproach{}
}

type gobApproach struct {
}

// Name of the approach
func (g *gobApproach) Name() string {
	return "GOB"
}

// Pack serialises the instance to a byte slice
func (g *gobApproach) Pack(data any) ([]byte, error) {
	gd, err := g.toGobDataBytes(data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(gd); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

// Unpack deserialises an instance from the byte slice
func (g *gobApproach) Unpack(data []byte, opts ...func(opt *TypeRegistryOptions)) (any, error) {
	var buf = bytes.NewBuffer(data)

	decoder := gob.NewDecoder(buf)

	var gd gobData

	err := decoder.Decode(&gd)
	if err != nil {
		return nil, err
	}

	v, err := g.fromGobDataBytes(&gd, opts...)
	if err != nil {
		return nil, err
	}
	return v, err

}

type gobData struct {
	DataType int8
	TypeName string
	Data     []byte
}

// toGobDataBytes serialises data types to []byte using gob encoding
func (g *gobApproach) toGobDataBytes(data any) (*gobData, error) {
	if data == nil {
		return &gobData{DataType: nilType, Data: []byte{}}, nil
	}

	var buf bytes.Buffer

	switch v := data.(type) {
	case []byte:
		return &gobData{DataType: byteSliceType, Data: v}, nil
	case int8:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: int8Type, Data: buf.Bytes()}, err
	case *int8:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: pint8Type, Data: buf.Bytes()}, err
	case int16:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: int16Type, Data: buf.Bytes()}, err
	case *int16:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: pint16Type, Data: buf.Bytes()}, err
	case int32:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: int32Type, Data: buf.Bytes()}, err
	case *int32:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: pint32Type, Data: buf.Bytes()}, err
	case int64:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: int64Type, Data: buf.Bytes()}, err
	case *int64:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: pint64Type, Data: buf.Bytes()}, err
	case uint8:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: uint8Type, Data: buf.Bytes()}, err
	case *uint8:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: puint8Type, Data: buf.Bytes()}, err
	case uint16:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: uint16Type, Data: buf.Bytes()}, err
	case *uint16:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: puint16Type, Data: buf.Bytes()}, err
	case uint32:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: uint32Type, Data: buf.Bytes()}, err
	case *uint32:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: puint32Type, Data: buf.Bytes()}, err
	case uint64:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: uint64Type, Data: buf.Bytes()}, err
	case *uint64:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: puint64Type, Data: buf.Bytes()}, err
	case float32:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: float32Type, Data: buf.Bytes()}, err
	case *float32:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: pfloat32Type, Data: buf.Bytes()}, err
	case float64:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: float64Type, Data: buf.Bytes()}, err
	case *float64:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: pfloat64Type, Data: buf.Bytes()}, err
	case bool:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: boolType, Data: buf.Bytes()}, err
	case *bool:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: pboolType, Data: buf.Bytes()}, err
	case time.Duration:
		err := binary.Write(&buf, binary.LittleEndian, v)
		return &gobData{DataType: durationType, Data: buf.Bytes()}, err
	case *time.Duration:
		err := binary.Write(&buf, binary.LittleEndian, *v)
		return &gobData{DataType: pdurationType, Data: buf.Bytes()}, err
	case string:
		_, err := buf.WriteString(v)
		return &gobData{DataType: stringType, Data: buf.Bytes()}, err
	case *string:
		_, err := buf.WriteString(*v)
		return &gobData{DataType: pstringType, Data: buf.Bytes()}, err
	case time.Time:
		b, err := v.GobEncode()
		return &gobData{DataType: timeType, Data: b}, err
	case *time.Time:
		b, err := v.GobEncode()
		return &gobData{DataType: ptimeType, Data: b}, err
	default:
		encoder := gob.NewEncoder(&buf)
		err := encoder.Encode(data)
		return &gobData{DataType: gobType, TypeName: fmt.Sprintf("%T", data), Data: buf.Bytes()}, err
	}
}

// ErrNoGobData raised when GOB serialisation approach has no data to deserialise
var ErrNoGobData = errors.New("no data provided to deserialise")

// ErrNoDeserialisableData raised when GOB serialisation approach has value data to deserialise
var ErrNoDeserialisableData = errors.New("no data found to deserialise")

func (g *gobApproach) fromGobDataBytes(data *gobData, opts ...func(o *TypeRegistryOptions)) (any, error) {

	if data == nil {
		return nil, ErrNoGobData
	}

	if len(data.Data) == 0 {
		switch data.DataType {
		case nilType:
			return nil, nil

		default:
			return nil, ErrNoDeserialisableData
		}
	}

	if data.DataType == byteSliceType {
		return data.Data, nil
	}

	var buf = bytes.NewBuffer(data.Data)

	switch data.DataType {
	case int8Type:
		var v int8
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case pint8Type:
		v := new(int8)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case int16Type:
		var v int16
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case pint16Type:
		v := new(int16)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case int32Type:
		var v int32
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case pint32Type:
		v := new(int32)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case int64Type:
		var v int64
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case pint64Type:
		v := new(int64)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case uint8Type:
		var v uint8
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case puint8Type:
		v := new(uint8)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case uint16Type:
		var v uint16
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case puint16Type:
		v := new(uint16)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case uint32Type:
		var v uint32
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case puint32Type:
		v := new(uint32)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case uint64Type:
		var v uint64
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case puint64Type:
		v := new(uint64)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case float32Type:
		var v float32
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case pfloat32Type:
		v := new(float32)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case float64Type:
		var v float64
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case pfloat64Type:
		v := new(float64)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case boolType:
		var v bool
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case pboolType:
		v := new(bool)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case durationType:
		var v time.Duration
		err := binary.Read(buf, binary.LittleEndian, &v)
		return v, err
	case pdurationType:
		v := new(time.Duration)
		err := binary.Read(buf, binary.LittleEndian, v)
		return v, err
	case gobType:
		var buf = bytes.NewBuffer(data.Data)

		v, err := CreateInstancePtr(data.TypeName, opts...)
		if err != nil {
			return nil, err
		}

		decoder := gob.NewDecoder(buf)

		// var ss = v.([]string)
		err = decoder.Decode(v)
		if err != nil {
			return nil, err
		}

		return v, nil
	default:
		panic("Ouch!")
	}
}
