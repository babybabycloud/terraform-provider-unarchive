package extract

import (
	"archive/tar"
	"container/list"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type tarHandler struct {
	rc *list.List
}

func (t *tarHandler) open(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	t.rc = list.New()
	t.rc.PushBack(file)
	return nil
}

func (t *tarHandler) generate(conf *Config) <-chan *item {
	r := t.rc.Back().Value.(io.Reader)
	reader := tar.NewReader(r)
	ch := make(chan *item)
	go func() {
		defer close(ch)
		for {
			header, err := reader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				tflog.Error(conf.Ctx, fmt.Sprintf("%s, continue", err))
				continue
			}

			ch <- &item{
				copyItem:  reader,
				name:      header.Name,
				isRegFile: header.Typeflag == tar.TypeReg,
			}
		}
	}()
	return ch
}

func (t *tarHandler) close() {
	for e := t.rc.Back(); e != nil; e = e.Prev() {
		e.Value.(io.Closer).Close()
	}
}
