# Unified Go Messaging SDK — Implementation Plan

## 1. Purpose

This project will provide an open-source Go library for sending messages through multiple messaging providers behind a small, consistent API.

The initial release will support four providers:

- **Africa's Talking** — the first SMS implementation and the reference adapter.
- **SMSLeopard** — a second local SMS implementation used to validate that the shared API is genuinely provider-independent.
- **Wasiliana** — a Kenyan SMS and USSD provider, documented at https://docs.wasiliana.com/, used to validate the adapter pattern against a `apiKey`-header based API ecosystem.
- **Safravo** — a Kenyan omnichannel provider for SMS plus the first official WhatsApp Business implementation, including text, templates, media, and webhook events.

These are starter providers, not the final provider list. Their purpose is to establish a stable core library, provider structure, test approach, and contribution process. Once those foundations are proven, maintainers and community contributors can add more providers without redesigning the SDK.

The project is a Go SDK, not a hosted messaging service, chatbot engine, or multi-tenant WhatsApp gateway. Applications using the SDK remain responsible for their own business logic, HTTP servers, databases, queues, and deployments.

## 2. Goals

The SDK should:

- Give Go applications a simple way to send SMS and WhatsApp messages.
- Allow common application code to depend on capabilities instead of a specific provider.
- Preserve access to provider-specific features when the common API is not sufficient.
- Normalize useful fields such as message IDs and delivery states without hiding the original provider response.
- Provide safe defaults for HTTP timeouts, error handling, logging, retries, and secret redaction.
- Make each provider package independently testable.
- Make adding a provider understandable and repeatable for new contributors.
- Keep the core package small and avoid unnecessary dependencies.

## 3. Non-goals for the initial release

The first release will not attempt to provide:

- A visual chatbot builder or conversation engine.
- Contact lists, marketing campaigns, or audience management.
- Message scheduling or a persistent job queue.
- A hosted webhook relay.
- Automatic failover between providers.
- One universal structure containing every feature from every provider.
- An unofficial WhatsApp Web integration.
- Guaranteed normalization of provider pricing or delivery terminology.

These features can be considered later, but they must not complicate the core provider SDK before the first four integrations are stable.

## 4. Core design decisions

### 4.1 Use small capability interfaces

Providers do not all support the same channels or message types. The SDK must therefore use focused interfaces instead of requiring every provider to implement one large interface.

```go
type SMSSender interface {
    SendSMS(ctx context.Context, message SMSMessage) (*SendResult, error)
}

type WhatsAppTextSender interface {
    SendWhatsAppText(ctx context.Context, message WhatsAppText) (*SendResult, error)
}

type WhatsAppTemplateSender interface {
    SendWhatsAppTemplate(ctx context.Context, message WhatsAppTemplate) (*SendResult, error)
}

type WhatsAppMediaSender interface {
    SendWhatsAppMedia(ctx context.Context, message WhatsAppMedia) (*SendResult, error)
}
```

Application code can accept only the capability it needs:

```go
func SendLoginCode(ctx context.Context, sender messaging.SMSSender, phone, code string) error {
    _, err := sender.SendSMS(ctx, messaging.SMSMessage{
        To:   phone,
        Body: "Your login code is " + code,
    })
    return err
}
```

Africa's Talking, SMSLeopard, and Wasiliana will initially implement `SMSSender`. Safravo will implement `SMSSender` and the applicable WhatsApp interfaces confirmed by the discovery gate.

### 4.2 Keep shared types intentionally small

Shared message types should contain only fields with consistent meaning across providers. Provider-only options belong in the provider package.

For example, the common SMS type may contain:

```go
type SMSMessage struct {
    To          string
    From        string
    Body        string
    CallbackURL string
    Metadata    map[string]string
}
```

If Africa's Talking supports a setting that has no equivalent elsewhere, its client should expose a provider-specific method or option rather than adding a misleading field to every message.

The initial common send method handles one recipient. Providers may expose batch-send methods from their concrete clients, but batch behavior will not be forced into the first shared interface.

### 4.3 Normalize conservatively

`SendResult` will expose fields applications commonly need while retaining the provider's original data.

```go
type SendResult struct {
    MessageID      string
    Status         MessageStatus
    Provider       string
    ProviderStatus string
    Raw            json.RawMessage
}
```

The normalized status set should remain small:

```go
const (
    StatusUnknown   MessageStatus = "unknown"
    StatusAccepted  MessageStatus = "accepted"
    StatusQueued    MessageStatus = "queued"
    StatusSent      MessageStatus = "sent"
    StatusDelivered MessageStatus = "delivered"
    StatusRead      MessageStatus = "read"
    StatusFailed    MessageStatus = "failed"
)
```

The exact status returned by the provider remains available in `ProviderStatus`. Unknown provider states must map to `StatusUnknown`; the SDK must not guess.

### 4.4 Use explicit provider construction

Each provider will have its own configuration and constructor:

```go
at, err := africastalking.New(africastalking.Config{
    Username: os.Getenv("AT_USERNAME"),
    APIKey:   os.Getenv("AT_API_KEY"),
})
```

The initial SDK will not use global provider registration through `init()`. Explicit constructors make dependencies, configuration errors, and tests easier to understand. A high-level factory can be considered later only if real usage shows a need for runtime provider selection.

### 4.5 Share infrastructure, not provider payloads

The core library may share HTTP request execution, safe logging, errors, and status types. Authentication headers, request bodies, response structures, and webhook payloads remain inside each provider package.

Shared webhook code should be extracted only after at least two provider implementations demonstrate the same behavior.

### 4.6 Use direct HTTP integrations by default

Africa's Talking, SMSLeopard, and Safravo will use their documented REST APIs through Go's `net/http` package. Safravo does not provide an official Go SDK, so its adapter will keep authentication, payload mapping, webhook handling, and raw-response access within the provider package.

### 4.7 Treat provider documentation review as an implementation gate

Provider APIs change. Before implementing any provider, confirm its current official documentation, authentication scheme, sandbox behavior, endpoints, limits, webhook security, and terms. Record the findings in that provider's package documentation or a short design note. Do not implement from the research summary alone.

### 4.8 Communication components only

This SDK is a messaging SDK. Some candidate providers (for example Wasiliana) also offer non-communication products such as airtime, payments, or billing APIs. Adapters must ignore those products entirely and implement only the communication components the SDK covers: SMS, USSD-related messaging where applicable, and WhatsApp. Payment, airtime, and similar components are out of scope and must not appear in the public API, capability matrix, or provider package.

## 5. Proposed repository layout

The final module path and public package name should be chosen before Phase 0 is merged. The layout below uses `messaging` as the conceptual root package.

```
.
├── go.mod
├── go.sum
├── message.go                    # Shared SMS and WhatsApp message types
├── result.go                     # SendResult and normalized statuses
├── errors.go                     # Shared typed errors and error categories
├── interfaces.go                 # Capability interfaces
├── options.go                    # Shared client/transport options
│
├── internal/
│   ├── transport/
│   │   ├── client.go             # HTTP execution, timeouts, safe retry hooks
│   │   ├── logging.go            # Redacted request diagnostics
│   │   └── client_test.go
│   └── testutil/
│       ├── server.go             # Helpers for local HTTP test servers
│       └── fixtures.go           # Fixture-loading helpers
│
├── fake/
│   ├── sender.go                 # Configurable test implementation
│   └── sender_test.go
│
├── africastalking/
│   ├── client.go
│   ├── sms.go
│   ├── types.go
│   ├── errors.go
│   ├── webhook.go
│   ├── client_test.go
│   └── testdata/
│
├── smsleopard/
│   ├── client.go
│   ├── sms.go
│   ├── types.go
│   ├── errors.go
│   ├── webhook.go
│   ├── client_test.go
│   └── testdata/
│
├── wasiliana/
│   ├── client.go
│   ├── sms.go
│   ├── types.go
│   ├── errors.go
│   ├── webhook.go
│   ├── client_test.go
│   └── testdata/
│
├── safravo/
│   ├── client.go
│   ├── sms.go
│   ├── whatsapp.go
│   ├── types.go
│   ├── errors.go
│   ├── webhook.go
│   ├── client_test.go
│   └── testdata/
│
├── examples/
│   ├── sms/
│   ├── whatsapp_text/
│   └── webhook_server/
│
├── docs/
│   ├── adding-a-provider.md
│   ├── provider-capabilities.md
│   └── design-decisions.md
│
├── .github/
│   ├── workflows/ci.yml
│   ├── ISSUE_TEMPLATE/
│   └── pull_request_template.md
│
├── README.md
├── CONTRIBUTING.md
├── SECURITY.md
├── CODE_OF_CONDUCT.md
├── CHANGELOG.md
├── LICENSE
└── implementation_plan.md
```

Files should be split when doing so improves clarity. The tree is a guide, not a requirement to create empty files before their corresponding phase.

## 6. Naming conventions

| Item | Convention | Example |
| --- | --- | --- |
| Package | Lowercase, no underscores | `smsleopard` |
| File | Lowercase snake case | `message_status.go` |
| Provider client | `Client` inside provider package | `safravo.Client` |
| Constructor | `New` | `safravo.New(config)` |
| Configuration | `Config` inside provider package | `africastalking.Config` |
| Request/message type | Describes the channel and operation | `WhatsAppTemplate` |
| Provider wire type | Unexported unless users need it | `sendSMSResponse` |
| Sentinel error | `Err` plus description | `ErrAuthentication` |
| Tests | Adjacent `_test.go` file | `sms_test.go` |
| Fixtures | Stored under `testdata/` | `send_success.json` |

Exported identifiers must have Go documentation comments. Abbreviations such as SMS and URL should use conventional Go capitalization.

## 7. Error model

Callers need to distinguish invalid input, authentication failures, rate limits, temporary provider failures, and permanent message rejection.

The core package should define a typed error with a small category set:

```go
type ErrorKind string

const (
    ErrorInvalidRequest  ErrorKind = "invalid_request"
    ErrorAuthentication  ErrorKind = "authentication"
    ErrorAuthorization   ErrorKind = "authorization"
    ErrorRateLimited     ErrorKind = "rate_limited"
    ErrorRejected        ErrorKind = "rejected"
    ErrorTemporary       ErrorKind = "temporary"
    ErrorTransport       ErrorKind = "transport"
    ErrorProvider        ErrorKind = "provider"
)

type ProviderError struct {
    Provider   string
    Kind       ErrorKind
    Code       string
    Message    string
    Retryable  bool
    HTTPStatus int
    Raw        json.RawMessage
}
```

Provider packages translate their error responses into this type. The original provider code and safe raw response are preserved for diagnosis. Error strings must never contain API keys, authorization headers, or message bodies.

Input validation should catch clear local errors, such as an empty recipient or body, before a network request. The SDK should avoid overly strict phone-number validation because acceptable formats may vary by provider and account.

## 8. HTTP, retry, and logging rules

### HTTP clients

- Accept an injected `http.Client` for testing and custom transports.
- Provide a sensible default timeout when the caller does not supply a client.
- Use `http.NewRequestWithContext` so cancellation and deadlines work.
- Close response bodies and limit the amount read into memory.
- Set an identifiable user agent containing the SDK version when practical.
- Do not mutate a caller-provided HTTP client.

### Retries

Sending a message can create a duplicate. Automatic retries must therefore be conservative:

- Do not automatically retry a send merely because the connection failed after the request may have reached the provider.
- Retry rate limits or temporary failures only when the operation is known to be safe, or when the provider supports an idempotency mechanism that is in use.
- Respect `Retry-After` when provided.
- Bound retry count and elapsed time.
- Document the final retry policy for each provider.

### Logging and privacy

- Logging is opt-in and uses `log/slog` rather than a custom logging framework.
- Never log credentials, authorization headers, full phone numbers, message bodies, template parameters, or complete webhook payloads by default.
- Diagnostic logs may include provider name, operation, HTTP status, duration, normalized error kind, and a redacted request ID.
- Tests must verify that representative secrets are redacted.

## 9. Webhook model

The SDK will help applications validate and parse provider callbacks, but it will not run an HTTP server for them.

Common event types should cover:

- Delivery status updates.
- Inbound SMS messages.
- Inbound WhatsApp messages.
- WhatsApp delivered and read receipts.

An event should retain both normalized and provider-specific information:

```go
type DeliveryEvent struct {
    Provider       string
    MessageID      string
    Status         MessageStatus
    ProviderStatus string
    OccurredAt     time.Time
    Recipient      string
    FailureCode    string
    FailureMessage string
    Raw            json.RawMessage
}
```

Each provider package owns its parsing and verification functions. Their APIs should accept the information required to verify the exact request, which may include the request URL, headers, and raw body. Signature verification must use constant-time comparison where appropriate.

The example webhook server must demonstrate:

- Limiting the request-body size.
- Reading and preserving the raw body.
- Verifying the provider signature before trusting the payload.
- Parsing the event.
- Returning the response expected by that provider.
- Handling duplicate callbacks idempotently at the application layer.

## 10. Implementation phases

### Phase 0 — Project foundation

**Goal:** Establish the public core, test tooling, CI, and open-source baseline.

**Work**

- Choose the module path, package name, supported Go versions, and license.
- Initialize `go.mod`.
- Add the shared interfaces and initial message/result/error types.
- Implement the internal HTTP transport rules required by the first provider.
- Implement a configurable fake SMS sender.
- Add test helpers based on `httptest.Server`.
- Configure CI to run:
  - `gofmt` verification
  - `go test ./...`
  - `go test -race ./...`
  - `go vet ./...`
  - the selected linter
- Add the initial README, contribution guide, security policy, code of conduct, and license.
- Record design decisions that affect the public API.

The fake sender should allow a test to preconfigure a result or error and should record received messages. It must not use magic phone-number or message-body values to decide whether a call succeeds.

**Tests**

- Compile-time assertions for fake capability implementations.
- Input validation tests.
- Error wrapping and `errors.As` tests.
- HTTP timeout and cancellation tests.
- Secret-redaction tests.
- Race tests for shared client state.

**Deliverable**

An application can depend on `SMSSender`, use the fake in its own tests, and run all repository checks without provider credentials.

**Milestone**

All CI checks pass and the public foundation is documented.

### Phase 1 — Africa's Talking SMS

**Goal:** Deliver the first real provider and establish the reference adapter structure.

**Discovery gate**

Before coding, verify and document:

- Current official SMS endpoint and environments.
- Authentication headers and required credentials.
- Sandbox account and simulator behavior.
- Sender-ID and recipient formatting rules.
- Single-recipient and bulk-send behavior.
- Response format, error format, and provider status values.
- Delivery-report and inbound-SMS callback formats.
- Whether callbacks can be authenticated or verified.
- Rate limits and provider retry guidance.

**Work**

- Add `africastalking.Config` with validation.
- Add `africastalking.New` with optional HTTP-client injection.
- Implement `SMSSender` for one-recipient SMS.
- Translate the common `SMSMessage` to the provider payload.
- Map provider results to `SendResult`.
- Translate documented provider failures to `ProviderError`.
- Add delivery-report parsing.
- Add inbound-SMS parsing if supported by the documented API.
- Preserve safe raw provider responses.
- Add a provider README or package example.
- Expose bulk sending only as an Africa's Talking-specific method if included.

**Tests**

- Constructor and configuration validation.
- Correct endpoint, headers, encoding, and request body.
- Successful sandbox-style response.
- Partial or per-recipient failure when applicable.
- Authentication, validation, rate-limit, server, malformed-response, timeout, and cancellation cases.
- Delivery and inbound webhook fixtures.
- Redaction of the API key, recipient, and message body.
- Compile-time assertion that the client implements `SMSSender`.

Tests use a local HTTP server and sanitized fixtures. Live sandbox tests must be optional, excluded from normal CI, and enabled only with explicit environment variables.

**Deliverable**

Users can send an SMS through Africa's Talking with the common interface and parse its supported callbacks.

**Milestone**

The first end-to-end provider example works in the official sandbox, and all offline tests pass without credentials.

### Phase 2 — SMSLeopard SMS

**Goal:** Add a second local provider and use real differences to validate the core design.

**Discovery gate**

Confirm the current official documentation for:

- Base URL and send endpoint.
- Basic Authentication construction and credential handling.
- JSON or form request encoding.
- Transactional versus bulk-message routing.
- Sender-ID requirements.
- Response, error, balance, and delivery-status formats.
- Callback configuration and verification support.
- Rate limits, retry guidance, and sandbox availability.

**Work**

- Add `smsleopard.Config`, constructor, and validation.
- Implement `SMSSender`.
- Implement provider request and response translation.
- Map SMSLeopard statuses and errors conservatively.
- Add delivery-report and inbound-message parsers where documented.
- Add examples and sanitized fixtures.
- Compare the implementation with Africa's Talking.
- Change the shared API only when both providers demonstrate a common need.
- Record important differences in the provider capability matrix.

**Tests**

- Cover the same categories as the Africa's Talking adapter, with additional tests for SMSLeopard's authentication, encoding, and provider-specific response shape.

**Design review**

At the end of this phase, review:

- Whether `SMSMessage` contains only genuinely shared fields.
- Whether `SendResult` represents both providers without losing important data.
- Whether normalized statuses and errors remain accurate.
- Whether the internal transport is reusable without provider conditionals.
- Whether any public API should change before WhatsApp support begins.

This is the preferred point for necessary breaking changes because the project has not reached v1.0.

**Deliverable**

The same application function can send through Africa's Talking or SMSLeopard by receiving an `SMSSender`.

**Milestone**

Two independently tested SMS providers pass CI and share no provider-specific wire structures.

### Phase 3 — Wasiliana SMS

**Goal:** Add a third local SMS provider and validate the adapter pattern against an `apiKey`-header based API ecosystem before the WhatsApp work begins.

**Discovery gate**

Confirm the current official documentation (https://docs.wasiliana.com/) for:

- Base URL and send endpoint (`POST https://api.wasiliana.com/api/v1/send/sms`).
- `apiKey` header authentication and credential handling.
- JSON request encoding and the documented request fields: `recipients`, `from`, `message`, and the optional `linkid`, `message_uid`, and `is_otp` fields.
- Sender-ID and recipient formatting rules, including Kenyan number formats.
- Response format, error format, and delivery-status values.
- Delivery-report callback configuration and verification support.
- Rate limits, retry guidance, and sandbox or test-account availability.

**Work**

- Add `wasiliana.Config`, constructor, and validation.
- Implement `SMSSender`.
- Implement `apiKey`-header authentication.
- Translate the common `SMSMessage` to the Wasiliana JSON payload, mapping `To`, `From`, and `Body` to `recipients`, `from`, and `message`.
- Map Wasiliana responses, statuses, and errors conservatively to `SendResult` and `ProviderError`.
- Add delivery-report parsing where documented.
- Add examples and sanitized fixtures.
- Compare the implementation with Africa's Talking and SMSLeopard.
- Change the shared API only when multiple providers demonstrate a common need.
- Record important differences in the provider capability matrix.

Wasiliana's non-communication products (airtime, payments) are out of scope per section 4.8 and must be ignored in the adapter, documentation, and capability matrix.

**Tests**

- Cover the same categories as the Africa's Talking adapter, with additional tests for Wasiliana's `apiKey` header, recipients-array encoding, optional `message_uid`/`is_otp` fields, and provider-specific response shape.

**Deliverable**

Users can send an SMS through Wasiliana with the common `SMSSender` interface.

**Milestone**

Three independently tested SMS providers pass CI and share no provider-specific wire structures.

### Phase 4 — Safravo SMS

**Goal:** Add a Kenyan omnichannel provider and validate the SDK against a different API ecosystem before extending Safravo to WhatsApp.

**Discovery and design gate**

- Confirm current official SMS API and webhook documentation.
- Confirm that SMS sending is exposed to third-party API clients and identify the required account or product configuration.
- Confirm authentication, Kenyan sender-ID and recipient requirements, sandbox or test support, status callbacks, webhook verification, rate limits, and idempotency behavior.

**Work**

- Add `safravo.Config`, constructor, and validation.
- Implement `SMSSender`.
- Support the documented sender identity options needed for SMS.
- Map Safravo message identifiers, statuses, and errors.
- Implement and test Safravo webhook verification when officially documented.
- Parse SMS delivery callbacks and inbound messages.
- Document Safravo sandbox or test-account limitations.

**Tests**

- Authentication and request mapping.
- Success, rejection, rate-limit, and provider-error responses.
- All supported status mappings, including unknown future statuses.
- Valid and invalid webhook authentication or signatures, when supported.
- Any verification cases where the original URL, headers, or payload affect verification.
- Local unit tests plus optional credential-gated integration tests.

**Deliverable**

Safravo can replace either local provider in application code that depends only on `SMSSender`.

**Milestone**

Four SMS implementations satisfy the same interface and the capability matrix accurately documents their differences.

### Phase 5 — Safravo WhatsApp and shared webhook events

**Goal:** Introduce WhatsApp without weakening or overloading the SMS API.

**Discovery gate**

Confirm current official rules for:

- WhatsApp sender and recipient identifiers.
- Customer-service windows and business-initiated messages.
- Approved content templates and parameter formats.
- Supported media types, size limits, and URLs.
- Embedded signup, account onboarding, and sandbox or test behavior.
- Delivery, read, failure, and inbound-message callbacks.
- Webhook authentication or signature verification.

**Work**

- Finalize the common WhatsApp text, template, and media types.
- Implement `WhatsAppTextSender` in the Safravo client.
- Implement `WhatsAppTemplateSender` using current Safravo and Meta template concepts.
- Implement `WhatsAppMediaSender` for supported media.
- Parse inbound WhatsApp messages.
- Parse sent, delivered, read, and failed status callbacks.
- Introduce common delivery/inbound event types only where the SMS and WhatsApp implementations demonstrate stable common fields.
- Keep Safravo-only fields in Safravo event structures or raw payloads.
- Add runnable examples for text, template, media, and webhook processing.

Interactive WhatsApp messages should be added in this phase only if the current Safravo API supports them cleanly and their data model is understood. Otherwise, they should be documented as a later enhancement.

**Tests**

- Text, template, and media request mapping.
- Required-field validation for each message kind.
- Template parameters and rejected-template errors.
- Status normalization, including read and unknown statuses.
- Inbound text and media fixtures.
- Signature verification and tampered-request cases.
- Redaction of phone numbers, message content, media URLs, and parameters.

**Deliverable**

Users can send the supported WhatsApp message types and safely process Safravo callbacks without treating WhatsApp as a special form of SMS.

**Milestone**

The SDK supports SMS through four providers and WhatsApp through its first provider, with documented capability interfaces for both channels.

### Phase 6 — Developer experience and contributor readiness

**Goal:** Make the project easy to adopt, review, and extend.

**Work**

- Complete the root README with installation and minimal examples.
- Publish a provider capability matrix showing:
  - SMS sending
  - WhatsApp text, template, and media support
  - Inbound-message parsing
  - Delivery callbacks
  - Signature verification
  - Sandbox or test support
- Write `docs/adding-a-provider.md` as a step-by-step adapter checklist.
- Add Go examples that are compiled by `go test` where practical.
- Document configuration, context cancellation, errors, retries, logging, privacy, and webhook security.
- Add issue templates for bugs, provider requests, and provider proposals.
- Add a pull-request checklist.
- Define API deprecation and provider maintenance policies.
- Add changelog automation or a clearly documented manual process.
- Audit exported names and documentation with `go doc`.

**Contributor provider checklist**

A provider contribution must:

- Link to current official API documentation.
- State whether the integration is official, community-supported, or unofficial.
- Implement only the capability interfaces it truly supports.
- Keep wire payloads in its provider package.
- Validate configuration without making a network request.
- Accept contexts and an injectable HTTP client.
- Translate documented errors and preserve safe provider details.
- Include sanitized success and failure fixtures.
- Test authentication, request mapping, errors, malformed responses, cancellation, and secret redaction.
- Add documentation and update the capability matrix.
- Pass formatting, tests, race detection, vetting, and linting.
- Avoid live credentials and network access in normal CI.

**Deliverable**

A contributor unfamiliar with the original implementation can add or maintain a provider by following repository documentation and existing adapter examples.

**Milestone**

All public packages are documented, examples compile, and the provider addition process has been tested through review of the four starter integrations.

### Phase 7 — Production hardening and v1.0

**Goal:** Stabilize the public API and prepare a trustworthy first major release.

**Work**

- Review public interfaces for the minimum stable v1 surface.
- Run the full test suite with the race detector.
- Add fuzz tests for webhook and provider-response parsers.
- Review retry behavior for duplicate-message risk.
- Review response-body limits and malformed payload handling.
- Review concurrent client use and any token caching.
- Review logs and errors for credentials and personal information.
- Add OpenTelemetry hooks only if they can remain optional and lightweight.
- Benchmark hot paths where performance is relevant.
- Verify examples against current provider sandboxes or test credentials.
- Document known provider limitations and operational requirements.
- Complete CHANGELOG.md and release notes.
- Tag v1.0.0 only after the compatibility and security review is complete.

**Milestone**

The public API is stable, all automated checks pass, the supported provider matrix is accurate, and v1.0.0 is ready for release.

## 11. Testing strategy

The test suite should use several layers:

- **Unit tests** — Test validation, mappings, status conversion, error conversion, redaction, and webhook parsing with table-driven tests.
- **HTTP contract tests** — Use `httptest.Server` to assert the exact method, path, headers, and encoded payload sent by each client. Return sanitized provider fixtures to test parsing. These tests run offline and form the main CI suite.
- **Compile-time interface checks** — Each provider should explicitly assert supported capabilities:

  ```go
  var _ messaging.SMSSender = (*Client)(nil)
  ```

  No assertion should exist for unsupported capabilities.
- **Fuzz tests** — Fuzz externally supplied data, especially webhook bodies, provider responses, timestamps, and status values. Parsers must return controlled errors rather than panic.
- **Integration tests** — Provider sandbox or test-credential tests must be opt-in. They should skip when required environment variables are absent and should never send to real users. Normal pull-request CI must not depend on provider uptime.
- **Security tests** — Include invalid signatures, oversized webhook bodies in examples, malformed JSON, unexpected content types, secret-redaction cases, and error responses that echo request data.

## 12. Definition of done

A phase or provider is complete only when:

- The agreed functionality is implemented.
- Public behavior is documented.
- Offline success and failure tests pass.
- Credentials are not required for normal tests.
- Context cancellation and HTTP timeouts work.
- Errors are useful and contain no secrets.
- Logs redact private message data.
- Fixtures are sanitized and stored under `testdata/`.
- The capability matrix is updated.
- Examples compile and match the current API.
- Formatting, tests, race detection, vetting, and linting pass.
- At least one maintainer other than the author has reviewed the change.

## 13. Future and community providers

After the four starter providers are stable, additional adapters can be added by maintainers or community contributors. The following list is an initial roadmap, not a promise that every provider will be implemented or permanently maintained by the initial project team.

### Candidate providers

- **Celcom Africa** — enterprise SMS and possible USSD/WhatsApp capabilities.
- **MoveSMS** — local bulk and transactional SMS.
- **WASMS** — SMS and WhatsApp services.
- **Infobip** — global SMS, WhatsApp, and other messaging channels.
- **Wati** — WhatsApp Business messaging.
- **Meta WhatsApp Cloud API** — direct official WhatsApp integration without an intermediary provider.

Provider priority should be based on user demand, accessible official documentation, sandbox availability, maintainer interest, and the ability to test the integration reliably.

Each added provider must follow the same capability-based design and contributor checklist. A provider must not force unrelated features into the core API.

### Experimental unofficial integrations

An integration based on whatsmeow or another WhatsApp Web protocol library must not be presented as equivalent to an official WhatsApp Business API. If accepted in the future, it should live in a clearly marked experimental module or package and document:

- Its unofficial status.
- Possible conflict with WhatsApp terms.
- Account suspension or banning risk.
- Protocol breakage and maintenance risk.
- Session storage and connection-lifecycle responsibilities.

It is deliberately outside the initial implementation phases.

## 14. Recommended implementation order

Work should be merged in small, reviewable changes. A practical sequence is:

1. Project metadata, CI, and core interfaces.
2. Shared results and typed errors.
3. HTTP transport and test utilities.
4. Fake sender.
5. Africa's Talking discovery note and SMS adapter.
6. Africa's Talking callbacks and examples.
7. SMSLeopard discovery note and SMS adapter.
8. Two-provider API review and any pre-v1 corrections.
9. Safravo discovery note and SMS adapter.
10. Safravo webhook verification and SMS callbacks.
11. WhatsApp common types and Safravo text support.
12. Safravo templates, media, and WhatsApp callbacks.
13. Contributor documentation and capability matrix.
14. Hardening, compatibility review, and v1.0 release.

Each change should leave the repository buildable and keep all existing tests passing. Large phases may be split among contributors, but shared public types should be agreed and merged before provider work that depends on them.