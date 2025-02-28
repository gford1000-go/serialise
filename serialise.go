package serialise

import (
	"bytes"
	"compress/flate"
	"errors"
	"io"
)

// TypeID identifies the types that are supported by serialisation
type TypeID int8

const (
	UnknownType TypeID = iota
	Int8Type
	Pint8Type
	Int8SliceType
	Int16Type
	Pint16Type
	Int16SliceType
	Int32Type
	Pint32Type
	Int32SliceType
	Int64Type
	Pint64Type
	Int64SliceType
	Uint8Type
	Puint8Type
	Uint16Type
	Puint16Type
	Uint16SliceType
	Uint32Type
	Puint32Type
	Uint32SliceType
	Uint64Type
	Puint64Type
	Uint64SliceType
	Float32Type
	Pfloat32Type
	Float32SliceType
	Float64Type
	Pfloat64Type
	Float64SliceType
	BoolType
	PboolType
	BoolSliceType
	DurationType
	PdurationType
	DurationSliceType
	StringType
	PstringType
	StringSliceType
	TimeType
	PtimeType
	ByteSliceType
	ByteSliceSliceType
	NilType
)

// Approach implements the mechanism to be used for serialisation
type Approach interface {
	// Name of the approach
	Name() string
	// Pack serialises the instance to a byte slice
	Pack(data any) ([]byte, error)
	// Unpack deserialises an instance from the byte slice
	//Unpack(data []byte, opts ...func(opt *TypeRegistryOptions)) (any, error)
	Unpack(data []byte) (any, error)
	// IsSerialisable returns true if the type of the instance is serialisable
	IsSerialisable(v any) bool
}

// Options adjust how serialisation is performed
type Options struct {
	// Approach specifies which serialisation method is to be used
	Approach Approach
	// Encryptor will encrypt the provided data
	Encryptor func(data []byte) ([]byte, error)
	// Decryptor will decrypt the provided data
	Decryptor func(data []byte) ([]byte, error)
}

// WithSerialisationApproach sets the serialisation approach to be used when calling ToBytes()
func WithSerialisationApproach(approach Approach) func(*Options) {
	return func(so *Options) {
		so.Approach = approach
	}
}

// ErrUnexpectedSerialisationError raised if an invalid serialisation approach is specified
var ErrUnexpectedSerialisationError = errors.New("unexpected error during serialisation")

// MinData is the current default serialisation approach
var defaultSerialisationApproach = NewMinDataApproach()

// Default returns the current serialisation approach that will be used
// by Pack, Unpack etc., if not set explicitly using the WithSerialisationApproach() option.
func Default() Approach {
	return defaultSerialisationApproach
}

// ToBytes returns a byte slice of the provded data.
func ToBytes(data any, opts ...func(*Options)) ([]byte, string, error) {

	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	// Defaults to the current defaultSerialisationApproach value if not specified via opts
	if o.Approach == nil {
		o.Approach = defaultSerialisationApproach
	}

	b, err := o.Approach.Pack(data)
	if err != nil {
		return nil, "", err
	}

	b, err = deflate(b)
	if err != nil {
		return nil, "", err
	}

	// Apply optional encryption
	if o.Encryptor != nil {
		b, err = o.Encryptor(b)
		if err != nil {
			return nil, "", err
		}
	}

	return b, o.Approach.Name(), nil
}

// ErrNoDataToDeserialise raised if nil or empty byte slice is used in FromBytes
var ErrNoDataToDeserialise = errors.New("no data provided for deserialisation")

// ErrInvalidSerialisationApproach raised if the serialisation approach is not valid
var ErrInvalidSerialisationApproach = errors.New("invalid serialisation approach provided")

// FromBytes returns deserialises the byte slice to an instance using the specified approach.
func FromBytes(data []byte, approach Approach, opts ...func(*Options)) (any, error) {

	if len(data) == 0 {
		return nil, ErrNoDataToDeserialise
	}

	if approach == nil {
		return nil, ErrInvalidSerialisationApproach
	}

	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	// Apply optional encryption
	var b []byte = data
	var err error
	if o.Decryptor != nil {
		b, err = o.Decryptor(b)
		if err != nil {
			return nil, err
		}
	}

	b, err = reflate(b)
	if err != nil {
		return nil, err
	}

	return approach.Unpack(b)
}

// ToBytesMany returns a byte slice of the provded data, individually packing all the items
// into a single byte array.  Use FromBytesMany to deserialise.
// An error will be generated if any of the items in the provided data are not serialisable
// by the selected Approach.
func ToBytesMany(data []any, opts ...func(*Options)) ([]byte, string, error) {

	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	// Defaults to the current defaultSerialisationApproach value if not specified via opts
	if o.Approach == nil {
		o.Approach = defaultSerialisationApproach
	}

	output := []byte{}

	b, err := ToBytesI64(int64(len(data)))
	if err != nil {
		return nil, "", err
	}

	output = append(output, b...)

	for _, item := range data {

		b, err := o.Approach.Pack(item)
		if err != nil {
			return nil, "", err
		}

		bl, err := ToBytesI64(int64(len(b)))
		if err != nil {
			return nil, "", err
		}

		output = append(output, bl...)
		output = append(output, b...)
	}

	output, err = deflate(output)
	if err != nil {
		return nil, "", err
	}

	// Apply optional encryption
	if o.Encryptor != nil {
		output, err = o.Encryptor(output)
		if err != nil {
			return nil, "", err
		}
	}

	return output, o.Approach.Name(), nil
}

func deflate(b []byte) ([]byte, error) {
	var flag byte = 0
	if len(b) > 25 { // Trading of time cost of Flate against space... for small []byte cost is too high
		oLen := len(b)
		var buf bytes.Buffer
		writer, _ := flate.NewWriter(&buf, flate.BestCompression)
		_, err := writer.Write(b)
		if err != nil {
			return nil, err
		}
		writer.Close()
		bf := buf.Bytes()

		if oLen > len(bf) {
			flag = 1
			b = bf
		}
	}
	return append([]byte{flag}, b...), nil
}

func reflate(b []byte) ([]byte, error) {
	if b[0] == 1 {
		r := flate.NewReader(bytes.NewReader(b[1:]))
		return io.ReadAll(r)
	} else {
		return b[1:], nil
	}
}

// FromBytesMany returns deserialises the byte slice to an array of instances using the specified Approach.
func FromBytesMany(data []byte, approach Approach, opts ...func(*Options)) ([]any, error) {

	if len(data) == 0 {
		return nil, ErrNoDataToDeserialise
	}

	if approach == nil {
		return nil, ErrInvalidSerialisationApproach
	}

	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	// Apply optional encryption
	var b []byte = data
	var err error
	if o.Decryptor != nil {
		b, err = o.Decryptor(b)
		if err != nil {
			return nil, err
		}
	}

	b, err = reflate(b)
	if err != nil {
		return nil, err
	}

	var sizeI64 = SizeOfI64()

	size, err := FromBytesI64(b[0:sizeI64])
	if err != nil {
		return nil, err
	}
	b = b[sizeI64:]

	output := make([]any, size)

	var offset int64
	for offset = 0; offset < size; offset++ {
		itemSize, err := FromBytesI64(b[0:sizeI64])
		if err != nil {
			return nil, err
		}

		itemData := b[sizeI64 : sizeI64+itemSize]
		v, err := approach.Unpack(itemData)
		if err != nil {
			return nil, err
		}

		output[offset] = v
		b = b[sizeI64+itemSize:]
	}

	return output, nil
}
