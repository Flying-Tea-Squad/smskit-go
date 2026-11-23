// Package wasiliana implements smskit capabilities for the Wasiliana SMS API
// (https://docs.wasiliana.com). Airtime, USSD and payments are out of scope.
//
// Wasiliana's send response carries no message ID, so SendResult.MessageID is
// always empty. SendSMS never retries: a failed connection may still have
// reached the provider, and a retry could deliver the message twice.
package wasiliana
