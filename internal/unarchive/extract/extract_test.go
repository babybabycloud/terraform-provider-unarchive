package extract

import (
	"archive/tar"
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	name := "/tmp/src.txt"
	var reader bytes.Buffer
	content := []byte("Test file")
	reader.Write(content)

	err := copy(name, &reader)
	assert.Nil(t, err)
	defer os.Remove(name)

	result, err := ioutil.ReadFile(name)
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
		Include: func(name string) bool {
			return strings.Contains(name, "include")
		},
		Exclude: func(name string) bool {
			return strings.Contains(name, "exclude")
		},
	}

	result := Extract(conf)
	expectFileName := "/tmp/include.txt"
	expect := []string{expectFileName}
	defer os.Remove(tarName)
	defer os.Remove(expectFileName)

	assert.EqualValues(t, expect, result.FileNames)
}

func createTestTarFile(t *testing.T, name string) {
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	content := []byte("Test content")

	writter := tar.NewWriter(file)
	defer writter.Close()

	header := &tar.Header{
		Name: "test/include.txt",
		Mode: 600,
		Size: int64(len(content)),
	}
	writter.WriteHeader(header)
	writter.Write(content)
	writter.Flush()

	header = &tar.Header{
		Name: "test/exclude.txt",
		Mode: 600,
		Size: int64(len(content)),
	}
	writter.WriteHeader(header)
	writter.Write(content)
	writter.Flush()
}
