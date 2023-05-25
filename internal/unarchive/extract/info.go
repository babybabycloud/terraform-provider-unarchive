package extract

import "fmt"

// ExtractInfo is the result of the extraction
type ExtractInfo struct {
	FileNames []string
	Msg       string
	Err       error
}

func newFailToOpenFile(name string, err error) ExtractInfo {
	return ExtractInfo{
		Msg: fmt.Sprintf("Failed to open zip file %s", name),
		Err: err,
	}
}
