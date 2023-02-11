// Extension to the uuid25 package that integrates third party modules
package uuid25ext

import (
	"github.com/google/uuid"
	"github.com/uuid25/go-uuid25"
)

// Creates a Uuid25 value from the UUID type of github.com/google/uuid.
func FromUUID(uuidValue uuid.UUID) uuid25.Uuid25 {
	return uuid25.FromBytes(uuidValue[:])
}

// Converts a Uuid25 value into the UUID type of github.com/google/uuid.
func ToUUID(uuid25 uuid25.Uuid25) uuid.UUID {
	uuid25Bytes := uuid25.ToBytes()
	uuidValue, err := uuid.FromBytes(uuid25Bytes[:])
	if err != nil {
		panic("unreachable")
	}
	return uuidValue
}

// Generates a random UUID (UUIDv4) value encoded in the Uuid25 format.
func NewV4() uuid25.Uuid25 {
	return FromUUID(uuid.New())
}

// Equivalent to [uuid25.FromBytes], re-exported for convenience.
func FromBytes(uuidBytes []byte) uuid25.Uuid25 {
	return uuid25.FromBytes(uuidBytes)
}

// Equivalent to [uuid25.Parse], re-exported for convenience.
func Parse(uuidString string) (uuid25.Uuid25, error) {
	return uuid25.Parse(uuidString)
}
