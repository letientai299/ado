package git_prs

import (
	"io"
	"net/url"
	"time"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest/_shared"
)

type SearchCriteria struct {
	// If set, search for pull requests that were created by this identity. Uuid.
	CreatorId *string

	// If specified, filters pull requests that created/closed before this date based on the
	// queryTimeRangeType specified.
	MaxTime *time.Time

	// If specified, filters pull requests that created/closed after this date based on the
	// queryTimeRangeType specified.
	MinTime *time.Time

	// The type of time range which should be used for minTime and maxTime. Defaults to Created if
	// unset.
	QueryTimeRangeType *models.PullRequestTimeRangeType

	// If set, search for pull requests whose target branch is in this repository. Uuid.
	RepositoryId *string

	// If set, search for pull requests that have this identity as a reviewer. Uuid.
	ReviewerId *string

	// If set, search for pull requests from this branch.
	SourceRefName *string

	// If set, search for pull requests whose source branch is in this repository. Uuid.
	SourceRepositoryId *string

	// If set, search for pull requests that are in this state. Defaults to Active if unset.
	Status *models.PullRequestStatus

	// If set, search for pull requests into this branch.
	TargetRefName *string
}

var _ _shared.Querier = (*SearchCriteria)(nil)

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
			[]byte("searchCriteria.queryTimeRangeType=" + url.QueryEscape(string(*s.QueryTimeRangeType))),
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
