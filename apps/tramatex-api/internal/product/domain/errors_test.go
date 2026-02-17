package domain

import (
	"errors"
	"testing"
)

func TestProductErrorCodes(t *testing.T) {
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
		var productErr ProductError
		if !errors.As(test.err, &productErr) {
			t.Fatalf("expected ProductError, got %T", test.err)
		}
		if productErr.Code != test.code {
			t.Fatalf("expected code %s, got %s", test.code, productErr.Code)
		}
	}
}

func TestProductErrorWrapping(t *testing.T) {
	cause := errors.New("root")
	err := WrapValidation("invalid", cause)

	var productErr ProductError
	if !errors.As(err, &productErr) {
		t.Fatalf("expected ProductError, got %T", err)
	}
	if productErr.Cause == nil {
		t.Fatalf("expected cause to be set")
	}
	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped cause")
	}
}

func TestProductErrorMessagesAndFormats(t *testing.T) {
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
			err:     NewNotFoundErrorf("missing %s", "product"),
			code:    ErrCodeNotFound,
			message: "missing product",
		},
		{
			name:    "not found wrap",
			err:     WrapNotFound("missing product", cause),
			code:    ErrCodeNotFound,
			message: "missing product",
			cause:   cause,
		},
		{
			name:    "not found wrap format",
			err:     WrapNotFoundf(cause, "missing %s", "product"),
			code:    ErrCodeNotFound,
			message: "missing product",
			cause:   cause,
		},
		{
			name:    "conflict format",
			err:     NewConflictErrorf("conflict %s", "sku"),
			code:    ErrCodeConflict,
			message: "conflict sku",
		},
		{
			name:    "conflict wrap",
			err:     WrapConflict("conflict sku", cause),
			code:    ErrCodeConflict,
			message: "conflict sku",
			cause:   cause,
		},
		{
			name:    "conflict wrap format",
			err:     WrapConflictf(cause, "conflict %s", "sku"),
			code:    ErrCodeConflict,
			message: "conflict sku",
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
			name:    "persistence wrap format",
			err:     WrapPersistencef(cause, "db %s", "down"),
			code:    ErrCodePersistence,
			message: "db down",
			cause:   cause,
		},
	}

	for _, testCase := range testCases {
		var productErr ProductError
		if !errors.As(testCase.err, &productErr) {
			t.Fatalf("%s: expected ProductError, got %T", testCase.name, testCase.err)
		}
		if productErr.Code != testCase.code {
			t.Fatalf("%s: expected code %s, got %s", testCase.name, testCase.code, productErr.Code)
		}
		if productErr.Message != testCase.message {
			t.Fatalf("%s: expected message %q, got %q", testCase.name, testCase.message, productErr.Message)
		}
		if testCase.cause == nil {
			if productErr.Error() != testCase.message {
				t.Fatalf("%s: expected error %q, got %q", testCase.name, testCase.message, productErr.Error())
			}
			continue
		}
		if !errors.Is(testCase.err, testCase.cause) {
			t.Fatalf("%s: expected wrapped cause", testCase.name)
		}
		expected := testCase.message + ": " + testCase.cause.Error()
		if productErr.Error() != expected {
			t.Fatalf("%s: expected error %q, got %q", testCase.name, expected, productErr.Error())
		}
	}
}
