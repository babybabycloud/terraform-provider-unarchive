package extract

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"testing"

	"github.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/common"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	name := "/tmp/des.txt"
	var reader bytes.Buffer
	content := []byte("Test file")
	reader.Write(content)

	itemInfo := &item{
		copyItem: &reader,
		name:     name,
		mode:     int64(fs.ModePerm),
	}
	err := copy(itemInfo)
	assert.Nil(t, err)
	defer os.Remove(name)

	result, err := os.ReadFile(name)
	assert.Nil(t, err)

	assert.Equal(t, content, result)
}

func TestExtract(t *testing.T) {
	tarName := "/tmp/test.tar"
	createTestTarFile(t, tarName)

	conf := &Config{
		Ctx:    context.TODO(),
		Name:   tarName,
		Type:   TAR,
		IsFlat: true,
		Outdir: "/tmp",
		Filters: []ConfigFilter{
			{
				Includes: common.Patterns{"include"},
				Excludes: common.Patterns{"exclude"},
			},
		},
	}

	result := Extract(conf)
	expectFileName := "/tmp/include.txt"
	expect := []string{expectFileName}
	defer os.Remove(tarName)
	defer os.Remove(expectFileName)

	assert.EqualValues(t, expect, result.FileNames)
}
