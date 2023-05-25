package extract

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargzGenerate(t *testing.T) {
	name := "/tmp/targzfile.tar.gz"
	createTestTargzFile(t, name)
	h := &targzHandler{}
	_ = h.open(name)
	defer h.close()
	conf := &Config{
		Ctx: context.TODO(),
	}
	filenames := h.generate(conf, func(i *item) string {
		return i.name
	})
	defer os.Remove(name)

	assert.NotEmpty(t, filenames)
}

func createTestTargzFile(t *testing.T, name string) {
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()

	content := []byte("Test content")

	gw := gzip.NewWriter(file)
	defer gw.Close()
	writter := tar.NewWriter(gw)
	defer writter.Close()

	header := &tar.Header{
		Name: "test/include.txt",
		Mode: 600,
		Size: int64(len(content)),
	}
	_ = writter.WriteHeader(header)
	_, _ = writter.Write(content)
	writter.Flush()

	header = &tar.Header{
		Name: "test/exclude.txt",
		Mode: 600,
		Size: int64(len(content)),
	}
	_ = writter.WriteHeader(header)
	_, _ = writter.Write(content)
	writter.Flush()
}
