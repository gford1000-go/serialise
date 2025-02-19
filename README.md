[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://en.wikipedia.org/wiki/MIT_License)
[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-green.svg)](https://godoc.org/github.com/gford1000-go/serialise)

# Serialise

Serialises basic types, pointers to types, and slices to a byte slice.

Allows the serialisation approach to be extended via the `Approach` interface.

Supports optional encryption of the byte slice using `aes-gcm`.

```go
func main() {
    data := []string{"Hello", "World!"}

    // Serialise the data using the default version of MinData serialisation
    b, name, _ := ToBytes(data, WithSerialisationApproach(NewMinDataApproach()))

    // Retrieve the Approach used for serialisation, from the returned name
    approach, _ := GetApproach(name)

    // Deserialise
    v, _ := FromBytes(b, approach)
}
```
