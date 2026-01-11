# Backend Tests Structure

## Overview

This directory contains all tests for the TramaTex backend, separated from the domain and application layers to maintain architectural purity and clean separation of concerns.

## Directory Structure

```
tests/
├── unit/
│   └── domain/
│       ├── user/
│       │   ├── email_test.go           (7 tests - Email Value Object)
│       │   ├── password_test.go        (8 tests - Password Value Object)
│       │   └── user_test.go            (9 tests - User Entity)
│       └── security/
│           └── jwt_test.go             (4 tests - TokenClaims Value Object)
└── integration/
    └── auth/
        └── login_test.go               (7 tests - LoginUseCase Integration)
```

## Running Tests

### All Tests
```bash
cd backend
go test ./tests/... -v
```

### Unit Tests Only
```bash
cd backend
go test ./tests/unit/... -v
```

### Integration Tests Only
```bash
cd backend
go test ./tests/integration/... -v
```

### Specific Package
```bash
cd backend
go test ./tests/unit/domain/user/... -v
```

### With Coverage
```bash
cd backend
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Organization

### Unit Tests (`tests/unit/domain/`)

Tests for **Domain Layer** Value Objects and Entities:

#### `user/email_test.go` (package `user_test`)
- Email validation (RFC format)
- Empty email rejection
- Very long email rejection
- Email immutability

#### `user/password_test.go` (package `user_test`)
- Password hashing with bcrypt (cost ≥ 10)
- Password verification (Matches)
- Password too short rejection (< 8 chars)
- Password too long rejection (> 72 chars)

#### `user/user_test.go` (package `user_test`)
- User creation with valid data
- User without email rejection
- User without password rejection
- User with role (User, Manager, Admin)
- User immutability

#### `security/jwt_test.go` (package `security_test`)
- TokenClaims creation with valid data
- TokenClaims with expired timestamp rejection
- TokenClaims with empty subject rejection
- TokenClaims String representation

### Integration Tests (`tests/integration/auth/`)

Tests for **Application Layer** Use Cases with mocked dependencies:

#### `auth/login_test.go` (package `auth_test`)
- LoginUseCase with valid credentials
- LoginUseCase with invalid email format
- LoginUseCase with user not found
- LoginUseCase with wrong password
- LoginUseCase with repository error propagation
- LoginUseCase with JWT generation error
- LoginUseCase token generation response

## Test Pattern: External Package Testing

All tests use Go's external package testing convention:

```go
// In tests/unit/domain/user/email_test.go
package user_test  // External package with _test suffix

import (
    "testing"
    "github.com/joran-cortez/tramatex/internal/domain/user"
)

func TestEmailNewWithValidFormat(t *testing.T) {
    e, err := user.NewEmail("user@example.com")  // Import tested package
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}
```

**Benefits:**
- ✅ Tests are external to domain layer
- ✅ Imported package is tested through its public API
- ✅ No circular imports
- ✅ Maintains clean architecture separation

## Mock Objects

Integration tests use mock implementations:

- `MockUserRepository` - Implements `user.UserRepository` interface
- `MockJWTService` - Implements `security.JWTService` interface

Mocks allow testing use cases without requiring actual database or JWT implementation.

## Code Coverage

Current coverage targets:

| Layer | Package | Target | Status |
|-------|---------|--------|--------|
| **Unit** | `domain/user` | ≥90% | ✅ Implemented |
| **Unit** | `domain/security` | ≥80% | ✅ Implemented |
| **Integration** | `application/auth` | ≥80% | ✅ Implemented |

## Test Count Summary

| Category | Count | Status |
|----------|-------|--------|
| Unit Tests (domain) | 28 | ✅ |
| Integration Tests | 7 | ✅ |
| **Total** | **35** | ✅ |

## Best Practices Applied

1. **TDD-First**: Tests designed before implementation
2. **Descriptive Names**: Test function names clearly describe what is tested
3. **AAA Pattern**: Arrange-Act-Assert structure in each test
4. **No Mock Unless Needed**: Mock only external dependencies (repository, JWT service)
5. **Table-Driven Tests**: Where applicable for multiple scenarios
6. **Error Messages**: Clear assertion messages for failures
7. **Independence**: Each test is independent and can run in any order

## Dependencies

Test files import from:
- Go standard library: `testing`, `time`, `errors`
- Internal packages: `github.com/joran-cortez/tramatex/internal/...`

No external testing frameworks required (uses Go's built-in `testing` package).

## Notes

- Domain layer (`backend/internal/domain/`) is **clean of test files**
- Application layer (`backend/internal/application/`) is **clean of test files**
- All tests reside in `tests/` directory for clean separation
- Tests verify behavior through public APIs only
- Mock objects simulate external dependencies (database, JWT)

---

Last updated: 2026-01-11
