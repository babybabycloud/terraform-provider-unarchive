package extract

import (
	"archive/zip"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type zipHandler struct {
	rc *zip.ReadCloser
}

func (z *zipHandler) open(name string) error {
	rc, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	z.rc = rc
	return nil
}

func (z *zipHandler) generate(conf *Config) <-chan *item {
	ch := make(chan *item)
	go func() {
		defer close(ch)
		for _, zipFile := range z.rc.File {
			file, err := zipFile.Open()
			if err != nil {
				tflog.Error(conf.Ctx, err.Error())
				continue
			}

			ch <- &item{
				copyItem:  file,
				name:      zipFile.Name,
				isRegFile: zipFile.Method == zip.Deflate,
				mode:      int64(zipFile.Mode()),
			}
		}
	}()
	return ch
}

func (z *zipHandler) close() {
	z.rc.Close()
}
