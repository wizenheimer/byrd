package handlers

import (
	"image"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type AIHandler struct {
	aiService interfaces.AIService
	logger    *logger.Logger
}

func NewAIHandler(aiService interfaces.AIService, logger *logger.Logger) *AIHandler {
	logger.Debug("creating new AI handler")

	return &AIHandler{
		aiService: aiService,
		logger:    logger.WithFields(map[string]interface{}{"module": "ai_handler"}),
	}
}

func (h *AIHandler) AnalyzeContentDifferences(c *fiber.Ctx) error {
	h.logger.Debug("analyzing content differences")

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid multipart form")
	}

	// Get files from form
	version1Files := form.File["version1"]
	version2Files := form.File["version2"]

	// Validate files
	if len(version1Files) == 0 || len(version2Files) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Both version files are required")
	}

	// Get first file from each
	file1Header := version1Files[0]
	file2Header := version2Files[0]

	// Process Version 1
	file1, err := file1Header.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open version1")
	}
	defer file1.Close()

	content1, err := io.ReadAll(file1)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read version1")
	}

	// Process Version 2
	file2, err := file2Header.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open version2")
	}
	defer file2.Close()

	content2, err := io.ReadAll(file2)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to read version2")
	}

	// Get profile from form values
	profiles := form.Value["profile_fields"]
	if len(profiles) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "profile_fields is required")
	}
	profileFieldsString := profiles[0]

	// Convert profile fields to slice
	profileFields := strings.Split(profileFieldsString, ",")

	result, err := h.aiService.AnalyzeContentDifferences(
		c.Context(),
		string(content1),
		string(content2),
		profileFields,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func (h *AIHandler) AnalyzeVisualDifferences(c *fiber.Ctx) error {
	h.logger.Debug("analyzing visual differences")

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid multipart form")
	}

	// Get files from form
	screenshots1 := form.File["screenshot1"]
	screenshots2 := form.File["screenshot2"]

	// Validate files
	if len(screenshots1) == 0 || len(screenshots2) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Both screenshots are required")
	}

	// Get first file from each
	file1Header := screenshots1[0]
	file2Header := screenshots2[0]

	// Process Screenshot1
	file1, err := file1Header.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open screenshot1")
	}
	defer file1.Close()

	img1, _, err := image.Decode(file1)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid image format for screenshot1")
	}

	// Process Screenshot2
	file2, err := file2Header.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open screenshot2")
	}
	defer file2.Close()

	img2, _, err := image.Decode(file2)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid image format for screenshot2")
	}

	// Get profile from form values
	profiles := form.Value["profile_fields"]
	if len(profiles) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "profile_fields is required")
	}
	profileString := profiles[0]

	// Convert profile fields to slice
	profileFields := strings.Split(profileString, ",")

	result, err := h.aiService.AnalyzeVisualDifferences(c.Context(), img1, img2, profileFields)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}
