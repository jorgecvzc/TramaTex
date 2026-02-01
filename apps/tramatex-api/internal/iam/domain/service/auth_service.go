package service

// AuthorizationService defines domain-level authorization checks.
type AuthorizationService interface {
	HasRole(userRole string, requiredRoles []string) bool
}
