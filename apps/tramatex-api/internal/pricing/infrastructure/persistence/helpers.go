package persistence

import "github.com/google/uuid"

func uuidArrayToStringArray(values []uuid.UUID) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, value.String())
	}
	return result
}

func stringArrayToUUIDs(values []string) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		parsed, err := uuid.Parse(value)
		if err != nil {
			continue
		}
		result = append(result, parsed)
	}
	return result
}
