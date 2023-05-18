package extract

import (
	"context"
	"path/filepath"
)

// TestFunc is a function with a string parameter to check includes and excludes
type TestFunc func(string) bool

// Config is an internal form of unarchive data source model
type Config struct {
	Ctx     context.Context
	Name    string
	Include TestFunc
	Exclude TestFunc
	Outdir  string
	IsFlat  bool
	Type    string
}

func (c *Config) isSkip(name string) bool {
	var isSkip bool
	return isSkip || !c.Include(name) || c.Exclude(name)
}

func (c *Config) correctFileName(filename string) string {
	var dir string
	if c.IsFlat {
		dir = filepath.Join(c.Outdir, filepath.Base(filename))
	} else {
		dir = filepath.Join(c.Outdir, filename)
	}
	return dir
}