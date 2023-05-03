package unarchive

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type config struct {
	FileName types.String `tfsdk:"file_name"`
	Output   types.String `tfsdk:"output"`
	Includes types.List   `tfsdk:"includes"`
	Excludes types.List   `tfsdk:"excludes"`
}

func (c config) decideOutputDir() string {
	outputDir, err := os.Getwd()
	if err != nil {
		outputDir = "./"
	}

	if !c.Output.IsNull() && !c.Output.IsUnknown() {
		outputDir = c.Output.ValueString()
	}
	return outputDir
}

func (z zipFileDataSourceModel) toPatterns(list types.List) patterns {
	var p patterns
	if !list.IsNull() && !list.IsUnknown() {
		patternsInType := make([]types.String, len(list.Elements()))
		p = make(patterns, len(list.Elements()))
		for index, value := range patternsInType {
			p[index] = value.ValueString()
		}
	}
	return p
}
