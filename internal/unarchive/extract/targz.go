package extract

import (
	"compress/gzip"
	"io"
)

type targzHandler struct {
	th tarHandler
}

func (t *targzHandler) open(name string) error {
	err := t.th.open(name)
	if err != nil {
		return err
	}

	gr, err := gzip.NewReader(t.th.rc.Back().Value.(io.Reader))
	if err != nil {
		return err
	}
	t.th.rc.PushBack(gr)
	return nil
}

func (t *targzHandler) generate(conf *Config) <-chan *item {
	return t.th.generate(conf)
}

func (t *targzHandler) close() {
	t.th.close()
}
