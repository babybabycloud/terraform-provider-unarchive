package unarchive

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const DEFAULT_DIR_MODE = 0740

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &zipFileDataSource{}
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

	msg, err := config.extract(ctx)
	if err != nil {
		resp.Diagnostics.AddError(msg, err.Error())
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

func (z zipFileDataSourceModel) copyFile(file *zip.File) error {
	outputDir := z.decideOutputDir()
	err := os.MkdirAll(outputDir, DEFAULT_DIR_MODE)
	if err != nil {
		return err
	}

	if file.Method == zip.Store {
		return nil
	}

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	err = os.MkdirAll(path.Join(outputDir, string(filepath.Separator), path.Dir(file.Name)), DEFAULT_DIR_MODE)
	if err != nil {
		return err
	}

	cf, err := os.Create(path.Join(outputDir, string(filepath.Separator), file.Name))
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

func (z zipFileDataSourceModel) extract(ctx context.Context) (string, error) {
	zipFile := z.FileName.ValueString()
	rc, err := zip.OpenReader(zipFile)
	if err != nil {
		return "Error occured when open zip file", err
	}
	defer rc.Close()

	ch := filesInSliceToChan(rc.File)

	ch = filter(ch, z.patterns(z.Includes), doesNameMatchPatterns)
	ch = filter(ch, z.patterns(z.Excludes), doesNotNameMatchPatterns)

	for file := range ch {
		err = z.copyFile(file)
		if err != nil {
			return "Error occured when copy file", err
		}
	}
	return "", nil
}

func (z zipFileDataSourceModel) patterns(list types.List) []string {
	var p []string
	if !list.IsNull() && !list.IsUnknown() {
		patterns := make([]types.String, len(list.Elements()))
		p = make([]string, len(list.Elements()))
		for index, value := range patterns {
			p[index] = value.ValueString()
		}
	}
	return p
}

func filesInSliceToChan(files []*zip.File) <-chan *zip.File {
	ch := make(chan *zip.File)
	go func() {
		defer close(ch)

		for _, file := range files {
			ch <- file
		}
	}()
	return ch
}

func filter(ch <-chan *zip.File, patterns []string, test func(filename string, patterns []string) bool) <-chan *zip.File {
	outCh := make(chan *zip.File)
	go func() {
		defer close(outCh)
		for file := range ch {
			matched := test(file.Name, patterns)
			if matched {
				outCh <- file
			}
		}
	}()
	return outCh
}

func doesNameMatchPatterns(name string, patterns []string) bool {
	for _, value := range patterns {
		matched, err := regexp.MatchString(value, name)
		if err != nil {
			tflog.Warn(context.Background(), fmt.Sprintf("%s. Ignore it", err.Error()))
		}
		if matched {
			return true
		}
	}
	return false
}

func doesNotNameMatchPatterns(name string, patterns []string) bool {
	return !doesNameMatchPatterns(name, patterns)
}
