package services

import "errors"

// Sentinel errors shared across services
var ErrCannotCancel = errors.New("cannot cancel order that has been shipped or delivered")
var ErrNotFound = errors.New("document not found")
