package unarchive

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const DEFAULT_DIR_MODE = 0740

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &zipFileDataSource{}
	// _ datasource.DataSourceWithConfigure = &zipFileDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewZipFileDataSource() datasource.DataSource {
	return &zipFileDataSource{}
}

// zipFileDataSource is the data source implementation.
type zipFileDataSource struct{}

// Metadata returns the data source type name.
func (d *zipFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zip_file"
}

// Schema defines the schema for the data source.
func (d *zipFileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *zipFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config zipFileDataSourceModel

	diag := req.Config.Get(ctx, &config)
	if diag.HasError() {
		for _, diagnotic := range diag.Errors() {
			resp.Diagnostics.AddError(diagnotic.Summary(), diagnotic.Detail())
		}
		return
	}

	err := config.extract()
	if err != nil {
		resp.Diagnostics.AddError("Error occurred when extract files!", err.Error())
		return
	}

	resp.State.Set(ctx, &config)
}

type zipFileDataSourceModel struct {
	FileName types.String `tfsdk:"file_name"`
	Output   types.String `tfsdk:"output"`
	Includes types.List   `tfsdk:"includes"`
	Excludes types.List   `tfsdk:"excludes"`
}

func (z zipFileDataSourceModel) extract() error {
	reader, err := zip.OpenReader(z.FileName.ValueString())
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		err = z.copyFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (z zipFileDataSourceModel) copyFile(file *zip.File) error {
	outputDir := z.decideOutputDir()
	err := os.MkdirAll(outputDir, DEFAULT_DIR_MODE)
	if err != nil {
		return nil
	}

	if file.Method == zip.Store {
		return nil
	}

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	err = os.MkdirAll(outputDir+string(filepath.Separator)+path.Dir(file.Name), DEFAULT_DIR_MODE)
	if err != nil {
		return err
	}

	cf, err := os.Create(outputDir + string(filepath.Separator) + file.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(cf, rc)
	if err != nil {
		return err
	}

	return nil
}

func (z zipFileDataSourceModel) decideOutputDir() string {
	outputDir, err := os.Getwd()
	if err != nil {
		outputDir = "./"
	}

	if !z.Output.IsNull() && !z.Output.IsUnknown() {
		outputDir = z.Output.ValueString()
	}
	return outputDir
}
