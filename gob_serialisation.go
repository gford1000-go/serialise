package serialise

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"encoding/gob"
// 	"errors"
// 	"fmt"
// 	"time"
// )

// // NewGOBApproach creates an Approach instance
// // that uses gob serialisation.
// func NewGOBApproach() Approach {
// 	return &gobApproach{}
// }

// type gobApproach struct {
// }

// // Name of the approach
// func (g *gobApproach) Name() string {
// 	return "GOB"
// }

// // IsSerialisable returns true if an instance of the specified type
// // can be serialised
// func (g *gobApproach) IsSerialisable(v any) (ok bool) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			ok = false
// 		}
// 	}()

// 	_, err := g.Pack(v)
// 	return err == nil
// }

// // Pack serialises the instance to a byte slice
// func (g *gobApproach) Pack(data any) ([]byte, error) {
// 	gd, err := g.toGobDataBytes(data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var buf bytes.Buffer
// 	encoder := gob.NewEncoder(&buf)
// 	if err := encoder.Encode(gd); err != nil {
// 		return nil, err
// 	} else {
// 		return buf.Bytes(), nil
// 	}
// }

// // Unpack deserialises an instance from the byte slice
// func (g *gobApproach) Unpack(data []byte, opts ...func(opt *TypeRegistryOptions)) (any, error) {
// 	var buf = bytes.NewBuffer(data)

// 	decoder := gob.NewDecoder(buf)

// 	var gd gobData

// 	err := decoder.Decode(&gd)
// 	if err != nil {
// 		return nil, err
// 	}

// 	v, err := g.fromGobDataBytes(&gd, opts...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return v, err

// }

// type gobData struct {
// 	DataType TypeID
// 	TypeName string
// 	Data     []byte
// }

// // toGobDataBytes serialises data types to []byte using gob encoding
// func (g *gobApproach) toGobDataBytes(data any) (*gobData, error) {
// 	if data == nil {
// 		return &gobData{DataType: NilType, Data: []byte{}}, nil
// 	}

// 	var buf bytes.Buffer

// 	switch v := data.(type) {
// 	case []byte:
// 		return &gobData{DataType: ByteSliceType, Data: v}, nil
// 	case int8:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: Int8Type, Data: buf.Bytes()}, err
// 	case *int8:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: Pint8Type, Data: buf.Bytes()}, err
// 	case int16:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: Int16Type, Data: buf.Bytes()}, err
// 	case *int16:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: Pint16Type, Data: buf.Bytes()}, err
// 	case int32:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: Int32Type, Data: buf.Bytes()}, err
// 	case *int32:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: Pint32Type, Data: buf.Bytes()}, err
// 	case int64:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: Int64Type, Data: buf.Bytes()}, err
// 	case *int64:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: Pint64Type, Data: buf.Bytes()}, err
// 	case uint8:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: Uint8Type, Data: buf.Bytes()}, err
// 	case *uint8:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: Puint8Type, Data: buf.Bytes()}, err
// 	case uint16:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: Uint16Type, Data: buf.Bytes()}, err
// 	case *uint16:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: Puint16Type, Data: buf.Bytes()}, err
// 	case uint32:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: Uint32Type, Data: buf.Bytes()}, err
// 	case *uint32:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: Puint32Type, Data: buf.Bytes()}, err
// 	case uint64:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: Uint64Type, Data: buf.Bytes()}, err
// 	case *uint64:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: Puint64Type, Data: buf.Bytes()}, err
// 	case float32:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: Float32Type, Data: buf.Bytes()}, err
// 	case *float32:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: Pfloat32Type, Data: buf.Bytes()}, err
// 	case float64:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: Float64Type, Data: buf.Bytes()}, err
// 	case *float64:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: Pfloat64Type, Data: buf.Bytes()}, err
// 	case bool:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: BoolType, Data: buf.Bytes()}, err
// 	case *bool:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: PboolType, Data: buf.Bytes()}, err
// 	case time.Duration:
// 		err := binary.Write(&buf, binary.LittleEndian, v)
// 		return &gobData{DataType: DurationType, Data: buf.Bytes()}, err
// 	case *time.Duration:
// 		err := binary.Write(&buf, binary.LittleEndian, *v)
// 		return &gobData{DataType: PdurationType, Data: buf.Bytes()}, err
// 	case string:
// 		_, err := buf.WriteString(v)
// 		return &gobData{DataType: StringType, Data: buf.Bytes()}, err
// 	case *string:
// 		_, err := buf.WriteString(*v)
// 		return &gobData{DataType: PstringType, Data: buf.Bytes()}, err
// 	case time.Time:
// 		b, err := v.GobEncode()
// 		return &gobData{DataType: TimeType, Data: b}, err
// 	case *time.Time:
// 		b, err := v.GobEncode()
// 		return &gobData{DataType: PtimeType, Data: b}, err
// 	default:
// 		encoder := gob.NewEncoder(&buf)
// 		err := encoder.Encode(data)
// 		return &gobData{DataType: GobType, TypeName: fmt.Sprintf("%T", data), Data: buf.Bytes()}, err
// 	}
// }

// // ErrNoGobData raised when GOB serialisation approach has no data to deserialise
// var ErrNoGobData = errors.New("no data provided to deserialise")

// // ErrNoDeserialisableData raised when GOB serialisation approach has value data to deserialise
// var ErrNoDeserialisableData = errors.New("no data found to deserialise")

// func (g *gobApproach) fromGobDataBytes(data *gobData, opts ...func(o *TypeRegistryOptions)) (any, error) {

// 	if data == nil {
// 		return nil, ErrNoGobData
// 	}

// 	if len(data.Data) == 0 {
// 		switch data.DataType {
// 		case NilType:
// 			return nil, nil

// 		default:
// 			return nil, ErrNoDeserialisableData
// 		}
// 	}

// 	if data.DataType == ByteSliceType {
// 		return data.Data, nil
// 	}

// 	var buf = bytes.NewBuffer(data.Data)

// 	switch data.DataType {
// 	case Int8Type:
// 		var v int8
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case Pint8Type:
// 		v := new(int8)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case Int16Type:
// 		var v int16
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case Pint16Type:
// 		v := new(int16)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case Int32Type:
// 		var v int32
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case Pint32Type:
// 		v := new(int32)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case Int64Type:
// 		var v int64
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case Pint64Type:
// 		v := new(int64)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case Uint8Type:
// 		var v uint8
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case Puint8Type:
// 		v := new(uint8)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case Uint16Type:
// 		var v uint16
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case Puint16Type:
// 		v := new(uint16)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case Uint32Type:
// 		var v uint32
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case Puint32Type:
// 		v := new(uint32)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case Uint64Type:
// 		var v uint64
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case Puint64Type:
// 		v := new(uint64)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case Float32Type:
// 		var v float32
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case Pfloat32Type:
// 		v := new(float32)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case Float64Type:
// 		var v float64
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case Pfloat64Type:
// 		v := new(float64)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case BoolType:
// 		var v bool
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case PboolType:
// 		v := new(bool)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case DurationType:
// 		var v time.Duration
// 		err := binary.Read(buf, binary.LittleEndian, &v)
// 		return v, err
// 	case PdurationType:
// 		v := new(time.Duration)
// 		err := binary.Read(buf, binary.LittleEndian, v)
// 		return v, err
// 	case GobType:
// 		var buf = bytes.NewBuffer(data.Data)

// 		v, err := CreateInstancePtr(data.TypeName, opts...)
// 		if err != nil {
// 			return nil, err
// 		}

// 		decoder := gob.NewDecoder(buf)

// 		// var ss = v.([]string)
// 		err = decoder.Decode(v)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return v, nil
// 	default:
// 		panic("Ouch!")
// 	}
// }
