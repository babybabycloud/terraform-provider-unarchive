package unarchive

type extractInfo struct {
	filenames []string
	msg       string
	err       error
}

func NewExtractInfoFromNameChan(ch <-chan string) extractInfo {
	filenames := make([]string, 0)
	for name := range ch {
		filenames = append(filenames, name)
	}
	return extractInfo{
		filenames: filenames,
	}
}
