package unarchive

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const DEFAULT_DIR_MODE = 0740

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &fileDataSource{}
)

// NewFileDataSource is a helper function to simplify the provider implementation.
func NewFileDataSource() datasource.DataSource {
	return &fileDataSource{}
}

// fileDataSource is the data source implementation.
type fileDataSource struct{}

// Metadata returns the data source type name.
func (d *fileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

// Schema defines the schema for the data source.
func (d *fileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"file_name": schema.StringAttribute{
				Required: true,
			},
			"output": schema.StringAttribute{
				Optional: true,
			},
			"includes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"excludes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"flat": schema.BoolAttribute{
				Optional: true,
			},
			"file_names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *fileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config fileDataSourceModel
	diag := req.Config.Get(ctx, &config)
	if diag.HasError() {
		for _, diagnotic := range diag.Errors() {
			resp.Diagnostics.AddError(diagnotic.Summary(), diagnotic.Detail())
		}
		return
	}

	file := NewZipFileWithFileDataSourceModel(config)
	info := file.extract()
	if info.err != nil {
		resp.Diagnostics.AddError(info.msg, info.err.Error())
		return
	}

	config.addFileNames(info.filenames)
	resp.State.Set(ctx, &config)
}

type fileDataSourceModel struct {
	FileName  types.String `tfsdk:"file_name"`
	Output    types.String `tfsdk:"output"`
	Includes  types.List   `tfsdk:"includes"`
	Excludes  types.List   `tfsdk:"excludes"`
	Flat      types.Bool   `tfsdk:"flat"`
	FileNames types.List   `tfsdk:"file_names"`
}

func (c fileDataSourceModel) decideOutputDir() string {
	outputDir, err := os.Getwd()
	if err != nil {
		outputDir = "."
	}

	if !c.Output.IsNull() {
		outputDir = c.Output.ValueString()
	}
	return outputDir
}

func (c fileDataSourceModel) includePatterns() testFunc {
	if !c.Includes.IsNull() {
		return func(name string) bool {
			return toPatterns(c.Includes).doesNameMatchPatterns(name)
		}
	}
	return func(_ string) bool {
		return true
	}
}

func (c fileDataSourceModel) excludePatterns() testFunc {
	if !c.Excludes.IsNull() {
		return func(name string) bool {
			return toPatterns(c.Excludes).doesNameMatchPatterns(name)
		}
	}
	return func(_ string) bool {
		return false
	}
}

func (c fileDataSourceModel) isFlat() bool {
	if c.Flat.IsNull() {
		return false
	}
	return c.Flat.ValueBool()
}

func (z *fileDataSourceModel) addFileNames(filenames []string) {
	values := make([]attr.Value, len(filenames))
	for index, value := range filenames {
		values[index] = types.StringValue(value)
	}

	z.FileNames, _ = types.ListValue(types.StringType, values)
}
