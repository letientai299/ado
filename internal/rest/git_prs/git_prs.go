package git_prs

import (
	"io"
	"strconv"
)

// ListQuery defines pagination and filtering options for listing pull requests.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git/pull-requests/get-pull-requests
type ListQuery struct {
	// Skip is the number of pull requests to skip (for pagination).
	Skip *int

	// Top is the maximum number of pull requests to return.
	// Azure DevOps has a default limit; use this to request fewer results.
	Top *int

	// SearchCriteria contains optional filters for the pull request list.
	SearchCriteria *SearchCriteria
}

// AppendTo writes the list query parameters to the writer.
func (q ListQuery) AppendTo(writer io.Writer) {
	if q.Top != nil {
		_, _ = writer.Write([]byte("&$top=" + strconv.Itoa(*q.Top)))
	}

	if q.Skip != nil {
		_, _ = writer.Write([]byte("&$skip=" + strconv.Itoa(*q.Skip)))
	}

	q.SearchCriteria.AppendTo(writer)
}
