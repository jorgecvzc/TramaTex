package domain

import (
	"errors"
	"testing"
)

func TestPartyErrorCodes(t *testing.T) {
	cases := []struct {
		err  error
		code ErrorCode
	}{
		{NewValidationError("bad input"), ErrCodeValidation},
		{NewNotFoundError("missing"), ErrCodeNotFound},
		{NewConflictError("conflict"), ErrCodeConflict},
		{NewPersistenceError("db"), ErrCodePersistence},
	}

	for _, test := range cases {
		var partyErr PartyError
		if !errors.As(test.err, &partyErr) {
			t.Fatalf("expected PartyError, got %T", test.err)
		}
		if partyErr.Code != test.code {
			t.Fatalf("expected code %s, got %s", test.code, partyErr.Code)
		}
	}
}

func TestPartyErrorWrapping(t *testing.T) {
	cause := errors.New("root")
	err := WrapValidation("invalid", cause)

	var partyErr PartyError
	if !errors.As(err, &partyErr) {
		t.Fatalf("expected PartyError, got %T", err)
	}
	if partyErr.Cause == nil {
		t.Fatalf("expected cause to be set")
	}
	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped cause")
	}
}

func TestPartyErrorMessagesAndFormats(t *testing.T) {
	cause := errors.New("root")

	testCases := []struct {
		name    string
		err     error
		code    ErrorCode
		message string
		cause   error
	}{
		{
			name:    "validation message",
			err:     NewValidationError("bad input"),
			code:    ErrCodeValidation,
			message: "bad input",
		},
		{
			name:    "validation format",
			err:     NewValidationErrorf("bad %s", "input"),
			code:    ErrCodeValidation,
			message: "bad input",
		},
		{
			name:    "validation wrap format",
			err:     WrapValidationf(cause, "bad %s", "input"),
			code:    ErrCodeValidation,
			message: "bad input",
			cause:   cause,
		},
		{
			name:    "not found format",
			err:     NewNotFoundErrorf("missing %s", "party"),
			code:    ErrCodeNotFound,
			message: "missing party",
		},
		{
			name:    "not found wrap",
			err:     WrapNotFound("missing party", cause),
			code:    ErrCodeNotFound,
			message: "missing party",
			cause:   cause,
		},
		{
			name:    "not found wrap",
			err:     WrapNotFoundf(cause, "missing %s", "party"),
			code:    ErrCodeNotFound,
			message: "missing party",
			cause:   cause,
		},
		{
			name:    "conflict format",
			err:     NewConflictErrorf("conflict %s", "role"),
			code:    ErrCodeConflict,
			message: "conflict role",
		},
		{
			name:    "conflict wrap",
			err:     WrapConflict("conflict role", cause),
			code:    ErrCodeConflict,
			message: "conflict role",
			cause:   cause,
		},
		{
			name:    "conflict wrap",
			err:     WrapConflictf(cause, "conflict %s", "role"),
			code:    ErrCodeConflict,
			message: "conflict role",
			cause:   cause,
		},
		{
			name:    "persistence format",
			err:     NewPersistenceErrorf("db %s", "down"),
			code:    ErrCodePersistence,
			message: "db down",
		},
		{
			name:    "persistence wrap",
			err:     WrapPersistence("db down", cause),
			code:    ErrCodePersistence,
			message: "db down",
			cause:   cause,
		},
		{
			name:    "persistence wrap",
			err:     WrapPersistencef(cause, "db %s", "down"),
			code:    ErrCodePersistence,
			message: "db down",
			cause:   cause,
		},
	}

	for _, testCase := range testCases {
		var partyErr PartyError
		if !errors.As(testCase.err, &partyErr) {
			t.Fatalf("%s: expected PartyError, got %T", testCase.name, testCase.err)
		}
		if partyErr.Code != testCase.code {
			t.Fatalf("%s: expected code %s, got %s", testCase.name, testCase.code, partyErr.Code)
		}
		if partyErr.Message != testCase.message {
			t.Fatalf("%s: expected message %q, got %q", testCase.name, testCase.message, partyErr.Message)
		}
		if testCase.cause == nil {
			if partyErr.Error() != testCase.message {
				t.Fatalf("%s: expected error %q, got %q", testCase.name, testCase.message, partyErr.Error())
			}
			continue
		}
		if !errors.Is(testCase.err, testCase.cause) {
			t.Fatalf("%s: expected wrapped cause", testCase.name)
		}
		expected := testCase.message + ": " + testCase.cause.Error()
		if partyErr.Error() != expected {
			t.Fatalf("%s: expected error %q, got %q", testCase.name, expected, partyErr.Error())
		}
	}
}
