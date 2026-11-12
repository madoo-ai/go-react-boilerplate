package entity

import "errors"

// ErrNotFound is returned when a record does not exist.
var ErrNotFound = errors.New("not found")

// CodedError carries a translation key that the HTTP layer maps to a status
// code and the React client maps to a localized message.
type CodedError struct {
	Code   string
	Fields map[string]string
}

func (e *CodedError) Error() string { return e.Code }

func NewError(code string) *CodedError { return &CodedError{Code: code} }

func ValidationError(fields map[string]string) *CodedError {
	return &CodedError{Code: "validation.failed", Fields: fields}
}

// Common application errors. The strings are translation keys, not English
// sentences, so every surface renders the message in the active locale.
var (
	ErrUnauthorized            = NewError("errors.unauthorized")
	ErrForbidden               = NewError("errors.forbidden")
	ErrTooManyRequests         = NewError("errors.tooManyRequests")
	ErrAccountPending          = NewError("auth.messages.pendingApproval")
	ErrAccountSuspended        = NewError("auth.messages.suspended")
	ErrInvalidCredentials      = NewError("errors.invalidCredentials")
	ErrIncorrectPassword       = NewError("errors.incorrectPassword")
	ErrInvalidToken            = NewError("errors.invalidToken")
	ErrInvalidCode             = NewError("errors.invalidCode")
	ErrUserNotFound            = NewError("errors.userNotFound")
	ErrRoleNotFound            = NewError("errors.roleNotFound")
	ErrEmailInUse              = NewError("errors.emailInUse")
	ErrRoleInUse               = NewError("errors.roleInUse")
	ErrCannotDeleteSelf        = NewError("errors.cannotDeleteSelf")
	ErrLastSuperadmin          = NewError("errors.lastSuperadmin")
	ErrProtectedUser           = NewError("errors.protectedUser")
	ErrImmutableRole           = NewError("errors.immutableRole")
	ErrDefaultRole             = NewError("errors.defaultRole")
	ErrConfirmPasswordRequired = NewError("errors.confirmPasswordRequired")
	ErrTwoFactorNotEnabled     = NewError("errors.twoFactorNotEnabled")
	ErrUnexpected              = NewError("errors.unexpected")
)
