package extract

import (
	"context"
	"path/filepath"

	"github.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/common"
	"github.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/model"
)

// ConfigFilter stores the expressions of include and exclude
type ConfigFilter struct {
	Includes common.Patterns
	Excludes common.Patterns
}

// FormFilterModel constructs an instance of ConfigFilter using model.FilterModel
func FromFilterModel(f model.FilterModel) ConfigFilter {
	var filter ConfigFilter
	if !f.Includes.IsNull() {
		filter.Includes = common.ToPatterns(f.Includes)
	}

	if !f.Excludes.IsNull() {
		filter.Excludes = common.ToPatterns(f.Excludes)
	}
	return filter
}

func (f ConfigFilter) isSkip(name string) bool {
	return !f.Includes.Include(name) || f.Excludes.Exclude(name)
}

// Config is an internal form of unarchive data source model
type Config struct {
	Ctx     context.Context
	Name    string
	Filters []ConfigFilter
	Outdir  string
	IsFlat  bool
	Type    string
}

// FromUnarchiveDataSourceModel creates a pointer refering to an instanc of Config
func FromUnarchiveDataSourceModel(m model.UnarchiveDataSourceModel) *Config {
	conf := new(Config)
	conf.Name = m.FileName.ValueString()
	conf.Outdir = m.DecideOutputDir()
	conf.IsFlat = m.IsFlat()
	conf.Type = m.Type.ValueString()

	filters := make([]ConfigFilter, len(m.Filters))
	for index, filter := range m.Filters {
		filters[index] = FromFilterModel(filter)
	}
	conf.Filters = filters
	return conf
}

func (c *Config) isSkip(name string) bool {
	for _, filter := range c.Filters {
		if !filter.isSkip(name) {
			return false
		}
	}
	return len(c.Filters) != 0
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
