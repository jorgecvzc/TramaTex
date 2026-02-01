package model

import "fmt"

// Permission represents a specific action a user can perform.
// This is a domain entity used by IAM for authorization decisions.
type Permission struct {
	id   string
	name string
}

// NewPermission creates a new permission with validation.
func NewPermission(id, name string) (*Permission, error) {
	if id == "" {
		return nil, fmt.Errorf("permission id cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("permission name cannot be empty")
	}
	return &Permission{id: id, name: name}, nil
}

// ID returns the permission ID.
func (p *Permission) ID() string {
	return p.id
}

// Name returns the permission name.
func (p *Permission) Name() string {
	return p.name
}
