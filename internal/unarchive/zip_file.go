package unarchive

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type zipFile struct {
	ctx     context.Context
	name    string
	include testFunc
	exclude testFunc
	outdir  string
	isFlat  bool
}

func NewZipFileWithFileDataSourceModel(model fileDataSourceModel) *zipFile {
	file := new(zipFile)
	file.ctx = context.TODO()
	file.name = model.FileName.ValueString()
	file.include = model.includePatterns()
	file.exclude = model.excludePatterns()
	file.outdir = model.decideOutputDir()
	file.isFlat = model.isFlat()
	return file
}

func (zf *zipFile) extract() extractInfo {
	rc, err := zip.OpenReader(zf.name)
	if err != nil {
		return extractInfo{
			msg: fmt.Sprintf("Failed to open zip file %s\n", zf.name),
			err: err,
		}
	}

	err = os.MkdirAll(zf.outdir, DEFAULT_DIR_MODE)
	if err != nil {
		return extractInfo{
			msg: fmt.Sprintf("Failed to create directory %s\n", zf.outdir),
			err: err,
		}
	}

	nameChan := make(chan string)
	go func() {
		defer close(nameChan)
		for _, file := range rc.File {

			if zf.isSkip(file.Name) {
				tflog.Debug(zf.ctx, fmt.Sprintf("Skip %s\n", file.Name))
				continue
			}

			filename := zf.correctFileName(file.Name)
			file.Name = filename
			// Directory
			if file.Method == zip.Store {
				continue
			}
			err = zf.copyFile(file)
			if err != nil {
				tflog.Error(zf.ctx, fmt.Sprintf("Failed to copy file %s\n", file.Name))
				continue
			}
			nameChan <- file.Name
		}
	}()

	return NewExtractInfoFromNameChan(nameChan)
}

func (zf *zipFile) isSkip(name string) bool {
	var isSkip bool
	return isSkip || !zf.include(name) || zf.exclude(name)
}

func (zf *zipFile) correctFileName(filename string) string {
	var dir string
	if zf.isFlat {
		dir = filepath.Join(zf.outdir, filepath.Base(filename))
	} else {
		dir = filepath.Join(zf.outdir, filename)
	}
	return dir
}

func (zf *zipFile) copyFile(file *zip.File) error {
	filename := path.Dir(file.Name)
	err := os.MkdirAll(filename, DEFAULT_DIR_MODE)
	if err != nil {
		return err
	}
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	cf, err := os.Create(file.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(cf, rc)
	if err != nil {
		return err
	}
	return nil
}
