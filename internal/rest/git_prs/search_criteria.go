package git_prs

import (
	"io"
	"net/url"
	"time"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest/_shared"
)

// SearchCriteria defines filtering options for listing pull requests.
// All fields are optional; only set fields will be included in the query.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests#gitpullrequestsearchcriteria
type SearchCriteria struct {
	// CreatorId filters pull requests created by this identity (GUID).
	CreatorId *string

	// MaxTime filters pull requests created/closed before this date.
	// The date field used depends on QueryTimeRangeType.
	MaxTime *time.Time

	// MinTime filters pull requests created/closed after this date.
	// The date field used depends on QueryTimeRangeType.
	MinTime *time.Time

	// QueryTimeRangeType specifies which date field to use for MinTime/MaxTime.
	// Defaults to "created" if not set.
	QueryTimeRangeType *models.PullRequestTimeRangeType

	// RepositoryId filters pull requests targeting this repository (GUID).
	RepositoryId *string

	// ReviewerId filters pull requests where this identity is a reviewer (GUID).
	ReviewerId *string

	// SourceRefName filters pull requests from this source branch.
	// Use full ref format: "refs/heads/branch-name".
	SourceRefName *string

	// SourceRepositoryId filters pull requests from this source repository (GUID).
	// Used for cross-repository or fork pull requests.
	SourceRepositoryId *string

	// Status filters pull requests by their lifecycle state.
	// Defaults to "active" if not set.
	Status *models.PullRequestStatus

	// TargetRefName filters pull requests targeting this branch.
	// Use full ref format: "refs/heads/branch-name".
	TargetRefName *string
}

var _ _shared.Querier = (*SearchCriteria)(nil)

// AppendTo writes the search criteria as query parameters to the writer.
func (s *SearchCriteria) AppendTo(w io.Writer) {
	if s == nil {
		return
	}

	if s.CreatorId != nil {
		_, _ = w.Write([]byte("&"))
		_, _ = w.Write([]byte("searchCriteria.creatorId=" + url.QueryEscape(*s.CreatorId)))
	}

	if s.MaxTime != nil {
		_, _ = w.Write([]byte("&"))
		_, _ = w.Write(
			[]byte("searchCriteria.maxTime=" + url.QueryEscape(s.MaxTime.Format(time.RFC3339))),
		)
	}

	if s.MinTime != nil {
		_, _ = w.Write([]byte("&"))
		_, _ = w.Write(
			[]byte("searchCriteria.minTime=" + url.QueryEscape(s.MinTime.Format(time.RFC3339))),
		)
	}

	if s.QueryTimeRangeType != nil {
		_, _ = w.Write([]byte("&"))
		_, _ = w.Write(
			[]byte(
				"searchCriteria.queryTimeRangeType=" + url.QueryEscape(
					string(*s.QueryTimeRangeType),
				),
			),
		)
	}

	if s.RepositoryId != nil {
		_, _ = w.Write([]byte("&"))
		_, _ = w.Write([]byte("searchCriteria.repositoryId=" + url.QueryEscape(*s.RepositoryId)))
	}

	if s.ReviewerId != nil {
		_, _ = w.Write([]byte("&"))
		_, _ = w.Write([]byte("searchCriteria.reviewerId=" + url.QueryEscape(*s.ReviewerId)))
	}

	if s.SourceRefName != nil {
		_, _ = w.Write([]byte("&"))
		_, _ = w.Write([]byte("searchCriteria.sourceRefName=" + url.QueryEscape(*s.SourceRefName)))
	}

	if s.SourceRepositoryId != nil {
		_, _ = w.Write([]byte("&"))
		_, _ = w.Write(
			[]byte("searchCriteria.sourceRepositoryId=" + url.QueryEscape(*s.SourceRepositoryId)),
		)
	}

	if s.Status != nil {
		_, _ = w.Write([]byte("&"))
		_, _ = w.Write([]byte("searchCriteria.status=" + url.QueryEscape(string(*s.Status))))
	}

	if s.TargetRefName != nil {
		_, _ = w.Write([]byte("&"))
		_, _ = w.Write([]byte("searchCriteria.targetRefName=" + url.QueryEscape(*s.TargetRefName)))
	}
}
