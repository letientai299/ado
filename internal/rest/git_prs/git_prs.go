package git_prs

import (
	"io"
	"strconv"
)

type ListQuery struct {
	Skip           *int
	Top            *int
	SearchCriteria *SearchCriteria
}

func (q ListQuery) AppendTo(writer io.Writer) {
	if q.Top != nil {
		_, _ = writer.Write([]byte("&$top=" + strconv.Itoa(*q.Top)))
	}

	if q.Skip != nil {
		_, _ = writer.Write([]byte("&$skip=" + strconv.Itoa(*q.Skip)))
	}

	q.SearchCriteria.AppendTo(writer)
}
