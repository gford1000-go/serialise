package serialise

import (
	"errors"
	"sync"
)

type serialisationRegistry struct {
	m   map[string]Approach
	lck sync.RWMutex
}

var registry *serialisationRegistry

func init() {
	registry = &serialisationRegistry{
		m: map[string]Approach{},
	}

	RegisterApproach(NewGOBApproach())
	RegisterApproach(NewMinDataApproach())
}

// RegisterApproach allows all registered Approach to be retrievable by Name()
func RegisterApproach(a Approach) {
	if a == nil {
		return
	}

	registry.lck.Lock()
	defer registry.lck.Unlock()

	registry.m[a.Name()] = a
}

// ErrUnknownApproach raised if the specified name is not found in the Approach registry
var ErrUnknownApproach = errors.New("specified Approach name is not registered")

// GetApproach returns the Approach with the specified name
func GetApproach(name string) (Approach, error) {
	registry.lck.RLock()
	defer registry.lck.RUnlock()

	a, ok := registry.m[name]
	if !ok {
		return nil, ErrUnknownApproach
	}
	return a, nil
}
