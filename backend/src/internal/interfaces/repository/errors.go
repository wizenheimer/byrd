package interfaces

import "errors"

// ------------- General Errors -------------
var (
	ErrInvalidLimit  = errors.New("invalid limit value for list operation")
	ErrInvalidOffset = errors.New("invalid offset value for list operation")
)

// ------------- Page History Repository Errors -------------
var (
	// ---- validation errors ----
	ErrPageIDsUnspecified = errors.New("pageIDs unspecified")

	// ---- remapped errors ----
	// case 1: remapping an existing error
	ErrFailedToScanPageHistory = errors.New("failed to scan page history")
	ErrInvalidPageHistory      = errors.New("invalid page history")

	// case 2: remapping a non error scenario to an error
	ErrPageHistoryNotFound = errors.New("no page history found for the page")
)

// ------------- Competitor Repository Errors -------------
var (
	// ---- validation errors ----
	ErrCompetitorNamesListEmpty = errors.New("competitor names list is empty")

	// --- remapped errors ----
	// remapping an existing error
	// remapping non error scenario to an error
	ErrCompetitorNotFound               = errors.New("competitor not found")
	ErrNoWorkspaceCompetitorsFound      = errors.New("no competitors found for the workspace")
	ErrFailedToScanCompetitors          = errors.New("failed to scan competitors")
	ErrFailedToConfirmCompetitorRemoval = errors.New("failed to confirm competitor removal")
)

// ------------- Page Repository Errors -------------
var (
	// ---- validation errors -----
	ErrPagesUnspecified                = errors.New("pages unspecified for competitor")
	ErrFailedToMarshallCaptureProfile  = errors.New("failed to marshal capture profile for page")
	ErrFailedToUnmarshalCaptureProfile = errors.New("failed to unmarshal capture profile for page")
	ErrFailedToMarshallDiffProfile     = errors.New("failed to marshal diff profile for page")
	ErrFailedToUnmarshalDiffProfile    = errors.New("failed to unmarshal diff profile for page")
	ErrInvalidBatchSize                = errors.New("invalid batch size")

	// ---- non fatal errors ----
	ErrFailedToConfirmPageRemoval = errors.New("failed to confirm page removal")

	// ---- remapped errors ----
	// case 1 : remapping an existing error
	// case 2 : remapping a non error scenario to an error
	ErrNoCompetitorPages               = errors.New("no pages found for the competitor")
	ErrPageNotFound                    = errors.New("page not found")
	ErrFailedToScanPages               = errors.New("failed to scan pages from pages table")
	ErrFailedToIterateOverPagesForScan = errors.New("failed to scan pages from pages table")
)

// ------------- User Repository Errors -------------

var (
	// ---- validation errors -----
	ErrInvalidUserEmail          = errors.New("user email invalid")
	ErrUserEmailsNotSpecified    = errors.New("user emails unspecified")
	ErrUserIDsNotSpecified       = errors.New("user ids unspecified")
	ErrCouldnotGetClerkUserEmail = errors.New("couldn't get user email from profile")

	// ---- non fatal errors ----
	ErrCouldnotScanUser            = errors.New("couldn't scan user")
	ErrNoWorkspaceFoundForUser     = errors.New("no workspaces found for the user")
	ErrCouldnotConfirmSyncStatus   = errors.New("couldn't confirm sync status")
	ErrCouldnotConfirmDeleteStatus = errors.New("couldn't confirm delete status")

	// ---- remapped errors ----
	// case 1 : remapping an existing error
	// case 2 : remapping a non error scenario to an error
	ErrUserNotFoundByID        = errors.New("user not found by ID")
	ErrUserNotFoundByEmail     = errors.New("user not found by email")
	ErrUserNotFoundByIDOrEmail = errors.New("user not found by ID or Email")
	ErrWorkspaceUsersNotFound  = errors.New("users not found in workspace")
)

// ------------- Workspace Repository Errors -------------

var (
	// ---- validation errors -----
	ErrNoWorkspaceSpecified = errors.New("no workspace specified")

	// ---- non fatal errors ----
	ErrFailedToConfirmUpdatedBillingEmail   = errors.New("failed to confirm billing email")
	ErrFailedToConfirmWorkspaceNameUpdate   = errors.New("failed to confirm workspace name update")
	ErrFailedToConfirmWorkspaceStatusUpdate = errors.New("failed to confirm workspace status update")
	ErrFailedToConfirmWorkspaceUpdate       = errors.New("failed to confirm workspace update")

	// ---- remapped errors ----
	// case 1 : remapping an existing error
	// case 2 : remapping a non error scenario to an error
	ErrWorkspaceNotFound = errors.New("workspace not found")
)
