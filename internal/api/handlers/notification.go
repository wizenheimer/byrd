package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
)

type NotificationHandler struct {
	notificationService interfaces.NotificationService
}

func NewNotificationHandler(notificationService interfaces.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) SendNotification(c *fiber.Ctx) error {
	var req models.NotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	result, err := h.notificationService.SendNotification(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}