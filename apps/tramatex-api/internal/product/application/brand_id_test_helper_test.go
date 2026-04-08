package application_test

import "github.com/google/uuid"

func uuidPtr(id uuid.UUID) *uuid.UUID { return &id }
