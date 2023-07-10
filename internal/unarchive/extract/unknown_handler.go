package extract

type unknownHandler struct{}

func (u unknownHandler) open(name string) error {
	panic("Unsupport type of handler")
}

func (u unknownHandler) generate(*Config, testAndCopy) []string {
	panic("Unsupport type of handler")
}

func (u unknownHandler) close() {
	panic("Unsupport type of handler")
}
