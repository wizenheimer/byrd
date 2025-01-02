package handlers

import "github.com/gofiber/fiber/v2"

type WorkspaceHandler struct {
	// TODO: populate the skeleton
}

func NewWorkspaceHandler() *WorkspaceHandler {
	return &WorkspaceHandler{}
}

// CreateWorkspace creates a new workspace
func (wh *WorkspaceHandler) CreateWorkspace(c *fiber.Ctx) error {
	return nil
}

// ListWorkspaces lists workspaces for a user
func (wh *WorkspaceHandler) ListWorkspaces(c *fiber.Ctx) error {
	return nil
}

// GetWorkspace gets a workspace by ID
func (wh *WorkspaceHandler) GetWorkspace(c *fiber.Ctx) error {
	return nil
}

// UpdateWorkspace updates a workspace by ID
func (wh *WorkspaceHandler) UpdateWorkspace(c *fiber.Ctx) error {
	return nil
}

// DeleteWorkspace deletes a workspace by ID
func (wh *WorkspaceHandler) DeleteWorkspace(c *fiber.Ctx) error {
	return nil
}

// ExitWorkspace exits a workspace by ID
func (wh *WorkspaceHandler) ExitWorkspace(c *fiber.Ctx) error {
	return nil
}

// JoinWorkspace joins a workspace by ID
func (wh *WorkspaceHandler) JoinWorkspace(c *fiber.Ctx) error {
	return nil
}

// ListWorkspaceUsers lists users for a workspace
func (wh *WorkspaceHandler) ListWorkspaceUsers(c *fiber.Ctx) error {
	return nil
}

// AddUserToWorkspace adds a user to a workspace
func (wh *WorkspaceHandler) AddUserToWorkspace(c *fiber.Ctx) error {
	return nil
}

// RemoveUserFromWorkspace removes a user from a workspace
func (wh *WorkspaceHandler) RemoveUserFromWorkspace(c *fiber.Ctx) error {
	return nil
}

// UpdateUserRoleInWorkspace updates user role in a workspace
func (wh *WorkspaceHandler) UpdateUserRoleInWorkspace(c *fiber.Ctx) error {
	return nil
}

// CreateCompetitorForWorkspace creates a competitor for a workspace
func (wh *WorkspaceHandler) CreateCompetitorForWorkspace(c *fiber.Ctx) error {
	return nil
}

// AddPageToCompetitor adds a page to a competitor
func (wh *WorkspaceHandler) AddPageToCompetitor(c *fiber.Ctx) error {
	return nil
}

// ListWorkspaceCompetitors lists competitors for a workspace
func (wh *WorkspaceHandler) ListWorkspaceCompetitors(c *fiber.Ctx) error {
	return nil
}

// ListPageHistory lists page history
func (wh *WorkspaceHandler) ListPageHistory(c *fiber.Ctx) error {
	return nil
}

// RemovePageFromCompetitor removes a page from a competitor
func (wh *WorkspaceHandler) RemovePageFromCompetitor(c *fiber.Ctx) error {
	return nil
}

// RemoveCompetitorFromWorkspace removes a competitor from a workspace
func (wh *WorkspaceHandler) RemoveCompetitorFromWorkspace(c *fiber.Ctx) error {
	return nil
}

// UpdatePageInCompetitor updates a page in a competitor
func (wh *WorkspaceHandler) UpdatePageInCompetitor(c *fiber.Ctx) error {
	return nil
}
