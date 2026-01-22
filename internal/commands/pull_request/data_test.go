package pull_request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEditPrInfo_EmptyTitle(t *testing.T) {
	// We use a simple command that empties the file to simulate empty title
	info := &prInfo{title: "original", desc: "desc"}
	// "truncate -s 0" or ">" or "cp /dev/null"
	// On macOS "truncate" might not be available, but ">" should work in "sh -c"
	// However, editPrInfo uses editor.New which runs Open(cmd, tmpFile.Name())
	// which runs sh -c 'cmd "$1"' -- filePath.
	// So if cmd is ">", it becomes sh -c '> "$1"' -- filePath, which empties the file.

	newInfo, err := editPrInfo(info, ">")
	assert.Error(t, err)
	assert.Equal(t, ErrEmptyTitle, err)
	assert.Nil(t, newInfo)
}

func TestEditPrInfo_Success(t *testing.T) {
	info := &prInfo{title: "original", desc: "desc"}
	// Replace content with "new title\n\nnew desc"
	// Using printf to handle newlines
	newInfo, err := editPrInfo(info, "printf 'new title\n\nnew desc' >")
	assert.NoError(t, err)
	assert.Equal(t, "new title", newInfo.title)
	assert.Equal(t, "new desc", newInfo.desc)
}
