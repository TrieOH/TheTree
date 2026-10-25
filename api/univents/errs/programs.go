package errs

import "errors"

// Program sign-ups, check-ins and attendance.
var (
	ErrTicketRequiredForActivity   = errors.New("a confirmed ticket is required to register for activities")
	ErrTicketRequiredForCheckpoint = errors.New("a confirmed ticket is required to pass this checkpoint")
	ErrAlreadyRegistered           = errors.New("you are already registered for this activity")
	ErrOccurrenceFull              = errors.New("this activity occurrence is full")
	ErrNotRegistered               = errors.New("you are not registered for this activity")
	ErrRegistrationLocked          = errors.New("registration is locked: attendance has been marked")
	ErrCheckInNotCheckpoint        = errors.New("check-in is only for checkpoint occurrences")
	ErrParticipationCancelled      = errors.New("cannot mark attendance: participation was cancelled")
)

// AccessLevelError is a ticket whose access level is below the program's.
type AccessLevelError struct{ Kind string }

func (e *AccessLevelError) Error() string {
	return "your ticket level does not grant access to this " + e.Kind
}
