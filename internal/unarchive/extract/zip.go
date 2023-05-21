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

func (z *zipHandler) generate(conf *Config, f testAndCopy) []string {
	filenames := make([]string, 0)
	for _, zipFile := range z.rc.File {
		file, err := zipFile.Open()
		if err != nil {
			tflog.Error(conf.Ctx, err.Error())
			continue
		}

		filename := f(&item{
			copyItem:  file,
			name:      zipFile.Name,
			isRegFile: zipFile.Method == zip.Deflate,
			mode:      int64(zipFile.Mode()),
		})
		if filename != "" {
			filenames = append(filenames, filename)
		}
	}
	return filenames
}

func (z *zipHandler) close() {
	z.rc.Close()
}
