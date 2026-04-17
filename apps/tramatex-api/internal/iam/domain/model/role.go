package model

import "fmt"

// Role represents a user role in the system.
type Role string

const (
	RoleAdmin      Role = "admin"
	RoleCommercial Role = "commercial"
	RoleCashier    Role = "cashier"
	RoleDesigner   Role = "designer"
	RoleWorkshop   Role = "workshop"
)

// IsValid validates if a role is one of the allowed roles.
func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleCommercial, RoleCashier, RoleDesigner, RoleWorkshop:
		return true
	default:
		return false
	}
}

// NewRole creates a Role from string and validates it.
func NewRole(value string) (Role, error) {
	role := Role(value)
	if !role.IsValid() {
		return "", fmt.Errorf("invalid role: %s", value)
	}
	return role, nil
}
