// Uuid25: 25-digit case-insensitive UUID encoding
//
// Uuid25 is an alternative UUID representation that shortens a UUID string to
// just 25 digits using the case-insensitive Base36 encoding. This library
// provides functionality to convert from the conventional UUID formats to
// Uuid25 and vice versa.
package uuid25

import (
	"database/sql/driver"
	"errors"
	"math"
)

// The primary value type containing the Uuid25 representation of a UUID.
//
// A valid value of this type must be constructed through FromBytes() or one of
// Parse*() functions.
type Uuid25 string

// Returns the 25-digit Uuid25 representation of this type.
func (uuid25 Uuid25) String() string {
	if len(uuid25) != 25 {
		// conduct O(1) quick check here because all other value receiver methods
		// directly or indirectly call String()
		panic("receiver not constructed properly")
	}
	return string(uuid25)
}

// Creates an instance from an array of Base36 digit values.
func fromDigitValues(digitValues []byte) (Uuid25, error) {
	if len(digitValues) != 25 {
		panic("invalid length of digit value array")
	}
	const digits = "0123456789abcdefghijklmnopqrstuvwxyz"
	const u128Max = "f5lxx1zz5pnorynqglhzmsp33" // 2^128 - 1

	var buffer [25]byte
	maybeTooLarge := true
	for i, e := range digitValues {
		if e >= 36 {
			return "", parseError // invalid digit value
		}
		buffer[i] = digits[e]
		if maybeTooLarge && buffer[i] > u128Max[i] {
			return "", parseError // 128-bit overflow
		} else if buffer[i] < u128Max[i] {
			maybeTooLarge = false
		}
	}
	return Uuid25(buffer[:]), nil
}

// Creates an instance from a 16-byte UUID binary representation.
func FromBytes(uuidBytes []byte) Uuid25 {
	if len(uuidBytes) != 16 {
		panic("the length of byte slice must be 16")
	}
	var buffer [25]byte
	if convertBase(uuidBytes[:], buffer[:], 256, 36) == nil {
		if uuid25, err := fromDigitValues(buffer[:]); err == nil {
			return uuid25
		}
	}
	panic("unreachable")
}

// Converts this type into the 16-byte binary representation of a UUID.
func (uuid25 Uuid25) ToBytes() [16]byte {
	var src [25]byte
	if decodeDigitChars(uuid25.String(), src[:], 36) == nil {
		var uuidBytes [16]byte
		if convertBase(src[:], uuidBytes[:], 36, 256) == nil {
			return uuidBytes
		}
	}
	panic("unreachable")
}

// Creates an instance from a UUID string representation.
//
// This method accepts the following formats:
//
//   - 25-digit Base36 Uuid25 format: `3ud3gtvgolimgu9lah6aie99o`
//   - 32-digit hexadecimal format without hyphens:
//     `40eb9860cf3e45e2a90eb82236ac806c`
//   - 8-4-4-4-12 hyphenated format: `40eb9860-cf3e-45e2-a90e-b82236ac806c`
//   - Hyphenated format with surrounding braces:
//     `{40eb9860-cf3e-45e2-a90e-b82236ac806c}`
//   - RFC 4122 URN format: `urn:uuid:40eb9860-cf3e-45e2-a90e-b82236ac806c`
func Parse(uuidString string) (Uuid25, error) {
	switch len(uuidString) {
	case 25:
		return ParseUuid25(uuidString)
	case 32:
		return ParseHex(uuidString)
	case 36:
		return ParseHyphenated(uuidString)
	case 38:
		return ParseBraced(uuidString)
	case 45:
		return ParseUrn(uuidString)
	default:
		return "", parseError
	}
}

// Creates an instance from the 25-digit Base36 Uuid25 format:
// `3ud3gtvgolimgu9lah6aie99o`.
func ParseUuid25(uuidString string) (Uuid25, error) {
	if len(uuidString) != 25 {
		return "", parseError
	}
	var buffer [25]byte
	if err := decodeDigitChars(uuidString, buffer[:], 36); err != nil {
		return "", parseError
	}
	return fromDigitValues(buffer[:])
}

// Creates an instance from the 32-digit hexadecimal format without hyphens:
// `40eb9860cf3e45e2a90eb82236ac806c`.
func ParseHex(uuidString string) (Uuid25, error) {
	if len(uuidString) != 32 {
		return "", parseError
	}
	var src [32]byte
	if err := decodeDigitChars(uuidString, src[:], 16); err != nil {
		return "", parseError
	}
	var buffer [25]byte
	if err := convertBase(src[:], buffer[:], 16, 36); err != nil {
		return "", parseError
	}
	return fromDigitValues(buffer[:])
}

// Creates an instance from the 8-4-4-4-12 hyphenated format:
// `40eb9860-cf3e-45e2-a90e-b82236ac806c`.
func ParseHyphenated(uuidString string) (Uuid25, error) {
	if len(uuidString) != 36 ||
		uuidString[8] != '-' ||
		uuidString[13] != '-' ||
		uuidString[18] != '-' ||
		uuidString[23] != '-' {
		return "", parseError
	}
	return ParseHex(
		uuidString[:8] +
			uuidString[9:13] +
			uuidString[14:18] +
			uuidString[19:23] +
			uuidString[24:])
}

// Creates an instance from the hyphenated format with surrounding braces:
// `{40eb9860-cf3e-45e2-a90e-b82236ac806c}`.
func ParseBraced(uuidString string) (Uuid25, error) {
	if len(uuidString) != 38 ||
		uuidString[0] != '{' ||
		uuidString[37] != '}' {
		return "", parseError
	}
	return ParseHyphenated(uuidString[1:37])
}

// Creates an instance from the RFC 4122 URN format:
// `urn:uuid:40eb9860-cf3e-45e2-a90e-b82236ac806c`.
func ParseUrn(uuidString string) (Uuid25, error) {
	if len(uuidString) != 45 ||
		(uuidString[0] != 'U' && uuidString[0] != 'u') ||
		(uuidString[1] != 'R' && uuidString[1] != 'r') ||
		(uuidString[2] != 'N' && uuidString[2] != 'n') ||
		(uuidString[3] != ':') ||
		(uuidString[4] != 'U' && uuidString[4] != 'u') ||
		(uuidString[5] != 'U' && uuidString[5] != 'u') ||
		(uuidString[6] != 'I' && uuidString[6] != 'i') ||
		(uuidString[7] != 'D' && uuidString[7] != 'd') ||
		(uuidString[8] != ':') {
		return "", parseError
	}
	return ParseHyphenated(uuidString[9:])
}

// Formats this type in the 32-digit hexadecimal format without hyphens:
// `40eb9860cf3e45e2a90eb82236ac806c`.
func (uuid25 Uuid25) ToHex() string {
	const digits = "0123456789abcdef"
	var src [25]byte
	if decodeDigitChars(uuid25.String(), src[:], 36) != nil {
		panic("unreachable")
	}
	var buffer [32]byte
	if convertBase(src[:], buffer[:], 36, 16) != nil {
		panic("unreachable")
	}
	for i, e := range buffer {
		buffer[i] = digits[e]
	}
	return string(buffer[:])
}

// Formats this type in the 8-4-4-4-12 hyphenated format:
// `40eb9860-cf3e-45e2-a90e-b82236ac806c`.
func (uuid25 Uuid25) ToHyphenated() string {
	s := uuid25.ToHex()
	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

// Formats this type in the hyphenated format with surrounding braces:
// `{40eb9860-cf3e-45e2-a90e-b82236ac806c}`.
func (uuid25 Uuid25) ToBraced() string {
	return "{" + uuid25.ToHyphenated() + "}"
}

// Formats this type in the RFC 4122 URN format:
// `urn:uuid:40eb9860-cf3e-45e2-a90e-b82236ac806c`.
func (uuid25 Uuid25) ToUrn() string {
	return "urn:uuid:" + uuid25.ToHyphenated()
}

// Implements the encoding.TextUnmarshaler interface.
func (uuid25 *Uuid25) UnmarshalText(text []byte) error {
	if uuid25 == nil {
		return errors.New("nil receiver")
	}
	result, err := Parse(string(text))
	*uuid25 = result
	return err
}

// Implements the encoding.TextMarshaler interface.
func (uuid25 Uuid25) MarshalText() (text []byte, err error) {
	return []byte(uuid25.String()), nil
}

// Implements the encoding.BinaryUnmarshaler interface.
func (uuid25 *Uuid25) UnmarshalBinary(data []byte) error {
	if uuid25 == nil {
		return errors.New("nil receiver")
	} else if len(data) == 16 {
		*uuid25 = FromBytes(data)
		return nil
	}
	return uuid25.UnmarshalText(data)
}

// Implements the encoding.BinaryMarshaler interface.
func (uuid25 Uuid25) MarshalBinary() (data []byte, err error) {
	return uuid25.MarshalText()
}

// Implements the sql.Scanner interface.
func (uuid25 *Uuid25) Scan(src any) error {
	if uuid25 == nil {
		return errors.New("nil receiver")
	}
	switch src := src.(type) {
	case string:
		return uuid25.UnmarshalText([]byte(src))
	case []byte:
		return uuid25.UnmarshalBinary(src)
	default:
		return errors.New("unsupported type conversion")
	}
}

// Implements the driver.Valuer interface.
func (uuid25 Uuid25) Value() (driver.Value, error) {
	return uuid25.String(), nil
}

// An error parsing a UUID string representation.
var parseError = errors.New("could not parse a UUID string")

// Converts a digit value array in `srcBase` to that in `dstBase`.
func convertBase(src []byte, dst []byte, srcBase uint, dstBase uint) error {
	if srcBase < 2 || srcBase > 256 || dstBase < 2 || dstBase > 256 {
		panic("invalid base")
	}

	// determine the number of `src` digits to read for each outer loop
	wordLen := 1
	wordBase := srcBase
	for wordBase <= math.MaxUint/(srcBase*dstBase) {
		wordLen += 1
		wordBase *= srcBase
	}

	for i := range dst {
		dst[i] = 0
	}

	if len(src) == 0 {
		return nil
	} else if len(dst) == 0 {
		return errors.New("too small dst")
	}

	dstUsed := len(dst) - 1 // storage to memorize range of `dst` filled

	// read `wordLen` digits from `src` for each outer loop
	wordHead := len(src) % wordLen
	if wordHead > 0 {
		wordHead -= wordLen
	}
	for ; wordHead < len(src); wordHead += wordLen {
		var carry uint = 0
		i := wordHead
		if i < 0 {
			i = 0
		}
		for ; i < wordHead+wordLen; i += 1 {
			e := uint(src[i])
			if e >= srcBase {
				panic("invalid src digit")
			}
			carry = carry*srcBase + e
		}

		// fill in `dst` from right to left, while carrying up prior result to left
		for i := len(dst) - 1; i >= 0; i -= 1 {
			carry += uint(dst[i]) * wordBase
			dst[i] = byte(carry % dstBase)
			carry /= dstBase

			// break inner loop when `carry` and remaining `dst` digits are all zero
			if carry == 0 && i <= dstUsed {
				dstUsed = i
				break
			}
		}
		if carry > 0 {
			return errors.New("too small dst")
		}
	}
	return nil
}

// An O(1) map from ASCII code points to Base36 digit values.
var decodeMap = [256]byte{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x01, 0x02, 0x03,
	0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16,
	0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d,
	0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
}

// Converts from a string of digit characters to an array of digit values.
func decodeDigitChars(src string, dst []byte, base byte) error {
	if base < 2 || base > 36 {
		panic("invalid base")
	}
	if len(src) != len(dst) {
		panic("invalid length of dst slice")
	}
	for i := 0; i < len(src); i += 1 {
		dst[i] = decodeMap[src[i]]
		if dst[i] >= base {
			return errors.New("invalid digit character")
		}
	}
	return nil
}
