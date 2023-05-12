package unarchive

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

const (
	DATA_DIR = "data"
)

var (
	FILE_NAME = path.Join("data", "master.zip")
)

func setupDownloadZip() error {
	resp, err := http.Get("https://github.com/babybabycloud/vali-helper/archive/refs/heads/master.zip")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		os.Mkdir(DATA_DIR, DEFAULT_DIR_MODE)
		file, err := os.Create(FILE_NAME)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, resp.Body)
		return err
	}
	return fmt.Errorf("Error with status code %d", resp.StatusCode)
}

func teardownRemoveDownloadZip() {
	os.RemoveAll(DATA_DIR)
}

func TestIsSkipNoFilter(t *testing.T) {
	model := fileDataSourceModel{
		FileName: types.StringValue(FILE_NAME),
	}
	file := NewZipFileWithFileDataSourceModel(model)
	result := file.isSkip("main.py")
	assert.Equal(t, false, result)
}

func TestIsSkipIncludeFilter(t *testing.T) {
	model := fileDataSourceModel{
		FileName: types.StringValue(FILE_NAME),
	}
	includes, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue(`\.py`),
	})
	model.Includes = includes
	file := NewZipFileWithFileDataSourceModel(model)

	result := file.isSkip("main.py")
	assert.Equal(t, false, result)
}

func TestIsSkipExcludeFilter(t *testing.T) {
	model := fileDataSourceModel{
		FileName: types.StringValue(FILE_NAME),
	}
	excludes, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue(`\.py`),
	})
	model.Excludes = excludes
	file := NewZipFileWithFileDataSourceModel(model)
	result := file.isSkip("main.py")
	assert.Equal(t, true, result)
}

func TestCorrectFileNameNotFlat(t *testing.T) {
	model := fileDataSourceModel{
		FileName: types.StringValue(FILE_NAME),
		Output:   types.StringValue(DATA_DIR),
	}
	file := NewZipFileWithFileDataSourceModel(model)
	filename := "a.txt"
	fullname := path.Join("your", "any", "dir", filename)
	result := file.correctFileName(fullname)
	assert.Equal(t, path.Join(DATA_DIR, fullname), result)
}

func TestCorrectFileNameWithFlat(t *testing.T) {
	model := fileDataSourceModel{
		FileName: types.StringValue(FILE_NAME),
		Output:   types.StringValue(DATA_DIR),
		Flat:     types.BoolValue(true),
	}
	file := NewZipFileWithFileDataSourceModel(model)
	filename := "a.txt"
	fullname := path.Join("your", "any", "dir", filename)
	result := file.correctFileName(fullname)
	assert.Equal(t, path.Join(DATA_DIR, filename), result)
}

func TestExtract(t *testing.T) {
	setupDownloadZip()

	model := fileDataSourceModel{
		FileName: types.StringValue(FILE_NAME),
		Flat:     types.BoolValue(true),
		Output:   types.StringValue(DATA_DIR),
	}
	includes, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue(`\.py`),
	})
	model.Includes = includes
	excludes, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue(`tests`),
	})
	model.Excludes = excludes
	file := NewZipFileWithFileDataSourceModel(model)
	result := file.extract()
	expect := extractInfo{
		filenames: []string{
			path.Join(DATA_DIR, "__init__.py"),
			path.Join(DATA_DIR, "helper.py"),
			path.Join(DATA_DIR, "validation.py"),
		},
	}
	assert.EqualValues(t, expect, result)

	teardownRemoveDownloadZip()
}
