package serialise

import (
	"fmt"
	"strings"
)

func Example() {

	data := []string{"Hello", "World!"}

	key := []byte("01234567890123456789012345678912")

	// Serialise the data using the default version of MinData serialisation,
	// encrypting the serialised byte slice
	b, name, _ := ToBytes(data, WithSerialisationApproach(NewMinDataApproach()), WithAESGCMEncryption(key))

	// Retrieve the Approach used for serialisation, from the returned name
	approach, _ := GetApproach(name)

	// Decrypt and deserialise
	v, _ := FromBytes(b, approach, WithAESGCMEncryption(key))

	fmt.Println(strings.Join(data, " ") == strings.Join(v.([]string), " "))
	// Output: true
}
