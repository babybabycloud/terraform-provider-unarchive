package unarchive

import (
	"os"
)

// type config struct {
// 	FileName types.String `tfsdk:"file_name"`
// 	Output   types.String `tfsdk:"output"`
// 	Includes types.List   `tfsdk:"includes"`
// 	Excludes types.List   `tfsdk:"excludes"`
// }

func (c zipFileDataSourceModel) decideOutputDir() string {
	outputDir, err := os.Getwd()
	if err != nil {
		outputDir = "./"
	}

	if !c.Output.IsNull() && !c.Output.IsUnknown() {
		outputDir = c.Output.ValueString()
	}
	return outputDir
}
