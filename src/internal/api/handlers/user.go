package handlers

import "github.com/gofiber/fiber/v2"

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	return nil
}
