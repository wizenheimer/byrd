package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/internal/domain/models"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type NotificationHandler struct {
	notificationService interfaces.NotificationService
	logger              *logger.Logger
}

func NewNotificationHandler(notificationService interfaces.NotificationService, logger *logger.Logger) *NotificationHandler {
	logger.Debug("creating new notification handler")

	return &NotificationHandler{
		notificationService: notificationService,
		logger:              logger.WithFields(map[string]interface{}{"module": "notification_handler"}),
	}
}

func (h *NotificationHandler) SendNotification(c *fiber.Ctx) error {
	h.logger.Debug("sending notification")

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
