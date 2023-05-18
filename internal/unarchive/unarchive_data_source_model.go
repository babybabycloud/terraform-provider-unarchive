package unarchive

import (
	"os"

	"gitee.com/babybabycloud/terraform-provider-unarchive/internal/unarchive/extract"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type unarchiveDataSourceModel struct {
	FileName  types.String `tfsdk:"file_name"`
	Output    types.String `tfsdk:"output"`
	Includes  types.List   `tfsdk:"includes"`
	Excludes  types.List   `tfsdk:"excludes"`
	Flat      types.Bool   `tfsdk:"flat"`
	Type      types.String `tfsdk:"type"`
	FileNames types.List   `tfsdk:"file_names"`
}

func (c unarchiveDataSourceModel) decideOutputDir() string {
	outputDir, err := os.Getwd()
	if err != nil {
		outputDir = "."
	}

	if !c.Output.IsNull() {
		outputDir = c.Output.ValueString()
	}
	return outputDir
}

func (c unarchiveDataSourceModel) includePatterns() extract.TestFunc {
	if !c.Includes.IsNull() {
		patterns := toPatterns(c.Includes)
		return func(name string) bool {
			return patterns.doesNameMatchPatterns(name)
		}
	}
	return func(_ string) bool {
		return true
	}
}

func (c unarchiveDataSourceModel) excludePatterns() extract.TestFunc {
	if !c.Excludes.IsNull() {
		patterns := toPatterns(c.Excludes)
		return func(name string) bool {
			return patterns.doesNameMatchPatterns(name)
		}
	}
	return func(_ string) bool {
		return false
	}
}

func (c unarchiveDataSourceModel) isFlat() bool {
	if c.Flat.IsNull() {
		return false
	}
	return c.Flat.ValueBool()
}

func (z *unarchiveDataSourceModel) addFileNames(filenames []string) {
	values := make([]attr.Value, len(filenames))
	for index, value := range filenames {
		values[index] = types.StringValue(value)
	}

	z.FileNames, _ = types.ListValue(types.StringType, values)
}
