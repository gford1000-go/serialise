package serialise

// import (
// 	"errors"
// 	"fmt"
// 	"reflect"
// 	"sync"
// 	"time"
// )

// // TypeRegistry manages a map of types to their type name
// type TypeRegistry struct {
// 	registry map[string]reflect.Type
// 	lck      sync.RWMutex
// }

// // AddTypeOf adds the type of the value to the registry
// func (r *TypeRegistry) AddTypeOf(v any) {
// 	r.lck.Lock()
// 	defer r.lck.Unlock()
// 	r.registry[getTypeName(v)] = reflect.TypeOf(v)
// }

// // ErrUnknownTypeName is raised if the requested type is not found in the registry
// var ErrUnknownTypeName = errors.New("requested type has not been registered")

// // GetType returns the registered type for the specified name
// func (r *TypeRegistry) GetType(name string) (reflect.Type, error) {
// 	r.lck.RLock()
// 	defer r.lck.RUnlock()

// 	if t, ok := r.registry[name]; ok {
// 		return t, nil
// 	}

// 	return nil, ErrUnknownTypeName

// }

// // NewTypeRegistry creates an instance of TypeRegistry
// func NewTypeRegistry() *TypeRegistry {
// 	return &TypeRegistry{
// 		registry: map[string]reflect.Type{},
// 	}
// }

// // TypeRegistryOptions allows overrides to the type registration behaviour
// type TypeRegistryOptions struct {
// 	Registry *TypeRegistry
// }

// func (t *TypeRegistryOptions) replace(overrides *TypeRegistryOptions) {
// 	t.Registry = overrides.Registry
// }

// func getTypeName(v any) string {
// 	return fmt.Sprintf("%T", v)
// }

// // defaultTypeRegistry is the common registry of types used by default
// var defaultTypeRegistry = NewTypeRegistry()

// // ErrCannotRegisterNilType raised if nil passed to RegisterType
// var ErrCannotRegisterNilType = errors.New("variable must not be nil in call to RegisterType")

// // ErrNoRegistry raised when there is no registry provided to operate on
// var ErrNoRegistry = errors.New("no registry provided into which to register type")

// // RegisterType allows registry of the type specified by the supplied value
// func RegisterType(v any, opts ...func(*TypeRegistryOptions)) error {
// 	if v == nil {
// 		return ErrCannotRegisterNilType
// 	}

// 	o := TypeRegistryOptions{}
// 	for _, opt := range opts {
// 		opt(&o)
// 	}

// 	if o.Registry == nil {
// 		o.Registry = defaultTypeRegistry
// 	}

// 	o.Registry.AddTypeOf(v)
// 	return nil
// }

// // GetRegisteredType returns an instance of the type specified by the name
// func GetRegisteredType(name string, opts ...func(*TypeRegistryOptions)) (reflect.Type, error) {

// 	o := TypeRegistryOptions{Registry: defaultTypeRegistry}
// 	for _, opt := range opts {
// 		opt(&o)
// 	}

// 	if o.Registry == nil {
// 		return nil, ErrNoRegistry
// 	}

// 	return o.Registry.GetType(name)
// }

// // CreateInstance returns an instance of the type specified by the name
// func CreateInstance(name string, opts ...func(*TypeRegistryOptions)) (any, error) {

// 	switch name {
// 	case "int":
// 		return int(0), nil
// 	case "int8":
// 		return int8(0), nil
// 	case "int16":
// 		return int16(0), nil
// 	case "int32":
// 		return int32(0), nil
// 	case "int64":
// 		return int64(0), nil
// 	case "*int":
// 		return new(int), nil
// 	case "*int8":
// 		return new(int8), nil
// 	case "*int16":
// 		return new(int16), nil
// 	case "*int32":
// 		return new(int32), nil
// 	case "*int64":
// 		return new(int64), nil
// 	case "uint":
// 		return uint(0), nil
// 	case "uint8":
// 		return uint8(0), nil
// 	case "uint16":
// 		return uint16(0), nil
// 	case "uint32":
// 		return uint32(0), nil
// 	case "uint64":
// 		return uint64(0), nil
// 	case "*uint":
// 		return new(uint), nil
// 	case "*uint8":
// 		return new(uint8), nil
// 	case "*uint16":
// 		return new(uint16), nil
// 	case "*uint32":
// 		return new(uint32), nil
// 	case "*uint64":
// 		return new(uint64), nil
// 	case "bool":
// 		return false, nil
// 	case "*bool":
// 		return new(bool), nil
// 	case "float32":
// 		return float32(0), nil
// 	case "*float32":
// 		return new(float32), nil
// 	case "float64":
// 		return float64(0), nil
// 	case "*float64":
// 		return new(float64), nil
// 	case "string":
// 		return "", nil
// 	case "*string":
// 		return new(string), nil
// 	case "time.Duration":
// 		return time.Duration(0), nil
// 	case "*time.Duration":
// 		return new(time.Duration), nil
// 	case "time.Time":
// 		return time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local), nil
// 	case "*time.Time":
// 		return new(time.Time), nil
// 	case "[]byte":
// 		return []byte{}, nil
// 	case "*[]byte":
// 		return new([]byte), nil
// 	default:
// 		t, err := GetRegisteredType(name, opts...)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return reflect.New(t).Elem().Interface(), nil
// 	}
// }

// func newPtr[T any]() *T {
// 	return new(T)
// }

// func newSlicePtr[T any]() *[]T {
// 	return new([]T)
// }

// // CreateInstance returns a pointer to an instance of the type specified by the name
// func CreateInstancePtr(name string, opts ...func(*TypeRegistryOptions)) (any, error) {

// 	switch name {
// 	case "int8":
// 		return newPtr[int8](), nil
// 	case "int16":
// 		return newPtr[int16](), nil
// 	case "int32":
// 		return newPtr[int32](), nil
// 	case "int64":
// 		return newPtr[int64](), nil
// 	case "uint8":
// 		return newPtr[uint8](), nil
// 	case "uint16":
// 		return newPtr[uint16](), nil
// 	case "uint32":
// 		return newPtr[uint32](), nil
// 	case "uint64":
// 		return newPtr[uint64](), nil
// 	case "bool":
// 		return newPtr[bool](), nil
// 	case "float32":
// 		return newPtr[float32](), nil
// 	case "float64":
// 		return newPtr[float64](), nil
// 	case "string":
// 		return newPtr[string](), nil
// 	case "time.Duration":
// 		return newPtr[time.Duration](), nil
// 	case "time.Time":
// 		return newPtr[time.Time](), nil
// 	case "[]byte":
// 		return newSlicePtr[byte](), nil
// 	default:
// 		t, err := GetRegisteredType(name, opts...)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return reflect.New(t).Interface(), nil
// 	}
// }
