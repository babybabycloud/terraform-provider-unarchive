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
	h.open(name)
	defer h.close()
	conf := &Config{
		Ctx: context.TODO(),
	}
	ch := h.generate(conf)

	items := make([]*item, 0)
	for element := range ch {
		items = append(items, element)
	}

	assert.NotEmpty(t, items)
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
