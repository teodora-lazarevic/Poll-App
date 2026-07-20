package services

import "errors"

// Domain Errors
var (
	ErrUnauthorized = errors.New("Unauthorized action")
	ErrInvalidCreds = errors.New("Invalid credentials")
	ErrUserNotFound = errors.New("User not found")
	ErrInvalidInput = errors.New("Invalid username, email, or password")
	ErrUserExists   = errors.New("User already exists in the database")

	ErrPollNotFound       = errors.New("Poll not found")
	ErrPollOptionNotFound = errors.New("Poll option not found")
	ErrOptionNotInPoll    = errors.New("Option does not belong to this poll")
	ErrDuplicateOption    = errors.New("Duplicate option text")
	ErrMinOptionsRequired = errors.New("A poll must have at least two options")
	ErrAlreadyVoted       = errors.New("User has already voted in this poll")
)
