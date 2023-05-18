package extract

import "io"

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
