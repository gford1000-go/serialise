package serialise

import (
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
	stringSliceType
	timeType
	ptimeType
	byteSliceType
	byteSliceSliceType
	nilType
	gobType
)

// Approach implements the mechanism to be used for serialisation
type Approach interface {
	// Name of the approach
	Name() string
	// Pack serialises the instance to a byte slice
	Pack(data any) ([]byte, error)
	// Unpack deserialises an instance from the byte slice
	Unpack(data []byte, opts ...func(opt *TypeRegistryOptions)) (any, error)
}

// SerialisationOptions adjust how serialisation is performed
type SerialisationOptions struct {
	// Approach specifies which serialisation method is to be used
	Approach Approach
	// RegistryOptions allows type registry overrides to be applied
	RegistryOptions *TypeRegistryOptions
}

// WithSerialisationApproach sets the serialisation approach to be used when calling ToBytes()
func WithSerialisationApproach(approach Approach) func(*SerialisationOptions) {
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

// MinData is the current default serialisation approach
var defaultSerialisationApproach = NewMinDataApproach()

// ToBytes returns a byte slice of the provded data.
// Currently only gob based serialisation is available, but other options may become available, hence
// the serialisation approach used is returned to guide future deserialisation.
func ToBytes(data any, opts ...func(*SerialisationOptions)) ([]byte, string, error) {

	o := SerialisationOptions{
		Approach: nil,
		RegistryOptions: &TypeRegistryOptions{
			Registry: defaultTypeRegistry,
		},
	}
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
	return b, o.Approach.Name(), nil
}

// ErrNoDataToDeserialise raised if nil or empty byte slice is used in FromBytes
var ErrNoDataToDeserialise = errors.New("no data provided for deserialisation")

// ErrInvalidSerialisationApproach raised if the serialisation approach is not valid
var ErrInvalidSerialisationApproach = errors.New("invalid serialisation approach provided")

// FromBytes returns deserialises the byte slice to an instance using the specified approach.
func FromBytes(data []byte, approach Approach, opts ...func(*SerialisationOptions)) (any, error) {

	if len(data) == 0 {
		return nil, ErrNoDataToDeserialise
	}

	if approach == nil {
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

	return approach.Unpack(data, func(opt *TypeRegistryOptions) { opt.replace(o.RegistryOptions) })
}
