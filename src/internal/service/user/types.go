package user

import (
	repo "github.com/wizenheimer/iris/src/internal/interfaces/repository"
	svc "github.com/wizenheimer/iris/src/internal/interfaces/service"
)

// compile time check if the interface is implemented
// TODO: reduce overhead by passing stuff by reference
var _ svc.UserService = (*userService)(nil)

// TODO: rethink retrieval methods
type userService struct {
	userRepository repo.UserRepository
}
