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

func (t *targzHandler) generate(conf *Config, f testAndCopy) []string {
	return t.th.generate(conf, f)
}

func (t *targzHandler) close() {
	t.th.close()
}
