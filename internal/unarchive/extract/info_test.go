package extract

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFailToOpenFile(t *testing.T) {
	name := "any.txt"
	err := errors.New("This is a fake error")

	info := newFailToOpenFile(name, err)
	expected := ExtractInfo{
		Msg: fmt.Sprintf("Failed to open zip file %s", name),
		Err: err,
	}

	assert.EqualValues(t, expected, info)
}
