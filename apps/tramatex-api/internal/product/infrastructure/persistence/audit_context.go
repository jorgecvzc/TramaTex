package persistence

import "context"

func actorIDFromContext(ctx context.Context) string {
	if actorID, ok := ctx.Value("userID").(string); ok && actorID != "" {
		return actorID
	}
	if actorID, ok := ctx.Value("actorID").(string); ok && actorID != "" {
		return actorID
	}
	return "system"
}
