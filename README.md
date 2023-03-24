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

The [uuid25ext] package integrates the popular [github.com/google/uuid] module
and adds functionality to generate a UUID value in the Uuid25 format.

```go
import "fmt"
import "github.com/google/uuid"
import "github.com/uuid25/go-uuid25/ext"

// convert from/to github.com/google/uuid module's UUID value
googleUuid, _ := uuid.Parse("f38a6b1f-576f-4c22-8d4a-5f72613483f6")
e := uuid25ext.FromUUID(googleUuid)
assert(e == "ef1zh7jc64vprqez41vbwe9km")
assert(uuid25ext.ToUUID(e) == googleUuid)

// generate new UUID in Uuid25 format
fmt.Println(uuid25ext.NewV4()) // e.g. "99wfqtl0z0yevxzpl4hv2dm5p"

func assert(c bool) { if !c { panic("assertion failed") } }
```

[uuid25ext]: https://pkg.go.dev/github.com/uuid25/go-uuid25/ext
[github.com/google/uuid]: https://pkg.go.dev/github.com/google/uuid

## License

Licensed under the Apache License, Version 2.0.

## See also

- [uuid25 package - github.com/uuid25/go-uuid25 - Go Packages](https://pkg.go.dev/github.com/uuid25/go-uuid25)
- [uuid25ext package - github.com/uuid25/go-uuid25/ext - Go Packages](https://pkg.go.dev/github.com/uuid25/go-uuid25/ext)
