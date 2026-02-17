package persistence

import (
	"testing"

	"github.com/google/uuid"
)

func TestUUIDStringHelpers(t *testing.T) {
	uuid1 := uuid.New()
	uuid2 := uuid.New()

	strings := uuidArrayToStringArray([]uuid.UUID{uuid1, uuid2})
	if len(strings) != 2 {
		t.Fatalf("expected 2 strings")
	}

	parsed := stringArrayToUUIDs([]string{uuid1.String(), "invalid", uuid2.String()})
	if len(parsed) != 2 {
		t.Fatalf("expected 2 parsed UUIDs")
	}
	if parsed[0] != uuid1 || parsed[1] != uuid2 {
		t.Fatalf("unexpected parsed UUIDs")
	}
}
