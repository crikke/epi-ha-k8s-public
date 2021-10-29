package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository struct{}

type Repository interface {
	GetUserById(ctx context.Context, id uuid.UUID) *User
	AddUser(ctx context.Context, user *User)
}

func (u *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) *User {
	return nil
}
