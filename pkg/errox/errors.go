package errox

import (
	"fmt"
)

// Sentinel errors for generic error classes.
//
// Note that error messages are the only reliable indicator of the error type
// in some cases, e.g., when embedded into a GraphQL response, thus changing
// them might break error matching for some clients, e.g., the UI.
var (
	NotFound = makeSentinel("not found")

	// AlreadyExists indicates that a request is syntactically correct but can't
	// be fulfilled due to a conflict with existing data. Such request could
	// have succeeded given a different system state.
	AlreadyExists = makeSentinel("already exists")

	// InvalidArgs indicates bad user request. Such request will always fail
	// regardless of the system state.
	InvalidArgs = makeSentinel("invalid arguments")

	// ReferencedByAnotherObject indicates that the requested object cannot be
	// removed or updated because it is referred to / in use by another object.
	// Future system state might allow the request.
	ReferencedByAnotherObject = makeSentinel("referenced by another object")

	// InvariantViolation indicates that some internal invariant has been
	// violated and the underlying system component is in an inconsistent state.
	// *Must not* be used for invalid user input.
	InvariantViolation = makeSentinel("invariant violation")

	// NoCredentials indicates that valid credentials are required but have not
	// been provided. The request can not be processed.
	NoCredentials = makeSentinel("credentials not found")

	// NoValidRole indicates that though user credentials have been provided,
	// they do not specify a valid role. This usually happens because of
	// misconfigured access control. The effect is similar to NoCredentials.
	NoValidRole = makeSentinel("no valid role")

	// NotAuthorized indicates that valid credentials have been provided but are
	// insufficient for processing the request.
	NotAuthorized = makeSentinel("not authorized")

	// NoAuthzConfigured indicates that authorization is not implemented for a
	// service. This is a programming error.
	NoAuthzConfigured = makeSentinel("service authorization is misconfigured")

	// When adding a new error please update the translators in this package (gRPC, etc.).
)

// GenericNoValidRole wraps ErrNoValidRole with a generic error message.
func GenericNoValidRole() error {
	return fmt.Errorf("access for this user is not authorized: %w, please contact your system administrator",
		NoValidRole)
}

func explain(err error, explanation string) error {
	return fmt.Errorf("%w: %s", err, explanation)
}

// NewErrNotAuthorized wraps ErrNotAuthorized into an explanation.
func NewErrNotAuthorized(explanation string) error {
	return explain(NotAuthorized, explanation)
}

// NewErrNoCredentials wraps ErrNoCredentials into an explanation.
func NewErrNoCredentials(explanation string) error {
	return explain(NoCredentials, explanation)
}

// NewErrInvariantViolation wraps ErrInvariantViolation into an explanation.
func NewErrInvariantViolation(explanation string) error {
	return explain(InvariantViolation, explanation)
}

// NewErrInvalidArgs wraps ErrInvalidArgs into an explanation.
func NewErrInvalidArgs(explanation string) error {
	return explain(InvalidArgs, explanation)
}
