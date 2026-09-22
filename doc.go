// Package smskit defines provider-independent messaging contracts.
//
// Provider adapters live in independent top-level packages. They may depend
// on this package and shared internal infrastructure, but provider-specific
// authentication and wire types do not belong in the root package.
package smskit
