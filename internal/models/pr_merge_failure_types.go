package models

// PullRequestMergeFailureType represents the type of failure that occurred
// during a pull request merge operation.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get#pullrequestmergefailuretype
type PullRequestMergeFailureType string

const (
	// PullRequestMergeFailureTypeNone indicates no failure occurred.
	// This is the default value when the merge has not failed.
	PullRequestMergeFailureTypeNone PullRequestMergeFailureType = "none"

	// PullRequestMergeFailureTypeUnknown indicates an unknown failure type.
	// The merge failed for an unspecified reason.
	PullRequestMergeFailureTypeUnknown PullRequestMergeFailureType = "unknown"

	// PullRequestMergeFailureTypeCaseSensitive indicates a case sensitivity conflict.
	// Files differ only in case (e.g., "File.txt" vs "file.txt"), which can cause
	// issues on case-insensitive file systems like Windows.
	PullRequestMergeFailureTypeCaseSensitive PullRequestMergeFailureType = "caseSensitive"

	// PullRequestMergeFailureTypeObjectTooLarge indicates an object exceeded size limits.
	// A file or object in the merge is too large for the repository.
	PullRequestMergeFailureTypeObjectTooLarge PullRequestMergeFailureType = "objectTooLarge"
)
