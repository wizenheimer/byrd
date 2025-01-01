package handlers

import (
	"io"

	"github.com/gofiber/fiber/v2"
	interfaces "github.com/wizenheimer/iris/src/internal/interfaces/service"
	api_models "github.com/wizenheimer/iris/src/internal/models/api"
	core_models "github.com/wizenheimer/iris/src/internal/models/core"
	"github.com/wizenheimer/iris/src/pkg/logger"
	"github.com/wizenheimer/iris/src/pkg/utils/path"
)

type DiffHandler struct {
	diffService interfaces.DiffService
	logger      *logger.Logger
}

func NewDiffHandler(diffService interfaces.DiffService, logger *logger.Logger) *DiffHandler {
	logger.Debug("creating new diff handler")

	return &DiffHandler{
		diffService: diffService,
		logger:      logger.WithFields(map[string]interface{}{"module": "diff_handler"}),
	}
}

// Get retrieves the diff analysis for a URL
func (h *DiffHandler) Get(c *fiber.Ctx) error {
	h.logger.Debug("creating new diff")

	var req api_models.URLDiffRequest
	var err error
	if err = c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Normalize the URL
	req.URL, err = path.NormalizeURL(req.URL)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.diffService.Get(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

// Compare compares the content of two versions of a URL
func (h *DiffHandler) Compare(c *fiber.Ctx) error {
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
	profiles := form.Value["profile"]
	if len(profiles) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Profile is required")
	}
	profile := profiles[0]

	// Compare contents
	htmlContent1 := core_models.ScreenshotHTMLContentResponse{
		Status:      "success",
		HTMLContent: string(content1),
	}

	htmlContent2 := core_models.ScreenshotHTMLContentResponse{
		Status:      "success",
		HTMLContent: string(content2),
	}

	result, err := h.diffService.Compare(c.Context(), &htmlContent1, &htmlContent2, profile, false)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}
