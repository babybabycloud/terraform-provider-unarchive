package model

import (
	"os"

	"gitee.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/common"
	"gitee.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/extract"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UnarchiveDataSourceModel is unarchive data source model
type UnarchiveDataSourceModel struct {
	FileName  types.String `tfsdk:"file_name"`
	Output    types.String `tfsdk:"output"`
	Includes  types.List   `tfsdk:"includes"`
	Excludes  types.List   `tfsdk:"excludes"`
	Flat      types.Bool   `tfsdk:"flat"`
	Type      types.String `tfsdk:"type"`
	FileNames types.List   `tfsdk:"file_names"`
}

// DecideOutputDir gets the final output directory
func (c UnarchiveDataSourceModel) DecideOutputDir() string {
	outputDir, err := os.Getwd()
	if err != nil {
		outputDir = "."
	}

	if !c.Output.IsNull() {
		outputDir = c.Output.ValueString()
	}
	return outputDir
}

// IncludePatterns return a function helps to filter which file is included
func (c UnarchiveDataSourceModel) IncludePatterns() extract.TestFunc {
	if !c.Includes.IsNull() {
		patterns := common.ToPatterns(c.Includes)
		return func(name string) bool {
			return patterns.DoesNameMatch(name)
		}
	}
	return func(_ string) bool {
		return true
	}
}

// ExcludePatterns return a function helps to filter which file is excluded
func (c UnarchiveDataSourceModel) ExcludePatterns() extract.TestFunc {
	if !c.Excludes.IsNull() {
		patterns := common.ToPatterns(c.Excludes)
		return func(name string) bool {
			return patterns.DoesNameMatch(name)
		}
	}
	return func(_ string) bool {
		return false
	}
}

// IsFlat is if the output file should be flatted
func (c UnarchiveDataSourceModel) IsFlat() bool {
	if c.Flat.IsNull() {
		return false
	}
	return c.Flat.ValueBool()
}

// AddFileNames adds all the output file names to the result
func (z *UnarchiveDataSourceModel) AddFileNames(filenames []string) {
	values := make([]attr.Value, len(filenames))
	for index, value := range filenames {
		values[index] = types.StringValue(value)
	}

	z.FileNames, _ = types.ListValue(types.StringType, values)
}
