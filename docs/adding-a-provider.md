# Adding a provider

Each adapter is a standalone top-level package that depends on the root
`smskit` package and may reuse code under `internal/`. Provider packages must
not import one another or move provider wire structures into the root package.

1. Complete a discovery note using current official documentation. Record authentication,
   endpoints, payloads, status/error formats, callbacks, limits, retry guidance, and test access.
2. Add a validated provider `Config` and `New` constructor. Configuration validation must not
   make a network request, and clients must accept an injected `*http.Client`.
3. Implement only the capability interfaces the provider genuinely supports.
4. Keep authentication, request/response structs, provider statuses, and webhook payloads in the
   provider package. Preserve safe raw provider data in shared results where useful.
5. Add compile-time assertions for each supported capability.
6. Add offline `httptest` contract tests for exact method, path, headers, encoding, success,
   authentication, validation, rate limits, malformed responses, cancellation, and redaction.
7. Store sanitized fixtures under the provider's `testdata/` directory. Live sandbox tests must
   be opt-in and skipped by normal CI.
8. Add provider documentation, examples, and an update to the capability matrix.

An adapter is ready for review only when formatting, tests, race detection, vetting, linting,
documentation, and the privacy checks pass.
