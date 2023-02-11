# Uuid25: 25-digit case-insensitive UUID encoding

[![GitHub tag](https://img.shields.io/github/v/tag/uuid25/go-uuid25)](https://github.com/uuid25/go-uuid25)
[![License](https://img.shields.io/github/license/uuid25/go-uuid25)](https://github.com/uuid25/go-uuid25/blob/main/LICENSE)

Uuid25 is an alternative UUID representation that shortens a UUID string to just
25 digits using the case-insensitive Base36 encoding. This library provides
functionality to convert from the conventional UUID formats to Uuid25 and vice
versa.

```go
import "github.com/uuid25/go-uuid25"

// convert from/to string
a, _ := uuid25.Parse("8da942a4-1fbe-4ca6-852c-95c473229c7d")
assert(a.String() == "8dx554y5rzerz1syhqsvsdw8t")
assert(a.ToHyphenated() == "8da942a4-1fbe-4ca6-852c-95c473229c7d")

// convert from/to 128-bit byte array
b := uuid25.FromBytes([]byte{
  0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
  0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
assert(b.String() == "f5lxx1zz5pnorynqglhzmsp33")
for _, x := range b.ToBytes() {
  assert(x == 0xff)
}

// convert from/to other popular textual representations
c := make([]uuid25.Uuid25, 4)
c[0], _ = uuid25.Parse("e7a1d63b711744238988afcf12161878")
c[1], _ = uuid25.Parse("e7a1d63b-7117-4423-8988-afcf12161878")
c[2], _ = uuid25.Parse("{e7a1d63b-7117-4423-8988-afcf12161878}")
c[3], _ = uuid25.Parse("urn:uuid:e7a1d63b-7117-4423-8988-afcf12161878")
for _, x := range c {
  assert(x.String() == "dpoadk8izg9y4tte7vy1xt94o")
}

d, _ := uuid25.Parse("dpoadk8izg9y4tte7vy1xt94o")
assert(d.ToHex() == "e7a1d63b711744238988afcf12161878")
assert(d.ToHyphenated() == "e7a1d63b-7117-4423-8988-afcf12161878")
assert(d.ToBraced() == "{e7a1d63b-7117-4423-8988-afcf12161878}")
assert(d.ToUrn() == "urn:uuid:e7a1d63b-7117-4423-8988-afcf12161878")

func assert(c bool) { if !c { panic("assertion failed") } }
```

## License

Licensed under the Apache License, Version 2.0.

## See also

- [uuid25 package - github.com/uuid25/go-uuid25 - pkg.go.dev](https://pkg.go.dev/github.com/uuid25/go-uuid25)
