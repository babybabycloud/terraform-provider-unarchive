package extract

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const DEFAULT_DIR_MODE = 0740

func copy(name string, r io.Reader) error {
	w, err := createFile(name)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

func createFile(name string) (io.WriteCloser, error) {
	filename := path.Dir(name)
	err := os.MkdirAll(filename, DEFAULT_DIR_MODE)
	if err != nil {
		return nil, err
	}

	cf, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	return cf, err
}

// Extract is the main function to extract files from an archive file
func Extract(conf *Config) ExtractInfo {
	err := os.MkdirAll(conf.Outdir, DEFAULT_DIR_MODE)
	if err != nil {
		return ExtractInfo{
			Msg: fmt.Sprintf("Failed to create directory %s", conf.Outdir),
			Err: err,
		}
	}

	handler := getHandler(conf.Type)

	err = handler.open(conf.Name)
	if err != nil {
		return newFailToOpenFile(conf.Name, err)
	}
	defer handler.close()

	ch := handler.generate(conf)

	fileNames := make([]string, 0)
	for item := range ch {
		if conf.isSkip(item.name) || !item.isRegFile {
			tflog.Info(conf.Ctx, fmt.Sprintf("Skip %s", item.name))
			continue
		}

		filename := conf.correctFileName(item.name)
		err = copy(filename, item.copyItem)
		if err != nil {
			tflog.Error(conf.Ctx, err.Error())
			continue
		}

		fileNames = append(fileNames, filename)
		// Close the CopyItem if it can be closed
		if canClose, ok := item.copyItem.(io.Closer); ok {
			canClose.Close()
		}
	}
	return ExtractInfo{
		FileNames: fileNames,
	}
}
