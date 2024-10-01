# Concise Reference Integrity Manifest and Module Identifiers
[![ci](https://github.com/jraman567/corim/actions/workflows/ci.yml/badge.svg)](https://github.com/jraman567/corim/actions/workflows/ci.yml)
[![cover ≥82%](https://github.com/jraman567/corim/actions/workflows/ci-go-cover.yml/badge.svg)](https://github.com/jraman567/corim/actions/workflows/ci-go-cover.yml)
[![linters](https://github.com/jraman567/corim/actions/workflows/linters.yml/badge.svg)](https://github.com/jraman567/corim/actions/workflows/linters.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jraman567/corim.svg)](https://pkg.go.dev/github.com/jraman567/corim)

The [`corim/corim`](corim) and [`corim/comid`](comid) packages provide a golang API for low-level manipulation of [Concise Reference Integrity Manifest (CoRIM)](https://datatracker.ietf.org/doc/draft-birkholz-rats-corim/) and Concise Module Identifier (CoMID) tags respectively.

> [!NOTE]
> These API are still in active development (as is the underlying CoRIM spec).
> They are **subject to change** in the future.

## Developer tips

Before requesting a PR (and routinely during the dev/test cycle), you are encouraged to run:
```
make presubmit
```
and check its output to make sure your code coverage figures are in line with the set target and that there are no newly introduced lint problems.

## Extending CoRIM/CoMID

The CoRIM specification provides a mechanism for adding extensions to the base
CoRIM schema. The `corim` and `comid` structs which can be extended, embed an
`Extensions` object  that allows registering a wrapper structure defining
extension fields. For field types that can be extended, i.e. `type choice`,
extensions can be implemented by calling an appropriate registration function
and giving it a new type or a value (for enums).

Please see [extensions documentation](extensions/README.md) for details.


