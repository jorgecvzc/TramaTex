package usecase

import (
	"context"
	"fmt"

	domain_repo "github.com/joran-cortez/tramatex/internal/iam/domain/repository"
)

// ListUsersUseCase lists all users.
type ListUsersUseCase struct {
	userRepo domain_repo.Repository
}

// NewListUsersUseCase creates a new list users use case.
func NewListUsersUseCase(userRepo domain_repo.Repository) *ListUsersUseCase {
	return &ListUsersUseCase{userRepo: userRepo}
}

// Execute returns a list of users.
func (uc *ListUsersUseCase) Execute(ctx context.Context) ([]UserDTO, error) {
	users, err := uc.userRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	result := make([]UserDTO, 0, len(users))
	for _, u := range users {
		result = append(result, UserDTO{
			ID:    u.ID(),
			Email: u.Email().Value(),
			Role:  string(u.Role()),
		})
	}

	return result, nil
}
