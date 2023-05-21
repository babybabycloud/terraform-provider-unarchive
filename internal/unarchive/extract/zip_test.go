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
	z.open(name)
	defer z.close()

	conf := &Config{
		Ctx: context.TODO(),
	}
	ch := z.generate(conf)

	items := make([]*item, 2)
	for element := range ch {
		items = append(items, element)
	}
	defer os.Remove(name)

	assert.NotEmpty(t, items)

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
	_, err = w.Write(content)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	writter.Flush()
}
