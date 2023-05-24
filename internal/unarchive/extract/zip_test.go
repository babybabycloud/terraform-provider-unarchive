package extract

import (
	"archive/zip"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZipGenerate(t *testing.T) {
	name := "/tmp/zip.zip"
	createZipFile(t, name)
	z := &zipHandler{}
	_ = z.open(name)
	defer z.close()

	conf := &Config{
		Ctx: context.TODO(),
	}
	filenames := z.generate(conf, func(i *item) string {
		return i.name
	})
	defer os.Remove(name)

	assert.NotEmpty(t, filenames)

}

func createZipFile(t *testing.T, name string) {
	file, err := os.Create(name)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer file.Close()

	content := []byte("Test Content")
	writter := zip.NewWriter(file)
	w, err := writter.Create("test/include.txt")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer writter.Close()

	_, err = w.Write(content)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	w, err = writter.Create("test/exclude.txt")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, _ = w.Write(content)
	writter.Flush()
}
