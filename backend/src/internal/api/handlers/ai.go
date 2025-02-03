// ./src/internal/api/handlers/ai.go
package handlers

import (
	"image"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type AIHandler struct {
	aiService ai.AIService
	processor *utils.MarkdownProcessor
	logger    *logger.Logger
}

func NewAIHandler(aiService ai.AIService, logger *logger.Logger) (*AIHandler, error) {
	processor, err := utils.NewMarkdownProcessor()
	if err != nil {
		return nil, err
	}

	ah := AIHandler{
		aiService: aiService,
		processor: processor,
		logger:    logger.WithFields(map[string]interface{}{"module": "ai_handler"}),
	}

	return &ah, nil
}

func (h *AIHandler) AnalyzeContentDifferences(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid multipart form", err.Error())
	}

	// Get files from form
	version1Files := form.File["version1"]
	version2Files := form.File["version2"]

	// Validate files
	if len(version1Files) == 0 || len(version2Files) == 0 {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Both versions are required", map[string]interface{}{
			"version_1": len(version1Files),
			"version_2": len(version2Files),
		})
	}

	// Get first file from each
	file1Header := version1Files[0]
	file2Header := version2Files[0]

	// Process Version 1
	file1, err := file1Header.Open()
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to open version1", err.Error())
	}
	defer file1.Close()

	content1, err := io.ReadAll(file1)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to read version1", err.Error())
	}

	// Process Version 2
	file2, err := file2Header.Open()
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to open version2", err.Error())
	}
	defer file2.Close()

	content2, err := io.ReadAll(file2)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to read version2", err.Error())
	}

	markdownContent1, err := h.processor.Process(string(content1))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to process markdown content 1", err.Error())
	}

	markdownContent2, err := h.processor.Process(string(content2))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to process markdown content 2", err.Error())
	}

	// Get profile from form values
	profiles := form.Value["profile_fields"]
	if len(profiles) == 0 {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "profile_fields is required", map[string]interface{}{
			"profiles": profiles,
		})
	}
	profileFieldsString := profiles[0]

	// Convert profile fields to slice
	profileFields := strings.Split(profileFieldsString, ",")

	result, err := h.aiService.AnalyzeContentDifferences(
		c.Context(),
		markdownContent1,
		markdownContent2,
		profileFields,
	)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Could not analyze content differences", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Content analysis successfully", result)
}

func (h *AIHandler) AnalyzeVisualDifferences(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid multipart form", err.Error())
	}

	// Get files from form
	screenshots1 := form.File["screenshot1"]
	screenshots2 := form.File["screenshot2"]

	// Validate files
	if len(screenshots1) == 0 || len(screenshots2) == 0 {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Both screenshots are required", map[string]interface{}{
			"screenshots_1_exists": len(screenshots1) == 0,
			"screenshots_2_exists": len(screenshots2) == 0,
		})
	}

	// Get first file from each
	file1Header := screenshots1[0]
	file2Header := screenshots2[0]

	// Process Screenshot1
	file1, err := file1Header.Open()
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to open screenshot1", err.Error())
	}
	defer file1.Close()

	img1, _, err := image.Decode(file1)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid image format for screenshot1", err.Error())
	}

	// Process Screenshot2
	file2, err := file2Header.Open()
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to open screenshot2", err.Error())
	}
	defer file2.Close()

	img2, _, err := image.Decode(file2)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid image format for screenshot2", err.Error())
	}

	// Get profile from form values
	profiles := form.Value["profile_fields"]
	if len(profiles) == 0 {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "profile_fields is required", map[string]interface{}{
			"profiles": profiles,
		})
	}
	profileString := profiles[0]

	// Convert profile fields to slice
	profileFields := strings.Split(profileString, ",")

	result, err := h.aiService.AnalyzeVisualDifferences(c.Context(), img1, img2, profileFields)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Could not analyze visual differences", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Analyzed visual difference successfully", result)
}

func (h *AIHandler) SummarizeContentDifferences(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Invalid multipart form", err.Error())
	}

	// Get files from form
	version1Files := form.File["version1"]
	version2Files := form.File["version2"]

	// Validate files
	if len(version1Files) == 0 || len(version2Files) == 0 {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "Both versions are required", map[string]interface{}{
			"version_1": len(version1Files),
			"version_2": len(version2Files),
		})
	}

	// Get first file from each
	file1Header := version1Files[0]
	file2Header := version2Files[0]

	// Process Version 1
	file1, err := file1Header.Open()
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to open version1", err.Error())
	}
	defer file1.Close()

	content1, err := io.ReadAll(file1)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to read version1", err.Error())
	}

	// Process Version 2
	file2, err := file2Header.Open()
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to open version2", err.Error())
	}
	defer file2.Close()

	content2, err := io.ReadAll(file2)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to read version2", err.Error())
	}

	markdownContent1, err := h.processor.Process(string(content1))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to process markdown content 1", err.Error())
	}

	markdownContent2, err := h.processor.Process(string(content2))
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Failed to process markdown content 2", err.Error())
	}

	// Get profile from form values
	profiles := form.Value["profile_fields"]
	if len(profiles) == 0 {
		return sendErrorResponse(c, h.logger, fiber.StatusBadRequest, "profile_fields is required", map[string]interface{}{
			"profiles": profiles,
		})
	}
	profileFieldsString := profiles[0]

	// Convert profile fields to slice
	profileFields := strings.Split(profileFieldsString, ",")

	result, err := h.aiService.AnalyzeContentDifferences(
		c.Context(),
		markdownContent1,
		markdownContent2,
		profileFields,
	)
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Could not analyze content differences", err.Error())
	}

	response, err := h.aiService.SummarizeChanges(c.Context(), []*models.DynamicChanges{result})
	if err != nil {
		return sendErrorResponse(c, h.logger, fiber.StatusInternalServerError, "Could not summarize changes", err.Error())
	}

	return sendDataResponse(c, fiber.StatusOK, "Content analysis successfully", response)
}
