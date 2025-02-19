package serialise

import (
	"errors"
)

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
	GobType
)

// Approach implements the mechanism to be used for serialisation
type Approach interface {
	// Name of the approach
	Name() string
	// Pack serialises the instance to a byte slice
	Pack(data any) ([]byte, error)
	// Unpack deserialises an instance from the byte slice
	Unpack(data []byte, opts ...func(opt *TypeRegistryOptions)) (any, error)
	// IsSerialisable returns true if the type of the instance is serialisable
	IsSerialisable(v any) bool
}

// SerialisationOptions adjust how serialisation is performed
type SerialisationOptions struct {
	// Approach specifies which serialisation method is to be used
	Approach Approach
	// RegistryOptions allows type registry overrides to be applied
	RegistryOptions *TypeRegistryOptions
	// Encryptor will encrypt the provided data
	Encryptor func(data []byte) ([]byte, error)
	// Decryptor will decrypt the provided data
	Decryptor func(data []byte) ([]byte, error)
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

	// Apply optional encryption
	var b []byte = data
	var err error
	if o.Decryptor != nil {
		b, err = o.Decryptor(b)
		if err != nil {
			return nil, err
		}
	}

	return approach.Unpack(b, func(opt *TypeRegistryOptions) { opt.replace(o.RegistryOptions) })
}
