package extract

import (
	"io"
)

type handler interface {
	open(name string) error
	generate(*Config) <-chan *item
	close()
}

type item struct {
	copyItem  io.Reader
	name      string
	isRegFile bool
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
		// Add unknown handler
		return nil
	}
}
