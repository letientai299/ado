package pull_request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrInfoParse(t *testing.T) {
	tests := []struct {
		name      string
		info      prInfo
		content   string
		wantTitle string
		wantDesc  string
		wantErr   error
	}{
		{
			name:    "empty title",
			info:    prInfo{},
			content: "\n\n",
			wantErr: ErrEmptyTitle,
		},
		{
			name: "unchanged on existing",
			info: prInfo{
				title: "old title",
				desc:  "old desc",
			},
			content: "old title\n\nold desc",
			wantErr: ErrUnchanged,
		},
		{
			name: "unchanged on new allows",
			info: prInfo{
				title: "old title",
				desc:  "old desc",
				isNew: true,
			},
			content:   "old title\n\nold desc",
			wantTitle: "old title",
			wantDesc:  "old desc",
		},
		{
			name: "updates with marker ignored",
			info: prInfo{
				title: "old title",
				desc:  "old desc",
			},
			content:   "new title\n\nnew desc\n" + prEditingMarker + "\nignored",
			wantTitle: "new title",
			wantDesc:  "new desc",
		},
		{
			name:      "title only",
			info:      prInfo{},
			content:   "new title",
			wantTitle: "new title",
			wantDesc:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.info.parse(tt.content)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantTitle, tt.info.title)
			require.Equal(t, tt.wantDesc, tt.info.desc)
		})
	}
}
