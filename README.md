# smskit-go

`github.com/Flying-Tea-Squad/smskit-go` is an open-source Go SDK for sending
SMS and WhatsApp messages through provider adapters behind a small,
provider-independent API.

## Status

The project is currently pre-v1. The public core is being established before
provider adapters are added. The API may change between pre-v1 releases.

## Requirements

- Go 1.22 or newer

The `go.mod` `go` directive is the minimum supported Go language and standard
library baseline. Newer Go releases are supported unless a release note says
otherwise.

## Package layout

The root package is named `smskit` and contains only shared messaging
contracts: capability interfaces, provider-independent message types, common
results, error categories, and reusable client options.

Each provider is an independent top-level package with a lowercase name, such
as `africastalking` or `safravo`. A provider may import the root `smskit`
package and approved packages under `internal/`, but provider packages must
not import one another. Provider authentication, wire payloads, statuses, and
webhook types stay in the provider package that owns them.

See [`plan.md`](plan.md) for the broader architecture and
[`docs/design-decisions.md`](docs/design-decisions.md) for the module,
package, and compatibility decisions.

## Checks

```sh
gofmt -l .
go test ./...
go vet ./...
```

## License

This project is licensed under the MIT License. See [`LICENSE`](LICENSE).
