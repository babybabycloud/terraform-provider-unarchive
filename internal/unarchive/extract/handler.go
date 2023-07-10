package extract

import (
	"io"
)

type testAndCopy func(*item) string

type handler interface {
	open(name string) error
	generate(*Config, testAndCopy) []string
	close()
}

type item struct {
	copyItem  io.Reader
	name      string
	isRegFile bool
	mode      int64
}

const (
	ZIP   = ".zip"
	TAR   = ".tar"
	TARGZ = ".tar.gz"
)

func getHandler(hType string) handler {
	switch hType {
	case ZIP:
		return &zipHandler{}
	case TAR:
		return &tarHandler{}
	case TARGZ:
		return &targzHandler{}
	default:
		return &unknownHandler{}
	}
}
