package serialise

import (
	"bytes"
	"encoding/gob"
	"errors"
)

const (
	unknownType int8 = iota
	int8Type
	pint8Type
	int8SliceType
	int16Type
	pint16Type
	int16SliceType
	int32Type
	pint32Type
	int32SliceType
	int64Type
	pint64Type
	int64SliceType
	uint8Type
	puint8Type
	uint16Type
	puint16Type
	uint16SliceType
	uint32Type
	puint32Type
	uint32SliceType
	uint64Type
	puint64Type
	uint64SliceType
	float32Type
	pfloat32Type
	float32SliceType
	float64Type
	pfloat64Type
	float64SliceType
	boolType
	pboolType
	boolSliceType
	durationType
	pdurationType
	durationSliceType
	stringType
	pstringType
	timeType
	ptimeType
	byteSliceType
	nilType
	gobType
)

// SerialisationApproach represents the available serialisation approaches
type SerialisationApproach int8

const (
	// UnknownSerialisation indicates no serialisation has been selected
	UnknownSerialisation SerialisationApproach = iota
	// GOB serialisation
	GOB
	// MinData indicates bespoke packing to minimise storage, with limited data types
	MinData
	// EndOfApproaches indicates the end of the available set of serialisation approaches
	EndOfApproaches
)

// SerialisationOptions adjust how serialisation is performed
type SerialisationOptions struct {
	// Approach specifies which serialisation method is to be used
	Approach SerialisationApproach
	// RegistryOptions allows type registry overrides to be applied
	RegistryOptions *TypeRegistryOptions
}

// WithSerialisationApproach sets the serialisation approach to be used when calling ToBytes()
func WithSerialisationApproach(approach SerialisationApproach) func(*SerialisationOptions) {
	return func(so *SerialisationOptions) {
		so.Approach = approach
	}
}

// WithTypeRegistryOptions sets the type registry options to be used when calling FromBytes()
func WithTypeRegistryOptions(opts ...func(*TypeRegistryOptions)) func(*SerialisationOptions) {

	o := &TypeRegistryOptions{}
	for _, opt := range opts {
		opt(o)
	}
	if o.Registry == nil {
		o.Registry = defaultTypeRegistry
	}

	return func(so *SerialisationOptions) {
		so.RegistryOptions = o
	}
}

// ErrUnexpectedSerialisationError raised if an invalid serialisation approach is specified
var ErrUnexpectedSerialisationError = errors.New("unexpected error during serialisation")

// GOB is the current default serialisation approach
var defaultSerialisationApproch = GOB

// ToBytes returns a byte slice of the provded data.
// Currently only gob based serialisation is available, but other options may become available, hence
// the serialisation approach used is returned to guide future deserialisation.
func ToBytes(data any, opts ...func(*SerialisationOptions)) ([]byte, SerialisationApproach, error) {

	o := SerialisationOptions{
		Approach: UnknownSerialisation,
		RegistryOptions: &TypeRegistryOptions{
			Registry: defaultTypeRegistry,
		},
	}
	for _, opt := range opts {
		opt(&o)
	}

	// Defaults to the current defaultSerialisationApproch value if not specified via opts
	if o.Approach == UnknownSerialisation {
		o.Approach = defaultSerialisationApproch
	}

	var packer func(any) ([]byte, error)

	switch o.Approach {
	case GOB:
		packer = func(data any) ([]byte, error) {
			g, err := toGobDataBytes(data)
			if err != nil {
				return nil, err
			}

			var buf bytes.Buffer
			encoder := gob.NewEncoder(&buf)
			if err := encoder.Encode(g); err != nil {
				return nil, err
			} else {
				return buf.Bytes(), nil
			}
		}
	case MinData:
		packer = toMinDataBytes
	default:
		return nil, UnknownSerialisation, ErrUnexpectedSerialisationError
	}

	b, err := packer(data)
	if err != nil {
		return nil, UnknownSerialisation, err
	}

	return b, o.Approach, err
}

// ErrNoDataToDeserialise raised if nil or empty byte slice is used in FromBytes
var ErrNoDataToDeserialise = errors.New("no data provided for deserialisation")

// ErrInvalidSerialisationApproach raised if the serialisation approach is not valid
var ErrInvalidSerialisationApproach = errors.New("invalid serialisation approach provided")

// FromBytes returns deserialises the byte slice to an instance using the specified approach.
func FromBytes(data []byte, approach SerialisationApproach, opts ...func(*SerialisationOptions)) (any, error) {

	if len(data) == 0 {
		return nil, ErrNoDataToDeserialise
	}

	if approach == UnknownSerialisation || approach >= EndOfApproaches {
		return nil, ErrInvalidSerialisationApproach
	}

	o := SerialisationOptions{
		RegistryOptions: &TypeRegistryOptions{
			Registry: defaultTypeRegistry,
		},
	}
	for _, opt := range opts {
		opt(&o)
	}

	// Note the options specified approach is ignored for deserialisation
	switch approach {
	case GOB:

		var buf = bytes.NewBuffer(data)

		decoder := gob.NewDecoder(buf)

		var g gobData

		err := decoder.Decode(&g)
		if err != nil {
			return nil, err
		}

		v, err := fromGobDataBytes(&g, func(opt *TypeRegistryOptions) { opt.replace(o.RegistryOptions) })
		if err != nil {
			return nil, err
		}
		return v, err
	case MinData:
		return fromMinDataBytes(data)
	default:
		return nil, ErrUnexpectedSerialisationError
	}
}
