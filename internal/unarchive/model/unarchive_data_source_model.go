package model

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FilterModel contains the expressions what include and exclude
type FilterModel struct {
	Includes types.List `tfsdk:"includes"`
	Excludes types.List `tfsdk:"excludes"`
}

// UnarchiveDataSourceModel is unarchive data source model
type UnarchiveDataSourceModel struct {
	FileName  types.String  `tfsdk:"file_name"`
	Output    types.String  `tfsdk:"output"`
	Filters   []FilterModel `tfsdk:"filters"`
	Flat      types.Bool    `tfsdk:"flat"`
	Type      types.String  `tfsdk:"type"`
	FileNames types.List    `tfsdk:"file_names"`
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
