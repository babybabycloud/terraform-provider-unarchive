package extract

import (
	"archive/tar"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTarGenerate(t *testing.T) {
	name := "/tmp/tarfile.tar"
	createTestTarFile(t, name)
	h := &tarHandler{}
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

func createTestTarFile(t *testing.T, name string) {
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()

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
