package extract

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const DEFAULT_DIR_MODE = 0740

func copy(info *item) error {
	w, err := createFile(info)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, info.copyItem)
	if err != nil {
		return err
	}
	return nil
}

func createFile(info *item) (io.WriteCloser, error) {
	filename := path.Dir(info.name)
	err := os.MkdirAll(filename, DEFAULT_DIR_MODE)
	if err != nil {
		return nil, err
	}

	cf, err := os.OpenFile(info.name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.FileMode(info.mode))
	if err != nil {
		return nil, err
	}

	return cf, nil
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

	filenames := handler.generate(conf, func(itemInfo *item) string {
		if conf.isSkip(itemInfo.name) || !itemInfo.isRegFile {
			tflog.Info(conf.Ctx, fmt.Sprintf("Skip %s", itemInfo.name))
			return ""
		}

		filename := conf.correctFileName(itemInfo.name)
		itemInfo.name = filename
		err = copy(itemInfo)
		if err != nil {
			tflog.Error(conf.Ctx, err.Error())
			return ""
		}

		// Close the CopyItem if it can be closed
		if canClose, ok := itemInfo.copyItem.(io.Closer); ok {
			canClose.Close()
		}
		return filename
	})

	return ExtractInfo{
		FileNames: filenames,
	}
}
