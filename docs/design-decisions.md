# Design decisions

## Module and public package boundary

The module path is:

```text
github.com/Flying-Tea-Squad/smskit-go
```

The root package name is `smskit`. It is the only shared public package and
contains provider-independent messaging contracts: focused capability
interfaces, common message and result types, normalized error categories, and
shared options.

Provider adapters are independent top-level packages with lowercase names.
They may import `smskit` and shared implementation details under `internal/`,
but they must not import another provider package. Authentication, provider
wire structures, provider-specific statuses, and webhook payloads remain
inside their owning provider package.

## Go support

Go 1.22 is the minimum supported version, as recorded by the `go` directive
in `go.mod`. Go 1.22 and later releases are supported unless explicitly
excluded in release notes. CI and local development should use the latest
available patch release for the selected Go version.

## Compatibility policy

The project is pre-v1, so minor and patch releases may contain breaking API
changes while the core contracts are being validated. Such changes must be
documented in release notes when they affect provider authors or applications.

Starting at v1.0.0, the project will follow semantic versioning: exported API
changes that break source compatibility require a new major version; additive
API changes use a minor version; and compatible fixes use a patch version.
Provider behavior and external provider APIs remain subject to each provider's
documented service contract.

## License

The repository is distributed under the MIT License. The complete license
text is in [`LICENSE`](../LICENSE).
